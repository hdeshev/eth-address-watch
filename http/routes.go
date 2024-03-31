package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"deshev.com/eth-address-watch/domain"
)

type Service interface {
	GetCurrentBlock() int
	GetTransactions(address string) []domain.Transaction
	Subscribe(address string) bool
}

type Router struct {
	*http.ServeMux

	log     *slog.Logger
	service Service
}

type Response struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"-"`
	Data    any    `json:"data,omitempty"`
}

func NewRouter(log *slog.Logger, s Service) *Router {
	mux := http.NewServeMux()

	r := &Router{
		ServeMux: mux,
		log:      log,
		service:  s,
	}

	mux.HandleFunc("/block", r.GetBlock)
	mux.HandleFunc("/transactions", r.GetTransactions)
	mux.HandleFunc("/subscribe", r.Subscribe)

	return r
}

func (r *Router) GetBlock(w http.ResponseWriter, req *http.Request) {
	resp := Response{
		Data: r.service.GetCurrentBlock(),
	}
	r.writeJSON(resp, w)
}

func (r *Router) GetTransactions(w http.ResponseWriter, req *http.Request) {
	address := req.URL.Query().Get("address")
	if address == "" {
		resp := Response{
			Message: "required address field missing",
			Code:    http.StatusBadRequest,
		}
		r.writeJSON(resp, w)
		return
	}

	resp := Response{
		Data: r.service.GetTransactions(address),
	}
	r.writeJSON(resp, w)
}

func (r *Router) Subscribe(w http.ResponseWriter, req *http.Request) {
	var body struct {
		Address string `json:"address"`
	}
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil || body.Address == "" {
		resp := Response{
			Message: "invalid subscribe request",
			Code:    http.StatusBadRequest,
		}
		r.writeJSON(resp, w)
		return
	}

	resp := Response{
		Data: r.service.Subscribe(body.Address),
	}
	r.writeJSON(resp, w)
}

func (r *Router) writeJSON(resp Response, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	if resp.Code == 0 {
		if resp.Message != "" {
			resp.Code = http.StatusInternalServerError
		} else {
			resp.Code = http.StatusOK
		}
	}
	w.WriteHeader(resp.Code)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		r.log.Error("failed to write JSON response", "error", err)
	}
}
