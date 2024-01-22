package client

import (
	"context"
	"errors"

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

func (c *Client) RecordPipelines(ctx context.Context, ps []*pb.Pipeline) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordPipelines(ctx)
	if err != nil {
		return err
	}

	for _, p := range ps {
		err := stream.Send(p)
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordJobs(ctx context.Context, jobs []*pb.Job) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordJobs(ctx)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		err := stream.Send(job)
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordSections(ctx context.Context, sections []*pb.Section) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordSections(ctx)
	if err != nil {
		return err
	}

	for _, section := range sections {
		err := stream.Send(section)
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordBridges(ctx context.Context, bridges []*pb.Bridge) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordBridges(ctx)
	if err != nil {
		return err
	}

	for _, bridge := range bridges {
		err := stream.Send(bridge)
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordTestReports(ctx context.Context, reports []*pb.TestReport) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordTestReports(ctx)
	if err != nil {
		return err
	}

	for _, tr := range reports {
		err := stream.Send(tr)
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordTestSuites(ctx context.Context, suites []*pb.TestSuite) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordTestSuites(ctx)
	if err != nil {
		return err
	}

	for _, ts := range suites {
		err := stream.Send(ts)
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordTestCases(ctx context.Context, cases []*pb.TestCase) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordTestCases(ctx)
	if err != nil {
		return err
	}

	for _, tc := range cases {
		err := stream.Send(tc)
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordLogEmbeddedMetrics(ctx context.Context, metrics []*pb.LogEmbeddedMetric) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordLogEmbeddedMetrics(ctx)
	if err != nil {
		return err
	}

	for _, m := range metrics {
		err := stream.Send(m)
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordTraces(ctx context.Context, traces []*pb.Trace) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordTraces(ctx)
	if err != nil {
		return err
	}

	for _, t := range traces {
		err := stream.Send(t)
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}
