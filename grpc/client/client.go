package client

import (
	"context"
	"fmt"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"

	"go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

type Client struct {
	conn grpc.ClientConnInterface
	stub servicepb.GitLabExporterClient

	metrics *grpcprom.ClientMetrics
}

func NewCLient(target string, opts ...grpc.DialOption) (*Client, error) {
	metrics := grpcprom.NewClientMetrics( /* opts ...grpcprom.ClientMetricsOption */ )

	opts = append(opts,
		grpc.WithChainUnaryInterceptor(
			metrics.UnaryClientInterceptor( /* opts ..grpcprom.Option */ ),
		),
		grpc.WithChainStreamInterceptor(
			metrics.StreamClientInterceptor( /* opts ..grpcprom.Option */ ),
		),
	)

	conn, err := grpc.NewClient(target, opts...)
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
		return fmt.Errorf("record commits: %w", err)
	}

	return nil
}

func RecordCoverageReports(c *Client, ctx context.Context, data []*typespb.CoverageReport) error {
	req := &servicepb.RecordCoverageReportsRequest{
		Data: data,
	}
	_, err := c.stub.RecordCoverageReports(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record coverage reports: %w", err)
	}

	return nil
}

func RecordCoveragePackages(c *Client, ctx context.Context, data []*typespb.CoveragePackage) error {
	req := &servicepb.RecordCoveragePackagesRequest{
		Data: data,
	}
	_, err := c.stub.RecordCoveragePackages(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record coverage packages: %w", err)
	}

	return nil
}

func RecordCoverageClasses(c *Client, ctx context.Context, data []*typespb.CoverageClass) error {
	req := &servicepb.RecordCoverageClassesRequest{
		Data: data,
	}
	_, err := c.stub.RecordCoverageClasses(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record coverage classes: %w", err)
	}

	return nil
}

func RecordCoverageMethods(c *Client, ctx context.Context, data []*typespb.CoverageMethod) error {
	req := &servicepb.RecordCoverageMethodsRequest{
		Data: data,
	}
	_, err := c.stub.RecordCoverageMethods(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record coverage methods: %w", err)
	}

	return nil
}

func RecordDeployments(c *Client, ctx context.Context, data []*typespb.Deployment) error {
	req := &servicepb.RecordDeploymentsRequest{
		Data: data,
	}
	_, err := c.stub.RecordDeployments(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record deployments: %w", err)
	}

	return nil
}

func RecordIssues(c *Client, ctx context.Context, data []*typespb.Issue) error {
	req := &servicepb.RecordIssuesRequest{
		Data: data,
	}
	_, err := c.stub.RecordIssues(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record issues: %w", err)
	}

	return nil
}

func RecordJobs(c *Client, ctx context.Context, data []*typespb.Job) error {
	req := &servicepb.RecordJobsRequest{
		Data: data,
	}
	_, err := c.stub.RecordJobs(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record jobs: %w", err)
	}

	return nil
}

func RecordMergeRequests(c *Client, ctx context.Context, data []*typespb.MergeRequest) error {
	req := &servicepb.RecordMergeRequestsRequest{
		Data: data,
	}
	_, err := c.stub.RecordMergeRequests(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record mergerequests: %w", err)
	}

	return nil
}

func RecordMergeRequestNoteEvents(c *Client, ctx context.Context, data []*typespb.MergeRequestNoteEvent) error {
	req := &servicepb.RecordMergeRequestNoteEventsRequest{
		Data: data,
	}
	_, err := c.stub.RecordMergeRequestNoteEvents(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record merge request note events: %w", err)
	}

	return nil
}

func RecordMetrics(c *Client, ctx context.Context, data []*typespb.Metric) error {
	req := &servicepb.RecordMetricsRequest{
		Data: data,
	}
	_, err := c.stub.RecordMetrics(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record metrics: %w", err)
	}

	return nil
}

func RecordPipelines(c *Client, ctx context.Context, data []*typespb.Pipeline) error {
	req := &servicepb.RecordPipelinesRequest{
		Data: data,
	}
	_, err := c.stub.RecordPipelines(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record pipelines: %w", err)
	}

	return nil
}

func RecordProjects(c *Client, ctx context.Context, data []*typespb.Project) error {
	req := &servicepb.RecordProjectsRequest{
		Data: data,
	}
	_, err := c.stub.RecordProjects(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record projects: %w", err)
	}

	return nil
}

func RecordSections(c *Client, ctx context.Context, data []*typespb.Section) error {
	req := &servicepb.RecordSectionsRequest{
		Data: data,
	}
	_, err := c.stub.RecordSections(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record sections: %w", err)
	}

	return nil
}

func RecordTestCases(c *Client, ctx context.Context, data []*typespb.TestCase) error {
	req := &servicepb.RecordTestCasesRequest{
		Data: data,
	}
	_, err := c.stub.RecordTestCases(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record testcases: %w", err)
	}

	return nil
}

func RecordTestReports(c *Client, ctx context.Context, data []*typespb.TestReport) error {
	req := &servicepb.RecordTestReportsRequest{
		Data: data,
	}
	_, err := c.stub.RecordTestReports(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record testreports: %w", err)
	}

	return nil
}

func RecordTestSuites(c *Client, ctx context.Context, data []*typespb.TestSuite) error {
	req := &servicepb.RecordTestSuitesRequest{
		Data: data,
	}
	_, err := c.stub.RecordTestSuites(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record testsuites: %w", err)
	}

	return nil
}

func RecordTraces(c *Client, ctx context.Context, data []*typespb.Trace) error {
	req := &servicepb.RecordTracesRequest{
		Data: data,
	}
	_, err := c.stub.RecordTraces(ctx, req /* opts ...grpc.CallOption */)
	if err != nil {
		return fmt.Errorf("record traces: %w", err)
	}

	return nil
}
