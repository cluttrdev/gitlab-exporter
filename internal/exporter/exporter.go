package exporter

import (
	"context"
	"errors"

	"google.golang.org/grpc"

	grpc_client "github.com/cluttrdev/gitlab-exporter/grpc/client"
	"github.com/cluttrdev/gitlab-exporter/internal/models"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	clients map[string]*grpc_client.Client
}

type EndpointConfig struct {
	Address string
	Options []grpc.DialOption
}

func New(endpoints []EndpointConfig) (*Exporter, error) {
	clients := make(map[string]*grpc_client.Client, len(endpoints))
	for _, cfg := range endpoints {
		c, err := grpc_client.NewCLient(context.Background(), cfg.Address, cfg.Options...)
		if err != nil {
			return nil, err
		}
		clients[cfg.Address] = c
	}

	return &Exporter{
		clients: clients,
	}, nil
}

func (e *Exporter) MetricsCollectorFor(endpoint string) prometheus.Collector {
	c, ok := e.clients[endpoint]
	if !ok {
		return nil
	}
	return c.MetricsCollector()
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

func (e *Exporter) ExportPipelines(ctx context.Context, data []*typespb.Pipeline) error {
	return export[typespb.Pipeline](e, ctx, data, grpc_client.RecordPipelines)
}

func (e *Exporter) ExportJobs(ctx context.Context, data []*typespb.Job) error {
	return export[typespb.Job](e, ctx, data, grpc_client.RecordJobs)
}

func (e *Exporter) ExportSections(ctx context.Context, data []*typespb.Section) error {
	return export[typespb.Section](e, ctx, data, grpc_client.RecordSections)
}

func (e *Exporter) ExportBridges(ctx context.Context, data []*typespb.Bridge) error {
	return export[typespb.Bridge](e, ctx, data, grpc_client.RecordBridges)
}

func (e *Exporter) ExportTestReports(ctx context.Context, data []*typespb.TestReport) error {
	return export[typespb.TestReport](e, ctx, data, grpc_client.RecordTestReports)
}

func (e *Exporter) ExportTestSuites(ctx context.Context, data []*typespb.TestSuite) error {
	return export[typespb.TestSuite](e, ctx, data, grpc_client.RecordTestSuites)
}

func (e *Exporter) ExportTestCases(ctx context.Context, data []*typespb.TestCase) error {
	return export[typespb.TestCase](e, ctx, data, grpc_client.RecordTestCases)
}

func (e *Exporter) ExportMetrics(ctx context.Context, data []*typespb.Metric) error {
	return export[typespb.Metric](e, ctx, data, grpc_client.RecordMetrics)
}

func (e *Exporter) ExportTraces(ctx context.Context, data []*typespb.Trace) error {
	return export[typespb.Trace](e, ctx, data, grpc_client.RecordTraces)
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
	Pipelines []*typespb.Pipeline
	Jobs      []*typespb.Job
	Sections  []*typespb.Section
	Bridges   []*typespb.Bridge
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
