package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"

	db2 "be/pkg/db"

	_ "github.com/go-sql-driver/mysql"
	migrate "github.com/rubenv/sql-migrate"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	ErrDBUserNotFound = errors.New("db username not found")
	ErrDBNameNotFound = errors.New("db name not found")
	ErrDBHostNotFound = errors.New("db host not found")

	ErrNamespaceNotFound  = errors.New("namespace variable not found")
	ErrSecretNameNotFound = errors.New("secret name not found")
)

type db struct {
	db *db2.Queries
}

func newDB(log *slog.Logger, conf *Config) (*db, error) {
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		return nil, ErrDBUserNotFound
	}
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, ErrDBHostNotFound
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, ErrDBNameNotFound
	}
	namespace := os.Getenv("NAMESPACE")
	if namespace == "" {
		return nil, ErrNamespaceNotFound
	}
	secretName := os.Getenv("SECRET_NAME")
	if secretName == "" {
		return nil, ErrNamespaceNotFound
	}

	// Initialize in-cluster lient
	cnf, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// Initialize k8s client
	cl, err := kubernetes.NewForConfig(cnf)
	if err != nil {
		return nil, err
	}

	// fetch secret
	sec, err := cl.CoreV1().Secrets(namespace).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not fetch secret: %w", err)
	}

	var dbPass string
	for _, s := range sec.Items {
		if byteSecret, ok := s.Data[secretName]; ok {
			dbPass = string(byteSecret)
		}
	}

	if dbPass == "" {
		return nil, fmt.Errorf("could not find secret in the secret store")
	}

	log.Info("Successfully fetched secret from secret store", "store_name", secretName)

	d, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbName))
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %w", err)
	}

	if err = d.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping db: %w", err)
	}

	log.Info("Database connection established - ping received")

	migNum, err := runMigrations(d, conf)
	if err != nil {
		return nil, err
	}

	log.Info("Migrations ran successfully", slog.Int("number", migNum))

	closers = append(closers, func() error {
		return d.Close()
	})

	return &db{
		db: db2.New(d),
	}, nil
}

func runMigrations(d *sql.DB, conf *Config) (int, error) {
	migrations := &migrate.FileMigrationSource{Dir: conf.WebServer.MigrationsFolder}
	n, err := migrate.Exec(d, "mysql", migrations, migrate.Up)
	if err != nil {
		return 0, err
	}

	return n, nil
}
