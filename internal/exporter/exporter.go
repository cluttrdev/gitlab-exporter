package exporter

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

func New(endpoints []grpc_client.EndpointConfig) (*Exporter, error) {
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

type recordFunc[T any] func(client *grpc_client.Client, ctx context.Context, data []*T) error

func export[T any](exporter *Exporter, ctx context.Context, data []*T, record recordFunc[T]) error {
	if len(data) == 0 {
		return nil
	}
	var errs error
	for _, client := range exporter.clients {
		if err := record(client, ctx, data); err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func (e *Exporter) ExportPipelines(ctx context.Context, data []*pb.Pipeline) error {
	return export[pb.Pipeline](e, ctx, data, grpc_client.RecordPipelines)
}

func (e *Exporter) ExportJobs(ctx context.Context, data []*pb.Job) error {
	return export[pb.Job](e, ctx, data, grpc_client.RecordJobs)
}

func (e *Exporter) ExportSections(ctx context.Context, data []*pb.Section) error {
	return export[pb.Section](e, ctx, data, grpc_client.RecordSections)
}

func (e *Exporter) ExportBridges(ctx context.Context, data []*pb.Bridge) error {
	return export[pb.Bridge](e, ctx, data, grpc_client.RecordBridges)
}

func (e *Exporter) ExportTestReports(ctx context.Context, data []*pb.TestReport) error {
	return export[pb.TestReport](e, ctx, data, grpc_client.RecordTestReports)
}

func (e *Exporter) ExportTestSuites(ctx context.Context, data []*pb.TestSuite) error {
	return export[pb.TestSuite](e, ctx, data, grpc_client.RecordTestSuites)
}

func (e *Exporter) ExportTestCases(ctx context.Context, data []*pb.TestCase) error {
	return export[pb.TestCase](e, ctx, data, grpc_client.RecordTestCases)
}

func (e *Exporter) ExportMetrics(ctx context.Context, data []*pb.Metric) error {
	return export[pb.Metric](e, ctx, data, grpc_client.RecordMetrics)
}

func (e *Exporter) ExportTraces(ctx context.Context, data []*pb.Trace) error {
	return export[pb.Trace](e, ctx, data, grpc_client.RecordTraces)
}

func (e *Exporter) ExportPipelineHierarchy(ctx context.Context, ph *models.PipelineHierarchy) error {
	data, err := flattenPipelineHierarchy(ph)
	if err != nil {
		return err
	}

	var errs error
	for _, client := range e.clients {
		if err := grpc_client.RecordPipelines(client, ctx, data.Pipelines); err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		if err := grpc_client.RecordJobs(client, ctx, data.Jobs); err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		if err := grpc_client.RecordSections(client, ctx, data.Sections); err != nil {
			errs = errors.Join(errs, err)
			continue
		}
		if err := grpc_client.RecordBridges(client, ctx, data.Bridges); err != nil {
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
