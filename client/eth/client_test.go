package eth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"deshev.com/eth-address-watch/config"
)

func TestClient_GetLatestBlock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/", r.URL.Path)

		var req map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)

		assert.Equal(t, "eth_blockNumber", req["method"])
		assert.Equal(t, nil, req["params"])

		response := map[string]interface{}{
			"result": "0x1234",
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	cfg := &config.Config{
		EthNodeURL:        server.URL,
		EthRequestTimeout: 1 * time.Second,
	}
	client := NewClient(cfg)

	block, err := client.GetLatestBlock(context.Background())
	require.NoError(t, err)

	assert.Equal(t, 0x1234, block)
}

func TestClient_GetBlock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/", r.URL.Path)

		var req map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&req)
		assert.NoError(t, err)

		assert.Equal(t, "eth_getBlockByNumber", req["method"])
		assert.Equal(t, []interface{}{"0x1234", true}, req["params"])

		response := map[string]interface{}{
			"result": map[string]interface{}{
				"number": "0x1234",
				"hash":   "0x1234567890abcdef",
				"transactions": []any{
					map[string]interface{}{
						"from": "0x1234",
						"to":   "0x5678",
					},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		assert.NoError(t, err)
	}))
	defer server.Close()

	cfg := &config.Config{
		EthNodeURL:        server.URL,
		EthRequestTimeout: 1 * time.Second,
	}
	client := NewClient(cfg)

	block, err := client.GetBlock(context.Background(), 0x1234)
	require.NoError(t, err)

	assert.Equal(t, "0x1234", block.Number)
	assert.Equal(t, "0x1234567890abcdef", block.Hash)
	assert.Equal(t, 1, len(block.Transactions))
	assert.Equal(t, "0x1234", block.Transactions[0].From)
	assert.Equal(t, "0x5678", block.Transactions[0].To)
}
