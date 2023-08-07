package app

import (
	"github.com/metalagman/devhands/internal/devhands/config"
	"github.com/metalagman/devhands/pkg/httpserver"
	"github.com/metalagman/devhands/pkg/logger"
)

type App struct {
	config config.Config
	logger logger.Logger
	stop   chan struct{}
	server *httpserver.Server
}

func New(cfg config.Config) (*App, error) {
	l := *logger.Global()

	opts := []httpserver.ServerOption{
		httpserver.WithListenAddr(cfg.Server.ListenAddr),
		httpserver.WithReadTimeout(cfg.Server.TimeoutRead),
		httpserver.WithWriteTimeout(cfg.Server.TimeoutWrite),
		httpserver.WithIdleTimeout(cfg.Server.TimeoutIdle),

		httpserver.WithHandler(Router()),
	}

	s := httpserver.New(opts...)

	a := &App{
		config: cfg,
		logger: l,
		stop:   make(chan struct{}),
		server: s,
	}

	l.Debug().Msgf("Listening on %s", cfg.Server.ListenAddr)

	return a, nil
}

func (a *App) Stop() {
	close(a.stop)
	logger.CheckErr(a.server.Stop())
}
