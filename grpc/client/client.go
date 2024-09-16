package client

import (
	"context"
	"fmt"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"

	"github.com/cluttrdev/gitlab-exporter/protobuf/servicepb"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type Client struct {
	conn grpc.ClientConnInterface
	stub servicepb.GitLabExporterClient

	metrics *grpcprom.ClientMetrics
}

func NewCLient(ctx context.Context, target string, opts ...grpc.DialOption) (*Client, error) {
	metrics := grpcprom.NewClientMetrics( /* opts ...grpcprom.ClientMetricsOption */ )

	opts = append(opts,
		grpc.WithChainUnaryInterceptor(
			metrics.UnaryClientInterceptor( /* opts ..grpcprom.Option */ ),
		),
		grpc.WithChainStreamInterceptor(
			metrics.StreamClientInterceptor( /* opts ..grpcprom.Option */ ),
		),
	)

	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, err
	}

	stub := servicepb.NewGitLabExporterClient(conn)

	return &Client{
		conn:    conn,
		stub:    stub,
		metrics: metrics,
	}, nil
}

func (c *Client) MetricsCollector() prometheus.Collector {
	return c.metrics
}

func RecordCommits(c *Client, ctx context.Context, data []*typespb.Commit) error {
	req := &servicepb.RecordCommitsRequest{
		Data: data,
	}
	_, err := c.stub.RecordCommits(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording commits: %w", err)
	}

	return nil
}

func RecordBridges(c *Client, ctx context.Context, data []*typespb.Bridge) error {
	req := &servicepb.RecordBridgesRequest{
		Data: data,
	}
	_, err := c.stub.RecordBridges(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording bridges: %w", err)
	}

	return nil
}

func RecordJobs(c *Client, ctx context.Context, data []*typespb.Job) error {
	req := &servicepb.RecordJobsRequest{
		Data: data,
	}
	_, err := c.stub.RecordJobs(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording jobs: %w", err)
	}

	return nil
}

func RecordMergeRequests(c *Client, ctx context.Context, data []*typespb.MergeRequest) error {
	req := &servicepb.RecordMergeRequestsRequest{
		Data: data,
	}
	_, err := c.stub.RecordMergeRequests(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording mergerequests: %w", err)
	}

	return nil
}

func RecordMergeRequestNoteEvents(c *Client, ctx context.Context, data []*typespb.MergeRequestNoteEvent) error {
	req := &servicepb.RecordMergeRequestNoteEventsRequest{
		Data: data,
	}
	_, err := c.stub.RecordMergeRequestNoteEvents(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording merge request note events: %w", err)
	}

	return nil
}

func RecordMetrics(c *Client, ctx context.Context, data []*typespb.Metric) error {
	req := &servicepb.RecordMetricsRequest{
		Data: data,
	}
	_, err := c.stub.RecordMetrics(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording metrics: %w", err)
	}

	return nil
}

func RecordPipelines(c *Client, ctx context.Context, data []*typespb.Pipeline) error {
	req := &servicepb.RecordPipelinesRequest{
		Data: data,
	}
	_, err := c.stub.RecordPipelines(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording pipelines: %w", err)
	}

	return nil
}

func RecordProjects(c *Client, ctx context.Context, data []*typespb.Project) error {
	req := &servicepb.RecordProjectsRequest{
		Data: data,
	}
	_, err := c.stub.RecordProjects(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording projects: %w", err)
	}

	return nil
}

func RecordSections(c *Client, ctx context.Context, data []*typespb.Section) error {
	req := &servicepb.RecordSectionsRequest{
		Data: data,
	}
	_, err := c.stub.RecordSections(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording sections: %w", err)
	}

	return nil
}

func RecordTestCases(c *Client, ctx context.Context, data []*typespb.TestCase) error {
	req := &servicepb.RecordTestCasesRequest{
		Data: data,
	}
	_, err := c.stub.RecordTestCases(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording testcases: %w", err)
	}

	return nil
}

func RecordTestReports(c *Client, ctx context.Context, data []*typespb.TestReport) error {
	req := &servicepb.RecordTestReportsRequest{
		Data: data,
	}
	_, err := c.stub.RecordTestReports(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording testreports: %w", err)
	}

	return nil
}

func RecordTestSuites(c *Client, ctx context.Context, data []*typespb.TestSuite) error {
	req := &servicepb.RecordTestSuitesRequest{
		Data: data,
	}
	_, err := c.stub.RecordTestSuites(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording testsuites: %w", err)
	}

	return nil
}

func RecordTraces(c *Client, ctx context.Context, data []*typespb.Trace) error {
	req := &servicepb.RecordTracesRequest{
		Data: data,
	}
	_, err := c.stub.RecordTraces(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording traces: %w", err)
	}

	return nil
}

func RecordUsers(c *Client, ctx context.Context, data []*typespb.User) error {
	req := &servicepb.RecordUsersRequest{
		Data: data,
	}
	_, err := c.stub.RecordUsers(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("error recording users: %w", err)
	}

	return nil
}
