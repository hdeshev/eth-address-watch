package eth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"deshev.com/eth-address-watch/config"
	"deshev.com/eth-address-watch/domain"
)

type Client struct {
	nodeURL        string
	requestTimeout time.Duration
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		nodeURL:        cfg.EthNodeURL,
		requestTimeout: cfg.EthRequestTimeout,
	}
}

func (c *Client) GetBlock(ctx context.Context, blockNumber int) (*domain.Block, error) {
	req := blockByNumberCall(blockNumber)
	var result struct {
		Result domain.Block `json:"result"`
	}
	err := jsonRPCRequest(ctx, c, req, &result)
	if err != nil {
		return nil, err
	}

	return &result.Result, nil
}

func (c *Client) GetLatestBlock(ctx context.Context) (int, error) {
	req := blockNumberCall()
	var result struct {
		Result string `json:"result"`
	}
	err := jsonRPCRequest(ctx, c, req, &result)
	if err != nil {
		return 0, err
	}

	blockNumber := 0
	_, err = fmt.Sscanf(result.Result, "0x%x", &blockNumber)
	if err != nil {
		return 0, fmt.Errorf("block number parse error: %w", err)
	}
	return blockNumber, nil
}

func jsonRPCRequest[R any](ctx context.Context, c *Client, call rpcMethodCall, result *R) error {
	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	payload := bytes.Buffer{}
	if err := json.NewEncoder(&payload).Encode(call); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.nodeURL, &payload)
	if err != nil {
		return fmt.Errorf("http request create error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request execute error: %w", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("http response parse error: %w", err)
	}
	return nil
}

type rpcMethodCall struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	ID      int    `json:"id"`
}

func blockNumberCall() rpcMethodCall {
	return rpcMethodCall{
		JSONRPC: "2.0",
		Method:  "eth_blockNumber",
		ID:      1,
	}
}

func blockByNumberCall(blockNumber int) rpcMethodCall {
	hexBlockNumber := fmt.Sprintf("0x%x", blockNumber)
	return rpcMethodCall{
		JSONRPC: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []any{hexBlockNumber, true},
		ID:      1,
	}
}
