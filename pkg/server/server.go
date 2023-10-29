package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/pprof"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/internal/healthz"
)

type ServerConfig struct {
	Address string

	LivenessCheck  healthz.Check
	ReadinessCheck healthz.Check

	Debug bool
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

	// health check endpoints
	health := healthz.NewHandler()
	health.SetLivenessCheck(s.cfg.LivenessCheck)
	health.SetReadinessCheck(s.cfg.ReadinessCheck)

	mux.HandleFunc("/healthz/live", health.LiveEndpoint)
	mux.HandleFunc("/healthz/ready", health.ReadyEndpoint)

	// metrics endpoint
	mux.HandleFunc("/metrics", s.MetricsHandler)

	// debug endpoints
	if s.cfg.Debug {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}

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

	<-ctx.Done()
	return srv.Shutdown(ctx)
}
