package domain

import (
	"context"
	"log/slog"

	"deshev.com/eth-address-watch/config"
)

type Watcher struct {
	log    *slog.Logger
	config *config.Config
}

func NewWatcher(log *slog.Logger, cfg *config.Config) *Watcher {
	return &Watcher{
		log:    log,
		config: cfg,
	}
}

func (w *Watcher) Start(ctx context.Context) error {
	w.log.Info("starting watcher")
	return nil
}
