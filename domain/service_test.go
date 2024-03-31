package domain

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"

	"deshev.com/eth-address-watch/config"
)

func Test_Service_Start(t *testing.T) {
	log := slog.Default()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := &config.Config{}

	blockC := make(chan *Block)
	s := NewService(log, cfg, blockC)
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
	cfg := &config.Config{}
	blockC := make(chan *Block)
	s := NewService(log, cfg, blockC)

	b := &Block{
		Number:       "0x11",
		NumberParsed: 0x11,
	}
	s.processBlock(b)

	assert.Equal(t, 0x11, s.GetCurrentBlock())
}
