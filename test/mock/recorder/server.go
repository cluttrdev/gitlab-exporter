package recorder_mock

import (
	"context"
	"net"

	"google.golang.org/grpc"

	"github.com/cluttrdev/gitlab-exporter/protobuf/servicepb"
)

type Recorder struct {
	servicepb.UnimplementedGitLabExporterServer

	server    *grpc.Server
	datastore *datastore
}

func New() *Recorder {
	return &Recorder{
		server:    grpc.NewServer(),
		datastore: &datastore{},
	}
}

func (s *Recorder) Datastore() Datastore {
	return s.datastore
}

func (r *Recorder) Reset() {
	r.datastore = &datastore{}
}

func (r *Recorder) Serve(lis net.Listener) error {
	servicepb.RegisterGitLabExporterServer(r.server, r)
	return r.server.Serve(lis)
}

func (r *Recorder) GracefulStop() {
	r.server.GracefulStop()
}

func (r *Recorder) Stop() {
	r.server.Stop()
}

func (r *Recorder) RecordPipelines(ctx context.Context, req *servicepb.RecordPipelinesRequest) (*servicepb.RecordSummary, error) {
	r.datastore.pipelines = append(r.datastore.pipelines, req.Data...)
	var count int32 = int32(len(req.Data))
	return &servicepb.RecordSummary{RecordedCount: count}, nil
}

func (r *Recorder) RecordJobs(ctx context.Context, req *servicepb.RecordJobsRequest) (*servicepb.RecordSummary, error) {
	r.datastore.jobs = append(r.datastore.jobs, req.Data...)
	var count int32 = int32(len(req.Data))
	return &servicepb.RecordSummary{RecordedCount: count}, nil
}

func (r *Recorder) RecordSections(ctx context.Context, req *servicepb.RecordSectionsRequest) (*servicepb.RecordSummary, error) {
	r.datastore.sections = append(r.datastore.sections, req.Data...)
	var count int32 = int32(len(req.Data))
	return &servicepb.RecordSummary{RecordedCount: count}, nil
}

func (r *Recorder) RecordBridges(ctx context.Context, req *servicepb.RecordBridgesRequest) (*servicepb.RecordSummary, error) {
	r.datastore.bridges = append(r.datastore.bridges, req.Data...)
	var count int32 = int32(len(req.Data))
	return &servicepb.RecordSummary{RecordedCount: count}, nil
}

func (r *Recorder) RecordTraces(ctx context.Context, req *servicepb.RecordTracesRequest) (*servicepb.RecordSummary, error) {
	r.datastore.traces = append(r.datastore.traces, req.Data...)
	var count int32 = int32(len(req.Data))
	return &servicepb.RecordSummary{RecordedCount: count}, nil
}

func (r *Recorder) RecordMetrics(ctx context.Context, req *servicepb.RecordMetricsRequest) (*servicepb.RecordSummary, error) {
	r.datastore.metrics = append(r.datastore.metrics, req.Data...)
	var count int32 = int32(len(req.Data))
	return &servicepb.RecordSummary{RecordedCount: count}, nil
}

func (r *Recorder) RecordTestReports(ctx context.Context, req *servicepb.RecordTestReportsRequest) (*servicepb.RecordSummary, error) {
	r.datastore.testreports = append(r.datastore.testreports, req.Data...)
	var count int32 = int32(len(req.Data))
	return &servicepb.RecordSummary{RecordedCount: count}, nil
}

func (r *Recorder) RecordTestSuites(ctx context.Context, req *servicepb.RecordTestSuitesRequest) (*servicepb.RecordSummary, error) {
	r.datastore.testsuites = append(r.datastore.testsuites, req.Data...)
	var count int32 = int32(len(req.Data))
	return &servicepb.RecordSummary{RecordedCount: count}, nil
}

func (r *Recorder) RecordTestCases(ctx context.Context, req *servicepb.RecordTestCasesRequest) (*servicepb.RecordSummary, error) {
	r.datastore.testcases = append(r.datastore.testcases, req.Data...)
	var count int32 = int32(len(req.Data))
	return &servicepb.RecordSummary{RecordedCount: count}, nil
}
