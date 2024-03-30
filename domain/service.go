package domain

import (
	"log/slog"

	"deshev.com/eth-address-watch/config"
)

type Service struct{}

func NewService(log *slog.Logger, c *config.Config) *Service {
	return &Service{}
}

// last parsed block
func (s *Service) GetCurrentBlock() int {
	return 0
}

// add address to observer
func (s *Service) Subscribe(address string) bool {
	return true
}

type Transaction struct {
	BlockNumber int
	From        string
	To          string
	Value       string
	Gas         string
	GasPrice    string
	Input       string
}

// list of inbound or outbound transactions for an address
func (s *Service) GetTransactions(address string) []Transaction {
	return []Transaction{}
}
