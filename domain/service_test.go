package domain

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Service_Start(t *testing.T) {
	log := slog.Default()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	blockC := make(chan *Block)
	s := NewService(log, blockC)
	end := make(chan struct{})
	var err error
	go func() {
		err = s.Start(ctx)
		end <- struct{}{}
	}()
	cancel()
	<-end

	assert.NoError(t, err)
}

func Test_GetCurrentBlock(t *testing.T) {
	log := slog.Default()
	blockC := make(chan *Block)
	s := NewService(log, blockC)

	b := &Block{
		Number:       "0x11",
		NumberParsed: 0x11,
	}
	s.processBlock(b)

	assert.Equal(t, 0x11, s.GetCurrentBlock())
}

func Test_Subscribe(t *testing.T) {
	log := slog.Default()
	blockC := make(chan *Block)
	s := NewService(log, blockC)

	s.Subscribe("0x1111")

	b := &Block{
		Number:       "0x11",
		NumberParsed: 0x11,
		Transactions: []*Transaction{
			{From: "0x1111", To: "0x1112"},
			{From: "0x1112", To: "0x1111"},
			{From: "0x2111", To: "0x2112"},
		},
	}
	s.processBlock(b)

	subscribedTxs := s.GetTransactions("0x1111")
	assert.Equal(t, 2, len(subscribedTxs))
	otherTxs := s.GetTransactions("0x2111")
	assert.Equal(t, 0, len(otherTxs))
}
