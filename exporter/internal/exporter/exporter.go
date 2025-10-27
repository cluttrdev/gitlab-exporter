package exporter

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	tracepb_v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/prometheus/client_golang/prometheus"
	grpc_client "go.cluttr.dev/gitlab-exporter/grpc/client"
	"go.cluttr.dev/gitlab-exporter/internal/exporter/messages"
	"go.cluttr.dev/gitlab-exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
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

type convertFunc[T any, M proto.Message] func(data T) M

func convert[T any, M proto.Message](data []T, cfun convertFunc[T, M]) []M {
	msgs := make([]M, 0, len(data))
	for _, d := range data {
		msgs = append(msgs, cfun(d))
	}
	return msgs
}

type recordFunc[T proto.Message] func(client *grpc_client.Client, ctx context.Context, data []T) error

func export[T proto.Message](exp *Exporter, ctx context.Context, data []T, record recordFunc[T]) error {
	if len(data) == 0 {
		return nil
	}

	// split data into batches to keep max message size
	batches, err := createBatches(data)
	if err != nil {
		return err
	}

	// for each client, export batches concurrently in an error group
	var wg sync.WaitGroup
	errChan := make(chan error)
	for _, client := range exp.clients {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// set up context aware error group and limit number of
			// active goroutines to send too much data at once
			eg, ctx := errgroup.WithContext(ctx)
			eg.SetLimit(100) // max 200 MiB is sent at once

			for _, batch := range batches {
				eg.Go(func() error {
					return record(client, ctx, batch)
				})
			}
			if err := eg.Wait(); err != nil {
				errChan <- err
			}
		}()
	}

	// wait for all client goroutines to finish
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	// collect errors from client goroutines
	var errs error
loop:
	for {
		select {
		case <-done:
			break loop
		case err := <-errChan:
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

func createBatches[T proto.Message](data []T) ([][]T, error) {
	const maxChunkSize int = 2 * 1024 * 1024 // 2 MiB
	var batches [][]T
	var i, j int
	for i < len(data) {
		j = rangeEndForSize(data, maxChunkSize, i)
		if !(i < j) {
			return nil, fmt.Errorf("empty range")
		}

		batches = append(batches, data[i:j])
		i = j
	}
	return batches, nil
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

func (e *Exporter) ExportCoverageReports(ctx context.Context, data []types.CoverageReport) error {
	msgs := convert(data, messages.NewCoverageReport)
	return export(e, ctx, msgs, grpc_client.RecordCoverageReports)
}

func (e *Exporter) ExportCoveragePackages(ctx context.Context, data []types.CoveragePackage) error {
	msgs := convert(data, messages.NewCoveragePackage)
	return export(e, ctx, msgs, grpc_client.RecordCoveragePackages)
}

func (e *Exporter) ExportCoverageClasses(ctx context.Context, data []types.CoverageClass) error {
	msgs := convert(data, messages.NewCoverageClass)
	return export(e, ctx, msgs, grpc_client.RecordCoverageClasses)
}

func (e *Exporter) ExportCoverageMethods(ctx context.Context, data []types.CoverageMethod) error {
	msgs := convert(data, messages.NewCoverageMethod)
	return export(e, ctx, msgs, grpc_client.RecordCoverageMethods)
}

func (e *Exporter) ExportDeployments(ctx context.Context, data []types.Deployment) error {
	msgs := convert(data, messages.NewDeployment)
	return export(e, ctx, msgs, grpc_client.RecordDeployments)
}

func (e *Exporter) ExportIssues(ctx context.Context, data []types.Issue) error {
	msgs := convert(data, messages.NewIssue)
	return export(e, ctx, msgs, grpc_client.RecordIssues)
}

func (e *Exporter) ExportJobs(ctx context.Context, data []types.Job) error {
	msgs := convert(data, messages.NewJob)
	return export(e, ctx, msgs, grpc_client.RecordJobs)
}

func (e *Exporter) ExportMergeRequests(ctx context.Context, data []types.MergeRequest) error {
	msgs := convert(data, messages.NewMergeRequest)
	return export(e, ctx, msgs, grpc_client.RecordMergeRequests)
}

func (e *Exporter) ExportMergeRequestNoteEvents(ctx context.Context, data []types.MergeRequestNoteEvent) error {
	msgs := convert(data, messages.NewMergeRequestNoteEvent)
	return export(e, ctx, msgs, grpc_client.RecordMergeRequestNoteEvents)
}

func (e *Exporter) ExportMetrics(ctx context.Context, data []types.Metric) error {
	msgs := convert(data, messages.NewMetric)
	return export(e, ctx, msgs, grpc_client.RecordMetrics)
}

func (e *Exporter) ExportPipelines(ctx context.Context, data []types.Pipeline) error {
	msgs := convert(data, messages.NewPipeline)
	return export(e, ctx, msgs, grpc_client.RecordPipelines)
}

func (e *Exporter) ExportProjects(ctx context.Context, data []types.Project) error {
	msgs := convert(data, messages.NewProject)
	return export(e, ctx, msgs, grpc_client.RecordProjects)
}

func (e *Exporter) ExportRunners(ctx context.Context, data []types.Runner, fetchedAt time.Time) error {
	msgs := convert(data, messages.NewRunner)
	record := func(client *grpc_client.Client, ctx context.Context, data []*typespb.Runner) error {
		return grpc_client.RecordRunners(client, ctx, data, fetchedAt)
	}
	return export(e, ctx, msgs, record)
}

func (e *Exporter) ExportSections(ctx context.Context, data []types.Section) error {
	msgs := convert(data, messages.NewSection)
	return export(e, ctx, msgs, grpc_client.RecordSections)
}

func (e *Exporter) ExportTestCases(ctx context.Context, data []types.TestCase) error {
	msgs := convert(data, messages.NewTestCase)
	return export(e, ctx, msgs, grpc_client.RecordTestCases)
}

func (e *Exporter) ExportTestReports(ctx context.Context, data []types.TestReport) error {
	msgs := convert(data, messages.NewTestReport)
	return export(e, ctx, msgs, grpc_client.RecordTestReports)
}

func (e *Exporter) ExportTestSuites(ctx context.Context, data []types.TestSuite) error {
	msgs := convert(data, messages.NewTestSuite)
	return export(e, ctx, msgs, grpc_client.RecordTestSuites)
}

func (e *Exporter) ExportPipelineSpans(ctx context.Context, data []types.Pipeline) error {
	spans := convert(data, messages.NewPipelineSpan)

	batches, err := createBatches(spans)
	if err != nil {
		return err
	}

	var msgs []*typespb.Trace
	for _, batch := range batches {
		msgs = append(msgs, &typespb.Trace{
			Data: &tracepb_v1.TracesData{
				ResourceSpans: []*tracepb_v1.ResourceSpans{
					messages.NewResourceSpan(map[string]string{
						"service.name": "gitlab_ci.pipeline",
					}, batch),
				},
			},
		})
	}

	return export(e, ctx, msgs, grpc_client.RecordTraces)
}

func (e *Exporter) ExportJobSpans(ctx context.Context, data []types.Job) error {
	spans := convert(data, messages.NewJobSpan)

	batches, err := createBatches(spans)
	if err != nil {
		return err
	}

	var msgs []*typespb.Trace
	for _, batch := range batches {
		msgs = append(msgs, &typespb.Trace{
			Data: &tracepb_v1.TracesData{
				ResourceSpans: []*tracepb_v1.ResourceSpans{
					messages.NewResourceSpan(map[string]string{
						"service.name": "gitlab_ci.job",
					}, batch),
				},
			},
		})
	}

	return export(e, ctx, msgs, grpc_client.RecordTraces)
}

func (e *Exporter) ExportSectionSpans(ctx context.Context, data []types.Section) error {
	spans := convert(data, messages.NewSectionSpan)

	batches, err := createBatches(spans)
	if err != nil {
		return err
	}

	var msgs []*typespb.Trace
	for _, batch := range batches {
		msgs = append(msgs, &typespb.Trace{
			Data: &tracepb_v1.TracesData{
				ResourceSpans: []*tracepb_v1.ResourceSpans{
					messages.NewResourceSpan(map[string]string{
						"service.name": "gitlab_ci.section",
					}, batch),
				},
			},
		})
	}

	return export(e, ctx, msgs, grpc_client.RecordTraces)
}
