package web

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"be/pkg/backend"
	"be/pkg/config"
	"be/pkg/storage"

	"github.com/gin-gonic/gin"
)

type Web struct {
	conf *config.Config
	gin  *gin.Engine
	log  *slog.Logger

	handlers map[config.Mode]func()
}

var ErrModeNotFound = errors.New("mode not found")

func NewWeb(ctx context.Context, conf *config.Config, log *slog.Logger, storer storage.Storer) (*Web, error) {
	w := &Web{
		conf: conf,
		log:  log,
	}

	if conf.Log.Level == config.Debug {
		gin.SetMode(gin.DebugMode)
	}

	g := gin.Default()
	g.LoadHTMLGlob(conf.Server.FEAssetsPath)
	g.Static("/assets", conf.Server.StaticAssetsPath)

	backendGroup := g.Group("/api/v1")

	back, err := backend.NewBackend(log, conf)
	if err != nil {
		return nil, err
	}

	wh := webHandlers{
		log:     log,
		conf:    conf,
		ctx:     ctx,
		storage: storer,
		back:    back,
	}

	w.handlers = map[config.Mode]func(){
		config.Frontend: wh.frontendHandlers(g),
		config.Backend:  wh.backendHandlers(backendGroup),
	}

	w.gin = g

	return w, nil
}

func (w *Web) Serve() error {
	handlers, ok := w.handlers[w.conf.Mode]
	if !ok {
		return ErrModeNotFound
	}

	handlers()

	return w.gin.Run(fmt.Sprintf("%s:%s", w.conf.Server.Host, w.conf.Server.Port))
}
