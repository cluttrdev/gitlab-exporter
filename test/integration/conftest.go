package integration_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"go.cluttr.dev/gitlab-exporter/internal/exporter"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"go.cluttr.dev/gitlab-exporter/test/mock/recorder"
)

func setupGitLab(t *testing.T) (*http.ServeMux, *gitlab.Client) {
	mux := http.NewServeMux()
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	client, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL: srv.URL,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	return mux, client
}

func setupExporter(t *testing.T) (*exporter.Exporter, *recorder_mock.Recorder) {
	rec := recorder_mock.New()
	t.Cleanup(rec.GracefulStop)

	const bufSize int = 4 * 1024 * 1024
	lis := bufconn.Listen(bufSize)
	go func() {
		if err := rec.Serve(lis); err != nil {
			t.Log(err)
		}
	}()

	exp, err := exporter.New([]exporter.EndpointConfig{
		{
			Address: "bufnet",
			Options: []grpc.DialOption{
				grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
					return lis.Dial()
				}),
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			},
		},
	})
	if err != nil {
		t.Fatalf("failed to create exporter: %v", err)
	}

	return exp, rec
}
