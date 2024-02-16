package client

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/cluttrdev/gitlab-exporter/grpc/exporterpb"
)

type EndpointConfig struct {
	Address string
	Options []grpc.DialOption
}

type Client struct {
	conn   grpc.ClientConnInterface
	client pb.GitLabExporterClient
}

func NewCLient(cfg EndpointConfig) (*Client, error) {
	conn, err := grpc.Dial(cfg.Address, cfg.Options...)
	if err != nil {
		return nil, err
	}

	client := pb.NewGitLabExporterClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

func send[T any](stream grpc.ClientStream, data []*T) error {
	for _, msg := range data {
		if err := stream.SendMsg(msg); err != nil {
			return err
		}
	}

	if err := stream.CloseSend(); err != nil {
		return err
	}

	m := new(pb.RecordSummary)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}

func RecordPipelines(c *Client, ctx context.Context, data []*pb.Pipeline) error {
	stream, err := c.client.RecordPipelines(ctx)
	if err != nil {
		return err
	}

	return send[pb.Pipeline](stream, data)
}

func RecordJobs(c *Client, ctx context.Context, data []*pb.Job) error {
	stream, err := c.client.RecordJobs(ctx)
	if err != nil {
		return err
	}

	return send[pb.Job](stream, data)
}

func RecordSections(c *Client, ctx context.Context, data []*pb.Section) error {
	stream, err := c.client.RecordSections(ctx)
	if err != nil {
		return err
	}

	return send[pb.Section](stream, data)
}

func RecordBridges(c *Client, ctx context.Context, data []*pb.Bridge) error {
	stream, err := c.client.RecordBridges(ctx)
	if err != nil {
		return err
	}

	return send[pb.Bridge](stream, data)
}

func RecordTestReports(c *Client, ctx context.Context, data []*pb.TestReport) error {
	stream, err := c.client.RecordTestReports(ctx)
	if err != nil {
		return err
	}

	return send[pb.TestReport](stream, data)
}

func RecordTestSuites(c *Client, ctx context.Context, data []*pb.TestSuite) error {
	stream, err := c.client.RecordTestSuites(ctx)
	if err != nil {
		return err
	}

	return send[pb.TestSuite](stream, data)
}

func RecordTestCases(c *Client, ctx context.Context, data []*pb.TestCase) error {
	stream, err := c.client.RecordTestCases(ctx)
	if err != nil {
		return err
	}

	return send[pb.TestCase](stream, data)
}

func RecordLogEmbeddedMetrics(c *Client, ctx context.Context, data []*pb.LogEmbeddedMetric) error {
	stream, err := c.client.RecordLogEmbeddedMetrics(ctx)
	if err != nil {
		return err
	}

	return send[pb.LogEmbeddedMetric](stream, data)
}

func RecordTraces(c *Client, ctx context.Context, data []*pb.Trace) error {
	stream, err := c.client.RecordTraces(ctx)
	if err != nil {
		return err
	}

	return send[pb.Trace](stream, data)
}
