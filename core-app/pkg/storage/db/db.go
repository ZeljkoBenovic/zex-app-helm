package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"be/pkg/config"
	"be/pkg/db"
	"be/pkg/util"

	migrate "github.com/rubenv/sql-migrate"
)

var (
	ErrDBPortNotFound = errors.New("db port not found")
	ErrDBUserNotFound = errors.New("db user not found")
	ErrDBPassNotFound = errors.New("db pass not found")
	ErrDBHostNotFound = errors.New("db host not found")
	ErrDBNameNotFound = errors.New("db name not found")
)

type DB struct {
	ctx      context.Context
	log      *slog.Logger
	conf     config.Config
	dbq      *db.Queries
	dbCloser func() error
}

type dbParams struct {
	user, pass, name, host, port string
	k8sPassFetchErr              error
}

func NewDb(ctx context.Context, log *slog.Logger, conf config.Config) (*DB, error) {
	log.Info("Fetching database")
	dbp := fetchDbParams(conf)
	if err := dbp.checkParams(); err != nil {
		return nil, fmt.Errorf("required db params not found: %w", err)
	}

	log.Info("Attempting to connect to database")

	d, err := sql.Open("mysql", fmt.Sprintf(dbp.createDbConnString()))
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %w", err)
	}

	if err := d.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping db: %w", err)
	}

	log.Info("Successfully connected to db")

	log.Info("Running database migrations")

	numberOfMigrations, err := runMigrations(d, conf)
	if err != nil {
		return nil, fmt.Errorf("could not run migrations: %w", err)
	}

	log.Info("Database migrations completed successfully", slog.Int("total_migrations", numberOfMigrations))

	return &DB{
		ctx:      ctx,
		log:      log,
		conf:     conf,
		dbq:      db.New(d),
		dbCloser: d.Close,
	}, nil
}

func (d *DB) Close() error {
	return d.dbCloser()
}

func (d *DB) GetTitle(id int32) (string, error) {
	resp, err := d.dbq.GetTitle(d.ctx, id)
	if err != nil {
		return "", fmt.Errorf("could not get title: %w", err)
	}

	if !resp.Valid {
		return "title not found", nil
	}

	return resp.String, nil
}

func (d *DB) GetAboutMe(id int32) (string, error) {
	resp, err := d.dbq.GetAboutMe(d.ctx, id)
	if err != nil {
		return "", fmt.Errorf("could not get aboutme: %w", err)
	}

	if !resp.Valid {
		return "about me not found", nil
	}

	return resp.String, nil
}

func (d *dbParams) createDbConnString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		d.user,
		d.pass,
		d.host,
		d.port,
		d.name,
	)
}

func (d *dbParams) checkParams() error {
	if d.port == "" {
		return ErrDBPortNotFound
	}

	if d.name == "" {
		return ErrDBNameNotFound
	}

	if d.host == "" {
		return ErrDBHostNotFound
	}

	if d.user == "" {
		return ErrDBUserNotFound
	}

	if d.k8sPassFetchErr != nil {
		return fmt.Errorf("could not fetch secret from k8s secret store: %w", d.k8sPassFetchErr)
	}

	if d.pass == "" {
		return ErrDBPassNotFound
	}

	return nil
}

func fetchDbParams(conf config.Config) dbParams {
	dConf := dbParams{}

	if envHost := os.Getenv("DB_HOST"); envHost != "" {
		dConf.host = envHost
	} else {
		dConf.host = conf.DB.DBHost
	}

	if envUser := os.Getenv("DB_USER"); envUser != "" {
		dConf.user = envUser
	} else {
		dConf.user = conf.DB.DBUser
	}

	if envDbName := os.Getenv("DB_NAME"); envDbName != "" {
		dConf.name = envDbName
	} else {
		dConf.name = conf.DB.DBName
	}

	if envDbPort := os.Getenv("DB_PORT"); envDbPort != "" {
		dConf.port = envDbPort
	} else if conf.DB.DBPort != "" {
		dConf.port = conf.DB.DBPort
	} else {
		dConf.port = "3306"
	}

	if dbPass := os.Getenv("DB_PASS"); dbPass != "" {
		dConf.pass = dbPass
	} else if conf.DB.DBPass != "" {
		dConf.pass = conf.DB.DBPass
	} else {
		pass, err := util.FetchPasswordFromK8SSecrets()
		if err != nil {
			dConf.k8sPassFetchErr = err
		}

		dConf.pass = pass
	}

	return dConf
}

func runMigrations(d *sql.DB, conf config.Config) (int, error) {
	migrations := &migrate.FileMigrationSource{Dir: conf.DB.MigrationsFolder}
	n, err := migrate.Exec(d, "mysql", migrations, migrate.Up)
	if err != nil {
		return 0, err
	}

	return n, nil
}
