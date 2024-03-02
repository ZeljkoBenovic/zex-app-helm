package main

import (
	"flag"
)

type mode int

const (
	invalid mode = iota
	frontend
	backend
)

func (m *mode) toString() string {
	var modeToString = map[mode]string{
		frontend: "Frontend",
		backend:  "Backend",
	}

	return modeToString[*m]
}

func toMode(strMode string) mode {
	var stringToMode = map[string]mode{
		"frontend": frontend,
		"backend":  backend,
		"":         invalid,
	}

	return stringToMode[strMode]
}

type Config struct {
	Function  mode
	WebServer WebServer
	Log       Log
}

type WebServer struct {
	Port             string `yaml:"port"`
	Host             string `yaml:"host"`
	FEAssetsPath     string `yaml:"fe_assets_path"`
	StaticAssetsPath string `yaml:"static_assets_path"`
	MigrationsFolder string `yaml:"migrations_folder"`
}

type Log struct {
	Level string `yaml:"level"`
	Json  bool   `json:"json_output"`
}

func NewConfig() *Config {
	var (
		buffMode = new(string)
		c        = new(Config)
	)

	flag.StringVar(&c.WebServer.Host, "host", "0.0.0.0", "web server host")
	flag.StringVar(&c.WebServer.Port, "port", "8080", "web server port")
	flag.StringVar(&c.Log.Level, "log-level", "info", "log level")
	flag.BoolVar(&c.Log.Json, "log-json", false, "output logs in json format")
	flag.StringVar(buffMode, "mode", "backend", "the mode of operation (frontend, backend)")
	flag.StringVar(&c.WebServer.FEAssetsPath, "fe-assets", "frontend/*", "folder containing frontend assets (frontend/*)")
	flag.StringVar(&c.WebServer.StaticAssetsPath, "static-assets", "static", "folder containing static assets (static/*)")
	flag.StringVar(&c.WebServer.MigrationsFolder, "migrations", "sqlc/migrations", "folder containing mysql migration files")
	flag.Parse()

	c.Function = toMode(*buffMode)

	return c
}
