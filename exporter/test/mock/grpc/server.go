package grpc_mock

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"

	"go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

type MockExporterServer struct {
	servicepb.UnimplementedGitLabExporterServer
	server *grpc.Server

	expectedPipelines []*typespb.Pipeline
	expectedJobs      []*typespb.Job
}

func NewMockExporterServer() *MockExporterServer {
	return &MockExporterServer{
		server: grpc.NewServer(),
	}
}

func (ms *MockExporterServer) Serve(lis net.Listener) error {
	servicepb.RegisterGitLabExporterServer(ms.server, ms)

	return ms.server.Serve(lis)
}

func (ms *MockExporterServer) GracefulStop() {
	ms.server.GracefulStop()
}

func (ms *MockExporterServer) Stop() {
	ms.server.Stop()
}

func (ms *MockExporterServer) ExpectPipelines(data []*typespb.Pipeline) {
	ms.expectedPipelines = data
}

func (ms *MockExporterServer) ExpectJobs(data []*typespb.Job) {
	ms.expectedJobs = data
}

func check[T any](srv *MockExporterServer, want []*T, got []*T, opts ...cmp.Option) (int32, error) {
	var errs error

	var recorded int32 = 0
	for i, got := range got {
		a, _ := json.Marshal(want[i])
		b, _ := json.Marshal(got)

		if diff := cmp.Diff(a, b, opts...); diff != "" {
			errs = errors.Join(errs, fmt.Errorf("Mismatch (-want, +got):\n%s", diff))
		} else {
			recorded++
		}
	}

	return recorded, errs
}

func (ms *MockExporterServer) RecordPipelines(ctx context.Context, r *servicepb.RecordPipelinesRequest) (*servicepb.RecordSummary, error) {
	n, err := check(ms, ms.expectedPipelines, r.Data)
	return &servicepb.RecordSummary{RecordedCount: n}, err
}

func (ms *MockExporterServer) RecordJobs(ctx context.Context, r *servicepb.RecordJobsRequest) (*servicepb.RecordSummary, error) {
	n, err := check(ms, ms.expectedJobs, r.Data)
	return &servicepb.RecordSummary{RecordedCount: n}, err
}
