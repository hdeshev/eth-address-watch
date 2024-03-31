package domain

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"deshev.com/eth-address-watch/config"
)

func Test_Watcher_Start(t *testing.T) {
	log := slog.Default()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := &config.Config{}

	client := stubClient(t, 0x11, []*Block{
		{Number: "0x11"},
	})

	blockC := make(chan *Block)
	w := NewWatcher(log, cfg, client, blockC)
	end := make(chan struct{})
	var err error
	go func() {
		err = w.Start(ctx)
		end <- struct{}{}
	}()
	cancel()
	<-end

	assert.NoError(t, err)
	assert.Equal(t, 0x10, w.lastBlock)
	assert.Equal(t, 0x11, w.nextBlock)
}

func Test_Tick_SingleBlock(t *testing.T) {
	log := slog.Default()

	cfg := &config.Config{}

	client := stubClient(t, 0x11, []*Block{
		{
			Number: "0x11",
			Transactions: []Transaction{
				{From: "0x1111", To: "0x1112"},
			},
		},
	})

	blockC := make(chan *Block, 1)
	w := NewWatcher(log, cfg, client, blockC)
	w.lastBlock = 0x10
	w.nextBlock = 0x10

	w.tick()

	assert.Equal(t, 0x11, w.lastBlock)
	assert.Equal(t, 0x11, w.nextBlock)

	close(blockC)
	block, hasBlock := <-blockC
	assert.True(t, hasBlock)
	assert.Equal(t, 0x11, block.NumberParsed)
	tx := block.Transactions[0]
	assert.Equal(t, "0x1111", tx.From)
	assert.Equal(t, "0x1112", tx.To)
	_, hasBlock = <-blockC
	assert.False(t, hasBlock)
}

func Test_Tick_CatchupMultipleBlocks(t *testing.T) {
	log := slog.Default()

	cfg := &config.Config{}

	client := stubClient(t, 0x12, []*Block{
		{
			Number: "0x11",
			Transactions: []Transaction{
				{From: "0x1111", To: "0x1112"},
			},
		},
		{
			Number: "0x12",
			Transactions: []Transaction{
				{From: "0x2111", To: "0x2112"},
			},
		},
	})

	blockC := make(chan *Block, 2)
	w := NewWatcher(log, cfg, client, blockC)
	w.lastBlock = 0x10
	w.nextBlock = 0x10

	w.tick()

	assert.Equal(t, 0x12, w.lastBlock)
	assert.Equal(t, 0x12, w.nextBlock)

	close(blockC)
	block1, hasBlock := <-blockC
	assert.True(t, hasBlock)
	assert.Equal(t, 0x11, block1.NumberParsed)
	tx1 := block1.Transactions[0]
	assert.Equal(t, "0x1111", tx1.From)
	assert.Equal(t, "0x1112", tx1.To)

	block2, hasBlock := <-blockC
	assert.True(t, hasBlock)
	assert.Equal(t, 0x12, block2.NumberParsed)
	tx2 := block2.Transactions[0]
	assert.True(t, hasBlock)
	assert.Equal(t, "0x2111", tx2.From)
	assert.Equal(t, "0x2112", tx2.To)
	_, hasBlock = <-blockC
	assert.False(t, hasBlock)
}

func Test_Tick_NoopIfCaughtUp(t *testing.T) {
	log := slog.Default()

	cfg := &config.Config{}

	client := stubClient(t, 0x12, []*Block{
		{
			Number: "0x12",
			Transactions: []Transaction{
				{From: "0x2111", To: "0x2112"},
			},
		},
	})

	blockC := make(chan *Block, 2)
	w := NewWatcher(log, cfg, client, blockC)
	w.lastBlock = 0x12
	w.nextBlock = 0x12

	w.tick()

	assert.Equal(t, 0x12, w.lastBlock)
	assert.Equal(t, 0x12, w.nextBlock)

	close(blockC)
	_, hasBlock := <-blockC
	assert.False(t, hasBlock)
}

func stubClient(t *testing.T, lastBlock int, blocks []*Block) ETHClient {
	t.Helper()

	client := &MockETHClient{}
	client.On("GetLatestBlock", mock.Anything).Return(lastBlock, nil)
	for _, block := range blocks {
		blockNumber := 0
		_, err := fmt.Sscanf(block.Number, "0x%x", &blockNumber)
		assert.NoError(t, err)

		client.On("GetBlock", mock.Anything, blockNumber).Return(block, nil)
	}

	client.On("GetBlock", mock.Anything, 101).Return(&Block{Number: "101"}, nil)
	return client
}

type MockETHClient struct {
	mock.Mock
}

func (m *MockETHClient) GetLatestBlock(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockETHClient) GetBlock(ctx context.Context, blockNumber int) (*Block, error) {
	args := m.Called(ctx, blockNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Block), args.Error(1)
}
