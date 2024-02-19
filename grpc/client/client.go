package client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/cluttrdev/gitlab-exporter/protobuf/servicepb"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type EndpointConfig struct {
	Address string
	Options []grpc.DialOption
}

type Client struct {
	conn   grpc.ClientConnInterface
	client servicepb.GitLabExporterClient
}

func NewCLient(cfg EndpointConfig) (*Client, error) {
	conn, err := grpc.Dial(cfg.Address, cfg.Options...)
	if err != nil {
		return nil, err
	}

	client := servicepb.NewGitLabExporterClient(conn)

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

	m := new(servicepb.RecordSummary)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return nil
}

func RecordPipelines(c *Client, ctx context.Context, data []*typespb.Pipeline) error {
	stream, err := c.client.RecordPipelines(ctx)
	if err != nil {
		return err
	}

	return send[typespb.Pipeline](stream, data)
}

func RecordJobs(c *Client, ctx context.Context, data []*typespb.Job) error {
	stream, err := c.client.RecordJobs(ctx)
	if err != nil {
		return err
	}

	return send[typespb.Job](stream, data)
}

func RecordSections(c *Client, ctx context.Context, data []*typespb.Section) error {
	stream, err := c.client.RecordSections(ctx)
	if err != nil {
		return err
	}

	return send[typespb.Section](stream, data)
}

func RecordBridges(c *Client, ctx context.Context, data []*typespb.Bridge) error {
	stream, err := c.client.RecordBridges(ctx)
	if err != nil {
		return err
	}

	return send[typespb.Bridge](stream, data)
}

func RecordTestReports(c *Client, ctx context.Context, data []*typespb.TestReport) error {
	stream, err := c.client.RecordTestReports(ctx)
	if err != nil {
		return err
	}

	return send[typespb.TestReport](stream, data)
}

func RecordTestSuites(c *Client, ctx context.Context, data []*typespb.TestSuite) error {
	stream, err := c.client.RecordTestSuites(ctx)
	if err != nil {
		return err
	}

	return send[typespb.TestSuite](stream, data)
}

func RecordTestCases(c *Client, ctx context.Context, data []*typespb.TestCase) error {
	stream, err := c.client.RecordTestCases(ctx)
	if err != nil {
		return err
	}

	return send[typespb.TestCase](stream, data)
}

func RecordMetrics(c *Client, ctx context.Context, data []*typespb.Metric) error {
	stream, err := c.client.RecordMetrics(ctx)
	if err != nil {
		return err
	}

	return send[typespb.Metric](stream, data)
}

func RecordTraces(c *Client, ctx context.Context, data []*typespb.Trace) error {
	stream, err := c.client.RecordTraces(ctx)
	if err != nil {
		return err
	}

	return send[typespb.Trace](stream, data)
}
