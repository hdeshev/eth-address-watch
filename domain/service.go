package domain

import (
	"context"
	"log/slog"
	"sync"
)

type TransactionStore = map[string][]*Transaction

type Service struct {
	mtx                sync.RWMutex
	log                *slog.Logger
	blockInput         <-chan *Block
	currentBlockNumber int

	store TransactionStore
}

func NewService(log *slog.Logger, blockInput <-chan *Block) *Service {
	return &Service{
		log:                log,
		blockInput:         blockInput,
		currentBlockNumber: 0,
		store:              TransactionStore{},
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
	s.mtx.Lock()
	defer s.mtx.Unlock()

	_, exists := s.store[address]
	if !exists {
		s.store[address] = []*Transaction{}
	} else {
		s.log.Info("address already subscribed", "address", address)
	}
	return true
}

// list of inbound or outbound transactions for an address
func (s *Service) GetTransactions(address string) []*Transaction {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

	return s.store[address]
}

func (s *Service) Start(ctx context.Context) error {
	s.log.Info("starting notification service")

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

	for _, tx := range block.Transactions {
		if _, exists := s.store[tx.From]; exists {
			s.store[tx.From] = append(s.store[tx.From], tx)
		}
		if _, exists := s.store[tx.To]; exists {
			s.store[tx.To] = append(s.store[tx.To], tx)
		}
	}
}
