package client_test

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	grpc_client "github.com/cluttrdev/gitlab-exporter/grpc/client"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"

	grpc_mock "github.com/cluttrdev/gitlab-exporter/test/mock/grpc"
)

const bufSize int = 1024 * 1024

func newServerAndClient() (*grpc_mock.MockExporterServer, *grpc_client.Client, error) {
	server := grpc_mock.NewMockExporterServer()

	listener := bufconn.Listen(bufSize)
	go func() {
		if err := server.Serve(listener); err != nil {
			fmt.Fprint(os.Stderr, err)
		}
	}()

	client, err := grpc_client.NewCLient(
		"passthrough://bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return listener.Dial()
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	return server, client, err
}

func Test_RecordPipelines(t *testing.T) {
	server, client, err := newServerAndClient()
	defer server.GracefulStop()
	if err != nil {
		t.Error(err)
	}

	data := []*typespb.Pipeline{
		{
			Id: 42,
		},
	}

	server.ExpectPipelines(data)

	ctx := context.Background()
	if err := grpc_client.RecordPipelines(client, ctx, data); err != nil {
		t.Error(err)
	}
}

func Test_RecordJobs(t *testing.T) {
	server, client, err := newServerAndClient()
	defer server.GracefulStop()
	if err != nil {
		t.Error(err)
	}

	data := []*typespb.Job{
		{
			Id:       42,
			Pipeline: nil,
		},
	}

	server.ExpectJobs(data)

	ctx := context.Background()
	if err := grpc_client.RecordJobs(client, ctx, data); err != nil {
		t.Error(err)
	}
}
