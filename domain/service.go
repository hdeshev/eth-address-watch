package domain

import (
	"context"
	"log/slog"
	"sync"

	"deshev.com/eth-address-watch/config"
)

type Service struct {
	mtx                sync.RWMutex
	log                *slog.Logger
	cfg                *config.Config
	blockInput         <-chan *Block
	currentBlockNumber int
}

func NewService(log *slog.Logger, c *config.Config, blockInput <-chan *Block) *Service {
	return &Service{
		log:                log,
		cfg:                c,
		blockInput:         blockInput,
		currentBlockNumber: 0,
	}
}

// last parsed block
func (s *Service) GetCurrentBlock() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.currentBlockNumber
}

// add address to observer
func (s *Service) Subscribe(address string) bool {
	return true
}

// list of inbound or outbound transactions for an address
func (s *Service) GetTransactions(address string) []Transaction {
	return []Transaction{}
}

func (s *Service) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case block := <-s.blockInput:
			s.processBlock(block)
		}
	}
}

func (s *Service) processBlock(block *Block) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	s.log.Info("service processing block", "block", block.NumberParsed, "transactions", len(block.Transactions))
	s.currentBlockNumber = block.NumberParsed
}
