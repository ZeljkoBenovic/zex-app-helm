package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidMode = errors.New("invalid mode selected")
	closers        []func() error
)

func main() {
	conf := NewConfig()
	log := setupLogger(conf)

	if err := RunWebServer(conf, log); err != nil {
		log.Error(err.Error())
	}
}

// RunWebServer will run the web server using provided config.
func RunWebServer(conf *Config, log *slog.Logger) error {

	if conf.Function == invalid {
		return ErrInvalidMode
	}

	if conf.Log.Level != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	g := gin.Default()
	g.LoadHTMLGlob(conf.WebServer.FEAssetsPath)
	g.Static("/static", conf.WebServer.StaticAssetsPath)
	be := g.Group("/api/v1")

	switch conf.Function {
	case frontend:
		if os.Getenv("BE_URL") == "" {
			log.Error("Backend URL is not defined", slog.String("env_var", "BE_URL"))
			os.Exit(1)
		}

		for _, hn := range handlersFE {
			g.Handle(hn(log, nil, conf))
		}
	case backend:
		d, err := newDB(log, conf)
		if err != nil {
			return err
		}

		for _, hn := range handlersBE {
			be.Handle(hn(log, d, conf))
		}
	}

	log.Info(fmt.Sprintf("%s server started", conf.Function.toString()), slog.String("host", conf.WebServer.Host), slog.String("port", conf.WebServer.Port))

	go handleTermination(log)

	return g.Run(fmt.Sprintf("%s:%s", conf.WebServer.Host, conf.WebServer.Port))
}

func setupLogger(conf *Config) *slog.Logger {
	var (
		logLvl = map[string]slog.Level{
			"info":  slog.LevelInfo,
			"debug": slog.LevelDebug,
		}
		log = &slog.Logger{}
	)

	switch conf.Log.Json {
	case true:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLvl[conf.Log.Level],
		}))
	default:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLvl[conf.Log.Level],
		}))
	}

	return log
}

func handleTermination(slog *slog.Logger) {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)

	<-sig
	slog.Debug("Received SIGTERM signal")
	slog.Info("Started shutdown process...")
	for _, cloerFn := range closers {
		if err := cloerFn(); err != nil {
			slog.Error("Closer func error", "err", err.Error())
		}
	}

	slog.Info("Shutdown process complete")
	os.Exit(0)
}
