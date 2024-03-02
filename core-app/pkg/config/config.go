package config

import (
	"errors"
	"flag"
)

var (
	ErrLogLevelNotSupported = errors.New("log level not supported")
	ErrModeNotFound         = errors.New("mode not found")
)

type LogLevel string
type Mode string

const (
	Debug LogLevel = "debug"
	Info  LogLevel = "info"
)
const (
	Backend  Mode = "backend"
	Frontend Mode = "frontend"
)

type Config struct {
	Mode   Mode   `yaml:"mode"`
	Log    Log    `yaml:"log"`
	Server Server `yaml:"server"`
	DB     DB     `yaml:"db"`
}

type Log struct {
	JsonOutput bool     `yaml:"json_output"`
	Level      LogLevel `yaml:"level"`
}

type Server struct {
	Host             string `yaml:"host"`
	Port             string `yaml:"port"`
	FEAssetsPath     string `yaml:"fe_assets_path"`
	StaticAssetsPath string `yaml:"static_assets_path"`
	BackendURL       string `yaml:"backend_url"`
}

type DB struct {
	DBUser           string `yaml:"user"`
	DBPass           string `yaml:"pass"`
	DBName           string `yaml:"name"`
	DBHost           string `yaml:"host"`
	DBPort           string `yaml:"port"`
	MigrationsFolder string `yaml:"migrations_folder"`
}

var logLevels = map[LogLevel]struct{}{
	Debug: {},
	Info:  {},
}
var modes = map[Mode]struct{}{
	Frontend: {},
	Backend:  {},
}

func NewConfig() (*Config, error) {
	var (
		c    = &Config{}
		ll   string
		mode string
	)

	flag.StringVar(&ll, "log-level", "info", "log level output")
	flag.BoolVar(&c.Log.JsonOutput, "json", false, "log output in json format")
	flag.StringVar(&c.Server.FEAssetsPath, "fe-assets", "frontend/*", "folder containing frontend assets (frontend/*)")
	flag.StringVar(&c.Server.StaticAssetsPath, "static-assets", "static", "folder containing static assets (static/*)")
	flag.StringVar(&c.DB.MigrationsFolder, "migrations", "sqlc/migrations", "folder containing mysql migration files")
	flag.StringVar(&mode, "mode", "backend", "the mode of operation (frontend, backend)")
	flag.StringVar(&c.DB.DBPass, "db-pass", "", "database password")
	flag.StringVar(&c.DB.DBUser, "db-user", "", "database username")
	flag.StringVar(&c.DB.DBHost, "db-host", "", "database host")
	flag.StringVar(&c.DB.DBPort, "db-port", "", "database port")
	flag.StringVar(&c.DB.DBName, "db-name", "", "database name")
	flag.StringVar(&c.Server.Host, "host", "0.0.0.0", "web server host")
	flag.StringVar(&c.Server.Port, "port", "8080", "web server port")
	flag.StringVar(&c.Server.BackendURL, "backend-url", "", "http url of the backend")

	flag.Parse()

	if _, ok := logLevels[LogLevel(ll)]; !ok {
		return nil, ErrLogLevelNotSupported
	}

	if _, ok := modes[Mode(mode)]; !ok {
		return nil, ErrModeNotFound
	}

	c.Log.Level = LogLevel(ll)
	c.Mode = Mode(mode)

	return c, nil
}
