package domain

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"deshev.com/eth-address-watch/config"
)

type Watcher struct {
	log       *slog.Logger
	config    *config.Config
	ethClient ETHClient

	blockOutput chan<- *Block
	lastBlock   int
	nextBlock   int
}

type ETHClient interface {
	GetLatestBlock(ctx context.Context) (int, error)
	GetBlock(ctx context.Context, blockNumber int) (*Block, error)
}

const (
	blockTickInterval = 10 * time.Second
	tickTimeout       = 5 * time.Second
)

func NewWatcher(log *slog.Logger, cfg *config.Config, client ETHClient, blockOutput chan<- *Block) *Watcher {
	return &Watcher{
		log:         log,
		config:      cfg,
		ethClient:   client,
		blockOutput: blockOutput,
		lastBlock:   0,
		nextBlock:   0,
	}
}

func (w *Watcher) Start(ctx context.Context) error {
	w.log.Info("starting watcher")
	blockNumber, err := w.ethClient.GetLatestBlock(ctx)
	if err != nil {
		return fmt.Errorf("error getting latest block: %w", err)
	}

	w.nextBlock = blockNumber
	w.lastBlock = blockNumber - 1
	w.log.Info("next block", "block", blockNumber)

	ticker := time.NewTicker(blockTickInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			w.tick()
		}
	}
}

func (w *Watcher) tick() {
	w.log.Info("ethereum watcher tick")

	ctx, cancel := context.WithTimeout(context.Background(), tickTimeout)
	defer cancel()

	blockNumber, err := w.ethClient.GetLatestBlock(ctx)
	if err != nil {
		w.log.Error("error getting latest block", "error", err)
		return
	}

	w.nextBlock = blockNumber
	for i := w.lastBlock + 1; i <= w.nextBlock; i++ {
		w.log.Info("ethereum watcher processing block", "block", i)

		block, err := w.ethClient.GetBlock(ctx, i)
		if err != nil {
			w.log.Error("error getting block", "block", i, "error", err)
			return
		}

		w.lastBlock = i
		number := 0
		_, err = fmt.Sscanf(block.Number, "0x%x", &number)
		if err != nil {
			w.log.Error("error parsing block number", "block", i, "raw_value", block.Number, "error", err)
			return
		}
		block.NumberParsed = number
		w.blockOutput <- block
	}
}
