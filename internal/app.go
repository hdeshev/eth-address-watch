package internal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"deshev.com/eth-address-watch/config"
	"deshev.com/eth-address-watch/domain"
	"deshev.com/eth-address-watch/http"
)

type Application struct {
	log    *slog.Logger
	config *config.Config
	ctx    context.Context

	service *domain.Service
	watcher *domain.Watcher
	server  *http.Server
}

func NewApplication(ctx context.Context, log *slog.Logger) *Application {
	cfg := config.New()
	service := domain.NewService(log, cfg)
	server := http.NewServer(log, service)
	watcher := domain.NewWatcher(log, cfg)

	return &Application{
		ctx:    ctx,
		log:    log,
		config: cfg,

		service: service,
		server:  server,
		watcher: watcher,
	}
}

func (a *Application) StartServer() error {
	//nolint:wrapcheck // boot errors are logged in main
	return a.server.Start(a.ctx)
}

func (a *Application) StartWatcher() error {
	//nolint:wrapcheck // boot errors are logged in main
	return a.watcher.Start(a.ctx)
}

func (a *Application) StartMonitor() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	select {
	case <-a.ctx.Done():
		return context.Canceled
	case <-quit:
		return context.Canceled
	}
}
