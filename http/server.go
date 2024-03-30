package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"deshev.com/eth-address-watch/domain"
)

const (
	port        = 9000
	httpTimeout = 10 * time.Second
)

type Server struct {
	log  *slog.Logger
	http *http.Server
}

func NewServer(log *slog.Logger, s *domain.Service) *Server {
	r := NewRouter(log, s)
	return &Server{
		log: log,
		http: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      r,
			ReadTimeout:  httpTimeout,
			WriteTimeout: httpTimeout,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Info("starting http server", "port", port)
	errCh := make(chan error, 1)
	go func() {
		err := s.http.ListenAndServe()
		if err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		s.Stop()
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) Stop() {
	s.log.Info("stopping http server")
	err := s.http.Shutdown(context.Background())
	if err != nil {
		s.log.Error("http: server: shutdown", err)
	}
}
