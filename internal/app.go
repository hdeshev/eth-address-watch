package internal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"deshev.com/eth-address-watch/client/eth"
	"deshev.com/eth-address-watch/config"
	"deshev.com/eth-address-watch/domain"
	"deshev.com/eth-address-watch/http"
)

const (
	blockBufferSize = 10
)

type Application struct {
	log    *slog.Logger
	config *config.Config
	ctx    context.Context

	service *domain.Service
	watcher *domain.Watcher
	server  *http.Server
	blockC  chan *domain.Block
}

func NewApplication(ctx context.Context, log *slog.Logger) *Application {
	cfg := config.New()
	blockC := make(chan *domain.Block, blockBufferSize)

	service := domain.NewService(log, blockC)
	server := http.NewServer(log, service)
	client := eth.NewClient(cfg)
	watcher := domain.NewWatcher(log, cfg, client, blockC)

	return &Application{
		ctx:    ctx,
		log:    log,
		config: cfg,

		service: service,
		server:  server,
		watcher: watcher,
		blockC:  blockC,
	}
}

func (a *Application) StartAPIServer() error {
	//nolint:wrapcheck // boot errors are logged in main
	return a.server.Start(a.ctx)
}

func (a *Application) StartBlockWatcher() error {
	//nolint:wrapcheck // boot errors are logged in main
	return a.watcher.Start(a.ctx)
}

func (a *Application) StartNotificationService() error {
	//nolint:wrapcheck // boot errors are logged in main
	return a.service.Start(a.ctx)
}

func (a *Application) StartSignalMonitor() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	select {
	case <-a.ctx.Done():
		return context.Canceled
	case <-quit:
		return context.Canceled
	}
}
