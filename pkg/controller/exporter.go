package controller

import (
	"context"
	"errors"

	grpc_client "github.com/cluttrdev/gitlab-exporter/grpc/client"
	pb "github.com/cluttrdev/gitlab-exporter/grpc/exporterpb"
	"github.com/cluttrdev/gitlab-exporter/internal/models"
)

type Exporter struct {
	clients []*grpc_client.Client
}

func NewExporter(endpoints []grpc_client.EndpointConfig) (*Exporter, error) {
	var clients []*grpc_client.Client
	for _, cfg := range endpoints {
		c, err := grpc_client.NewCLient(cfg)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}

	return &Exporter{
		clients: clients,
	}, nil
}

func (e *Exporter) RecordPipelines(ctx context.Context, data []*pb.Pipeline) error {
	var errs error
	for _, client := range e.clients {
		if err := client.RecordPipelines(ctx, data); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (e *Exporter) RecordJobs(ctx context.Context, data []*pb.Job) error {
	var errs error
	for _, client := range e.clients {
		if err := client.RecordJobs(ctx, data); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (e *Exporter) RecordSections(ctx context.Context, data []*pb.Section) error {
	var errs error
	for _, client := range e.clients {
		if err := client.RecordSections(ctx, data); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (e *Exporter) RecordBridges(ctx context.Context, data []*pb.Bridge) error {
	var errs error
	for _, client := range e.clients {
		if err := client.RecordBridges(ctx, data); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (e *Exporter) RecordTestReports(ctx context.Context, reports []*pb.TestReport) error {
	var errs error
	for _, client := range e.clients {
		if err := client.RecordTestReports(ctx, reports); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (e *Exporter) RecordTestSuites(ctx context.Context, suites []*pb.TestSuite) error {
	var errs error
	for _, client := range e.clients {
		if err := client.RecordTestSuites(ctx, suites); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (e *Exporter) RecordTestCases(ctx context.Context, cases []*pb.TestCase) error {
	var errs error
	for _, client := range e.clients {
		if err := client.RecordTestCases(ctx, cases); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (e *Exporter) RecordTraces(ctx context.Context, traces []*pb.Trace) error {
	var errs error
	for _, client := range e.clients {
		if err := client.RecordTraces(ctx, traces); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (e *Exporter) RecordPipelineHierarchy(ctx context.Context, ph *models.PipelineHierarchy) error {
	data, err := flattenPipelineHierarchy(ph)
	if err != nil {
		return err
	}
	var errs error
	for _, client := range e.clients {
		if err := client.RecordPipelines(ctx, data.Pipelines); err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		if err := client.RecordJobs(ctx, data.Jobs); err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		if err := client.RecordSections(ctx, data.Sections); err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		if err := client.RecordBridges(ctx, data.Bridges); err != nil {
			errs = errors.Join(errs, err)
			continue
		}
	}

	return errs
}

type pipelineData struct {
	Pipelines []*pb.Pipeline
	Jobs      []*pb.Job
	Sections  []*pb.Section
	Bridges   []*pb.Bridge
}

func flattenPipelineHierarchy(ph *models.PipelineHierarchy) (pipelineData, error) {
	var data pipelineData

	data.Pipelines = append(data.Pipelines, ph.Pipeline)
	data.Jobs = append(data.Jobs, ph.Jobs...)
	data.Sections = append(data.Sections, ph.Sections...)
	data.Bridges = append(data.Bridges, ph.Bridges...)

	for _, dph := range ph.DownstreamPipelines {
		d, err := flattenPipelineHierarchy(dph)
		if err != nil {
			return data, err
		}

		data.Pipelines = append(data.Pipelines, d.Pipelines...)
		data.Jobs = append(data.Jobs, d.Jobs...)
		data.Sections = append(data.Sections, d.Sections...)
		data.Bridges = append(data.Bridges, d.Bridges...)
	}

	return data, nil
}
