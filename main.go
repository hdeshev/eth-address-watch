package main

import (
	"context"
	"errors"
	"log/slog"

	"golang.org/x/sync/errgroup"

	"deshev.com/eth-address-watch/internal"
)

func main() {
	log := slog.Default()
	ops, ctx := errgroup.WithContext(context.Background())

	app := internal.NewApplication(ctx, log)
	log.Info("starting eth-address-watch")

	ops.Go(app.StartWatcher)
	ops.Go(app.StartServer)
	ops.Go(app.StartMonitor)

	err := ops.Wait()
	if !errors.Is(err, context.Canceled) {
		log.Error("server terminated abnormally", "boot_error", err)
	}
}
