package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"strings"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
)

type Server struct {
	recorder servicepb.GitLabExporterServer
	health   *health.Server
	metrics  *grpcprom.ServerMetrics
}

func New(recorder servicepb.GitLabExporterServer) *Server {
	healthServer := health.NewServer()
	healthServer.SetServingStatus("" /* system */, healthpb.HealthCheckResponse_NOT_SERVING)

	metricsServer := grpcprom.NewServerMetrics()

	return &Server{
		recorder: recorder,
		health:   healthServer,
		metrics:  metricsServer,
	}
}

func (s *Server) MetricsCollector() prometheus.Collector {
	return s.metrics
}

func (s *Server) SetServingStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) {
	s.health.SetServingStatus(service, status)
}

func (s *Server) ListenAndServe(ctx context.Context, addr string) error {
	// Parse address to determine network type
	network, address := parseAddress(addr)

	// setup grpc server
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(s.metrics.UnaryServerInterceptor()),
		grpc.ChainStreamInterceptor(s.metrics.StreamServerInterceptor()),
	)

	servicepb.RegisterGitLabExporterServer(grpcServer, s.recorder)
	healthpb.RegisterHealthServer(grpcServer, s.health)
	s.metrics.InitializeMetrics(grpcServer)

	// serve and monitor health
	g := &run.Group{}

	{ // serve grpc
		g.Add(func() error { // execute
			listener, err := net.Listen(network, address)
			if err != nil {
				return err
			}
			slog.Info(fmt.Sprintf("Listening on %s://%s", network, listener.Addr().String()))

			return grpcServer.Serve(listener)
		}, func(err error) { // interrupt
			// Cleanup unix socket if needed
			if network == "unix" {
				os.Remove(address)
			}
			s.health.Shutdown()
			grpcServer.GracefulStop()
			grpcServer.Stop()
		})
	}

	// { // monitor health
	// 	ctx, cancel := context.WithCancel(ctx)
	// 	g.Add(func() error { // execute
	// 		if err := s.getReady(ctx); err != nil {
	// 			return err
	// 		}
	//
	// 		return s.watchReadiness(ctx)
	// 	}, func(err error) { // interrupt
	// 		cancel()
	// 	})
	// }

	{ // context handler
		ctx, cancel := context.WithCancel(ctx)
		g.Add(func() error { // execute
			<-ctx.Done()
			return ctx.Err()
		}, func(err error) { // interrupt
			cancel()
		})
	}

	return g.Run()
}

// parseAddress parses the address string to determine network type and address
func parseAddress(addr string) (network, address string) {
	// Handle "unix:///path/to/socket"
	if strings.HasPrefix(addr, "unix://") {
		return "unix", strings.TrimPrefix(addr, "unix://")
	}
	// Handle bare unix socket path (starts with /)
	if strings.HasPrefix(addr, "/") {
		return "unix", addr
	}
	// Default to TCP for host:port format
	return "tcp", addr
}
