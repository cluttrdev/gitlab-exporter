package exporter

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	grpc_client "github.com/cluttrdev/gitlab-exporter/grpc/client"
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
		c, err := grpc_client.NewCLient(cfg.Address, cfg.Options...)
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

type recordFunc[T proto.Message] func(client *grpc_client.Client, ctx context.Context, data []T) error

func export[T proto.Message](exporter *Exporter, ctx context.Context, data []T, record recordFunc[T]) error {
	if len(data) == 0 {
		return nil
	}

	const maxChunkSize int = 2 * 1024 * 1024 // 2 MiB
	var errs error
	var i, j int
	for i < len(data) {
		j = rangeEndForSize(data, maxChunkSize, i)
		if !(i < j) {
			errs = errors.Join(errs, fmt.Errorf("empty range"))
			break
		}

		for _, client := range exporter.clients {
			if err := record(client, ctx, data[i:j]); err != nil {
				errs = errors.Join(errs, err)
			}
		}

		i = j
	}
	return errs
}

func rangeEndForSize[T proto.Message](data []T, size int, start int) int {
	if start < 0 {
		return 0
	}

	var end int = start
	var s int = 0
	for end < len(data) && s < size {
		s += proto.Size(data[end])
		end++
	}

	return end
}

func (e *Exporter) ExportCommits(ctx context.Context, data []*typespb.Commit) error {
	return export[*typespb.Commit](e, ctx, data, grpc_client.RecordCommits)
}

func (e *Exporter) ExportJobs(ctx context.Context, data []*typespb.Job) error {
	return export[*typespb.Job](e, ctx, data, grpc_client.RecordJobs)
}

func (e *Exporter) ExportMergeRequests(ctx context.Context, data []*typespb.MergeRequest) error {
	return export[*typespb.MergeRequest](e, ctx, data, grpc_client.RecordMergeRequests)
}

func (e *Exporter) ExportMergeRequestNoteEvents(ctx context.Context, data []*typespb.MergeRequestNoteEvent) error {
	return export[*typespb.MergeRequestNoteEvent](e, ctx, data, grpc_client.RecordMergeRequestNoteEvents)
}

func (e *Exporter) ExportMetrics(ctx context.Context, data []*typespb.Metric) error {
	return export[*typespb.Metric](e, ctx, data, grpc_client.RecordMetrics)
}

func (e *Exporter) ExportPipelines(ctx context.Context, data []*typespb.Pipeline) error {
	return export[*typespb.Pipeline](e, ctx, data, grpc_client.RecordPipelines)
}

func (e *Exporter) ExportProjects(ctx context.Context, data []*typespb.Project) error {
	return export[*typespb.Project](e, ctx, data, grpc_client.RecordProjects)
}

func (e *Exporter) ExportSections(ctx context.Context, data []*typespb.Section) error {
	return export[*typespb.Section](e, ctx, data, grpc_client.RecordSections)
}

func (e *Exporter) ExportTestCases(ctx context.Context, data []*typespb.TestCase) error {
	return export[*typespb.TestCase](e, ctx, data, grpc_client.RecordTestCases)
}

func (e *Exporter) ExportTestReports(ctx context.Context, data []*typespb.TestReport) error {
	return export[*typespb.TestReport](e, ctx, data, grpc_client.RecordTestReports)
}

func (e *Exporter) ExportTestSuites(ctx context.Context, data []*typespb.TestSuite) error {
	return export[*typespb.TestSuite](e, ctx, data, grpc_client.RecordTestSuites)
}

func (e *Exporter) ExportTraces(ctx context.Context, data []*typespb.Trace) error {
	return export[*typespb.Trace](e, ctx, data, grpc_client.RecordTraces)
}
