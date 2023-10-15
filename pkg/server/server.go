package server

import (
	"context"
	"errors"
	"log"
	"net/http"
)

type ServerConfig struct {
	Address string
}

type Server struct {
	cfg ServerConfig
}

func New(cfg ServerConfig) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) routes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", s.MetricsHandler)

	return mux
}

func (s *Server) Serve(ctx context.Context) error {
	mux := s.routes()

	srv := &http.Server{
		Addr:    s.cfg.Address,
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Println(err)
			}
		}
	}()

	select {
	case <-ctx.Done():
		srv.Shutdown(ctx)
	}

	return nil
}
