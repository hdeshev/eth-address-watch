package http

import (
	"log/slog"
	"net/http"
)

type Router struct {
	*http.ServeMux

	log *slog.Logger
}

type Service interface{}

func NewRouter(log *slog.Logger, s Service) *Router {
	mux := http.NewServeMux()

	r := &Router{
		ServeMux: mux,
		log:      log,
	}

	mux.HandleFunc("/block", r.Block)

	return r
}

func (r *Router) Block(w http.ResponseWriter, req *http.Request) {
	_, err := w.Write([]byte("TODO"))
	r.log.Info("block", "write_error", err)
}
