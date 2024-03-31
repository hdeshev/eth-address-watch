package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"deshev.com/eth-address-watch/domain"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetCurrentBlock() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockService) GetTransactions(address string) []domain.Transaction {
	args := m.Called()
	return args.Get(0).([]domain.Transaction)
}

func (m *MockService) Subscribe(address string) bool {
	args := m.Called()
	return args.Bool(0)
}

func Test_GetBlock(t *testing.T) {
	log := slog.Default()

	mockService := new(MockService)
	mockService.On("GetCurrentBlock").Return(1)

	router := NewRouter(log, mockService)

	req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/block", http.NoBody)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var parsedBody Response
	err := json.NewDecoder(rr.Body).Decode(&parsedBody)
	assert.NoError(t, err)
	assert.Equal(t, 1.0, parsedBody.Data)
}

func Test_GeTransactions(t *testing.T) {
	log := slog.Default()

	mockService := new(MockService)
	mockService.On("GetTransactions").Return([]domain.Transaction{
		{
			From:     "address1",
			To:       "address2",
			Value:    "1",
			Gas:      "1",
			GasPrice: "1",
		},
		{
			From:     "address1",
			To:       "address3",
			Value:    "1",
			Gas:      "1",
			GasPrice: "1",
		},
	})

	tests := []struct {
		name       string
		address    string
		wantStatus int
		wantTxs    string
		wantError  string
	}{
		{
			name:       "get transactions",
			address:    "address-1",
			wantStatus: http.StatusOK,
			wantTxs:    "address1->address2|address1->address3",
		},
		{
			name:       "invalid address",
			address:    "",
			wantStatus: http.StatusBadRequest,
			wantError:  "required address field missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := NewRouter(log, mockService)

			req, _ := http.NewRequestWithContext(context.TODO(), "GET", "/transactions?address="+tt.address, http.NoBody)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			var parsedBody Response
			err := json.NewDecoder(rr.Body).Decode(&parsedBody)
			assert.NoError(t, err)

			if tt.wantError != "" {
				assert.Equal(t, tt.wantError, parsedBody.Message)
			} else {
				txs, ok := parsedBody.Data.([]any)
				assert.True(t, ok)
				assert.Equal(t, tt.wantTxs, formatTxs(t, txs))
			}
		})
	}
}

func formatTxs(t *testing.T, txs []any) string {
	t.Helper()

	result := []string{}
	for i := range txs {
		m, ok := txs[i].(map[string]any)
		if !ok {
			assert.Fail(t, "unexpected value")
		}
		result = append(result, fmt.Sprintf("%s->%s", m["from"], m["to"]))
	}
	return strings.Join(result, "|")
}

func Test_Subscribe(t *testing.T) {
	tests := []struct {
		name        string
		requestBody string
		wantStatus  int
		wantResult  bool
		wantError   string
	}{
		{
			name:        "valid subscription",
			requestBody: `{"address":"address-1"}`,
			wantStatus:  http.StatusOK,
			wantResult:  true,
		},
		{
			name:        "invalid json",
			requestBody: "",
			wantStatus:  http.StatusBadRequest,
			wantError:   "invalid subscribe request",
		},
		{
			name:        "invalid request schema",
			requestBody: `{"missing-address-field":"dummy-data"}`,
			wantStatus:  http.StatusBadRequest,
			wantError:   "invalid subscribe request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.Default()

			mockService := new(MockService)
			mockService.On("Subscribe").Return(tt.wantResult)

			router := NewRouter(log, mockService)

			body := bytes.NewBufferString(tt.requestBody)
			req, _ := http.NewRequestWithContext(context.TODO(), "POST", "/subscribe", body)
			rr := httptest.NewRecorder()
			router.ServeHTTP(rr, req)

			assert.Equal(t, tt.wantStatus, rr.Code)
			var parsedBody Response
			err := json.NewDecoder(rr.Body).Decode(&parsedBody)
			assert.NoError(t, err)

			result := parsedBody.Data
			if tt.wantError != "" {
				assert.Equal(t, tt.wantError, parsedBody.Message)
			} else {
				assert.Equal(t, tt.wantResult, result)
			}
		})
	}
}
