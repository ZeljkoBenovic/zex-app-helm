package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"be/pkg/config"
	"be/pkg/storage"
	"be/pkg/web"
)

type Server interface {
	Serve() error
}

type App struct {
	ctx    context.Context
	conf   *config.Config
	log    *slog.Logger
	store  *storage.Storage
	webSrv Server
}

func NewApp() (*App, error) {
	a := &App{
		ctx: context.Background(),
	}

	conf, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	a.conf = conf

	a.setupLogger()

	store, err := storage.NewMockStorage()
	if err != nil {
		return nil, err
	}

	w, err := web.NewWeb(a.ctx, a.conf, a.log, store.DB())
	if err != nil {
		return nil, err
	}

	a.webSrv = w
	a.store = store

	go a.handleTerminationSignals()

	return a, nil
}

func (a *App) Run() error {
	a.log.Info("Starting application...")
	return a.webSrv.Serve()
}

func (a *App) setupLogger() {
	var (
		s      *slog.Logger
		logLvl = map[config.LogLevel]slog.Level{
			config.Debug: slog.LevelDebug,
			config.Info:  slog.LevelInfo,
		}
	)
	if a.conf.Log.JsonOutput {
		s = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLvl[a.conf.Log.Level],
		}))
	} else {
		s = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLvl[a.conf.Log.Level],
		}))
	}

	a.log = s
}

func (a *App) handleTerminationSignals() {
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Kill, os.Interrupt)

	<-sig
	a.log.Info("Application shutting down...")
	if err := a.store.Close(); err != nil {
		a.log.Error("Error closing store", "err", err.Error())
		return
	}

	a.log.Info("Application shutdown successful")
	os.Exit(0)
}
