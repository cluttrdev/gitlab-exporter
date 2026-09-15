package tasks

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/rest"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/logql"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/types"
)

func FetchProjectsPipelines(ctx context.Context, glab *gitlab.Client, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) ([]types.Pipeline, error) {
	gids := make([]string, 0, len(projectIds))
	for _, id := range projectIds {
		gids = append(gids, graphql.FormatId(id, graphql.GlobalIdProjectPrefix))
	}

	opts := graphql.GetPipelinesOptions{
		UpdatedAfter:  updatedAfter,
		UpdatedBefore: updatedBefore,
	}

	pipelinesFields, err := glab.GraphQL.GetProjectsPipelines(ctx, gids, opts)
	if errors.Is(err, context.Canceled) {
		return nil, err
	} else if err != nil {
		err = fmt.Errorf("get pipeline fields: %w", err)
	}

	pipelines := make([]types.Pipeline, 0, len(pipelinesFields))
	for _, pf := range pipelinesFields {
		p, err := graphql.ConvertPipeline(pf)
		if err != nil {
			slog.Error("error converting pipeline fields",
				slog.String("id", pf.Id),
				slog.String("projectId", pf.Project.Id),
				slog.String("error", err.Error()),
			)
			continue
		}
		pipelines = append(pipelines, p)
	}

	return pipelines, err
}

func FetchProjectsPipelinesJobs(ctx context.Context, glab *gitlab.Client, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) ([]types.Job, error) {
	gids := make([]string, 0, len(projectIds))
	for _, id := range projectIds {
		gids = append(gids, graphql.FormatId(id, graphql.GlobalIdProjectPrefix))
	}

	opts := graphql.GetPipelinesOptions{
		UpdatedAfter:  updatedAfter,
		UpdatedBefore: updatedBefore,
	}

	jobsFields, err := glab.GraphQL.GetProjectsPipelinesJobs(ctx, gids, opts)
	if errors.Is(err, context.Canceled) {
		return nil, err
	} else if err != nil {
		err = fmt.Errorf("get projects pipelines jobs: %w", err)
	}

	jobs := make([]types.Job, 0, len(jobsFields))
	for _, jf := range jobsFields {
		j, err := graphql.ConvertJob(jf)
		if err != nil {
			var jfId string
			if jf.Id != nil {
				jfId = *jf.Id
			}
			slog.Error("error converting job fields",
				slog.String("id", jfId),
				slog.String("pipelineId", jf.Pipeline.Id),
				slog.String("projectId", jf.Project.Id),
				slog.String("error", err.Error()),
			)
			continue
		}
		jobs = append(jobs, j)
	}

	return jobs, err
}

type FetchProjectsJobsLogDataOptions struct {
	ProjectJobLogQueries map[int64][]logql.MetricQuery
}

func FetchProjectsJobsLogData(ctx context.Context, glab *gitlab.Client, jobs []types.Job, opts FetchProjectsJobsLogDataOptions) ([]types.Section, []types.Metric, map[int64][]types.JobLogProperty, error) {
	// jobLogBundle groups one job's fan-out results.
	// jobId is carried explicitly so the reduce step can rebuild the per-job
	// properties map, including for errored jobs.
	type jobLogBundle struct {
		jobId      int64
		sections   []types.Section
		metrics    []types.Metric
		properties []types.JobLogProperty
	}

	bundles, err := fanOut(ctx, glab, jobs,
		func(ctx context.Context, glab *gitlab.Client, job types.Job) (jobLogBundle, error) {
			opt := FetchProjectJobLogDataOptions{
				Queries: opts.ProjectJobLogQueries[job.Pipeline.Project.Id],
			}

			sections, metrics, properties, err := FetchProjectJobLogData(ctx, glab, job, opt)
			return jobLogBundle{
				jobId:      job.Id,
				sections:   sections,
				metrics:    metrics,
				properties: properties,
			}, err
		},
		// collectOnError: FetchProjectJobLogData can return populated
		// sections/metrics/properties *with* an error (rewind + logql-query
		// paths).
		fanOutOptions{CollectOnError: true},
	)

	var (
		sections   []types.Section
		metrics    []types.Metric
		properties = make(map[int64][]types.JobLogProperty)
	)
	for _, b := range bundles {
		sections = append(sections, b.sections...)
		metrics = append(metrics, b.metrics...)
		properties[b.jobId] = append(properties[b.jobId], b.properties...)
	}

	return sections, metrics, properties, err
}

type FetchProjectJobLogDataOptions struct {
	Queries []logql.MetricQuery
}

func FetchProjectJobLogData(ctx context.Context, glab *gitlab.Client, job types.Job, opts FetchProjectJobLogDataOptions) ([]types.Section, []types.Metric, []types.JobLogProperty, error) {
	var (
		sections   []types.Section
		metrics    []types.Metric
		properties []types.JobLogProperty
	)

	jobRef := types.JobReference{
		Id:       job.Id,
		Name:     job.Name,
		Pipeline: job.Pipeline,
	}

	log, err := glab.Rest.GetJobLog(ctx, job.Pipeline.Project.Id, job.Id)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get job log: %w", err)
	} else if log == nil {
		return nil, nil, nil, nil
	}

	logData, err := rest.ParseJobLog(log)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse job log: %w", err)
	}

	var jobFinishedAtUnix int64 = 0
	if job.FinishedAt != nil && !job.FinishedAt.IsZero() {
		jobFinishedAtUnix = job.FinishedAt.Unix()
	}
	for secnum, secdat := range logData.Sections {
		if secdat.Start > 0 && secdat.End == 0 { // section started but not finished
			secdat.End = max(secdat.Start, jobFinishedAtUnix)
		}

		section := types.Section{
			Id:  jobRef.Id*1000 + int64(secnum),
			Job: jobRef,

			Name:       secdat.Name,
			StartedAt:  gitlab.Ptr(convertSectionTimestamp(secdat.Start)),
			FinishedAt: gitlab.Ptr(convertSectionTimestamp(secdat.End)),
			Duration:   time.Duration((secdat.End - secdat.Start) * int64(time.Second)),
		}

		sections = append(sections, section)
	}

	for iid, m := range logData.Metrics {
		metrics = append(metrics, types.Metric{
			Id:  fmt.Sprintf("%d-%d", jobRef.Id, iid+1),
			Iid: int64(iid + 1),
			Job: jobRef,

			Name:      m.Name,
			Labels:    m.Labels,
			Value:     m.Value,
			Timestamp: m.Timestamp,
		})
	}

	for _, p := range logData.Properties {
		properties = append(properties, types.JobLogProperty{
			Name:  p.Name,
			Value: p.Value,
		})
	}

	if _, err := log.Seek(0, 0); err != nil {
		return sections, metrics, properties, fmt.Errorf("rewind job log: %w", err)
	}
	logqlMetrics, err := queryJobLogQLMetrics(log, opts.Queries, jobRef, int64(len(logData.Metrics)))
	if err != nil {
		return sections, metrics, properties, fmt.Errorf("query job logql metrics: %w", err)
	}
	if job.FinishedAt != nil && !job.FinishedAt.IsZero() {
		for i := range len(logqlMetrics) {
			logqlMetrics[i].Timestamp = job.FinishedAt.UnixMilli()
		}
	}
	for _, m := range logqlMetrics {
		if m.Value > 0 {
			metrics = append(metrics, m)
		}
	}

	return sections, metrics, properties, nil
}

// convertSectionTimestamp converts a timestamp from a job log section to time.Time.
//
// The Unix timestamp for the job log sections is generated using the awk command:
//
//	$(awk 'BEGIN{srand(); print srand()}')
//
// When srand() is called without arguments, it seeds the random number generator
// using the current time and returns the previous seed value.
// So in `{srand(); print srand()}` the first srand() seeds with the current time,
// then the second srand() re-seeds and returns the previous seed.
//
// Depending on the awk implementation used in the jobs shell environment, the
// time unit used for seeding can differ, e.g.
// - GNU awk (gawk) typically uses time() which returns seconds since epoch
// - mawk uses a higher-precision timer for seeding, returning nanoseconds since epoch
func convertSectionTimestamp(ts int64) time.Time {
	const maxReasonableSeconds = 4294967295 // 2106-02-07T06:28:15.000Z, uint32_max

	// Try as seconds first
	if ts < maxReasonableSeconds {
		return time.Unix(ts, 0)
	}

	// Try as milliseconds (13 digits for current era)
	if ms := ts / 1000; ms < maxReasonableSeconds {
		return time.Unix(ms, (ts%1000)*1_000_000)
	}

	// Try as microseconds (16 digits)
	if us := ts / 1_000_000; us < maxReasonableSeconds {
		return time.Unix(us, (ts%1_000_000)*1000)
	}

	// Must be nanoseconds
	return time.Unix(0, ts)
}

func queryJobLogQLMetrics(log *bytes.Reader, queries []logql.MetricQuery, job types.JobReference, startIid int64) ([]types.Metric, error) {
	filters := make([]logql.LineFilter, 0, len(queries))
	for _, query := range queries {
		filters = append(filters, query.LineFilter)
	}

	counts, err := logql.Count(log, filters)
	if err != nil {
		return nil, err
	}

	metrics := make([]types.Metric, 0, len(queries))
	for i, query := range queries {
		iid := startIid + 1 + int64(i)
		metrics = append(metrics, types.Metric{
			Id:  fmt.Sprintf("%d-%d", job.Id, iid),
			Iid: iid,
			Job: job,

			Name:   query.Name,
			Labels: query.LabelAdd,
			Value:  float64(counts[i]),
			// Timestamp: 0,
		})
	}

	return metrics, nil
}

func FetchProjectsPipelinesTestReports(ctx context.Context, glab *gitlab.Client, pipelines []types.Pipeline) ([]types.TestReport, []types.TestSuite, []types.TestCase, error) {
	type testReportBundle struct {
		report types.TestReport
		suites []types.TestSuite
		cases  []types.TestCase
	}

	bundles, err := fanOut(ctx, glab, pipelines,
		func(ctx context.Context, glab *gitlab.Client, pipeline types.Pipeline) (testReportBundle, error) {
			report, summary, err := glab.Rest.GetPipelineTestReport(ctx, pipeline.Project.Id, pipeline.Id)
			if err != nil {
				return testReportBundle{}, fmt.Errorf("get project pipeline test report: %w", err)
			}

			tr, ts, tc, err := rest.ConvertTestReport(report, summary, pipeline)
			if err != nil {
				return testReportBundle{}, fmt.Errorf("convert project pipeline test report: %w", err)
			}

			return testReportBundle{report: tr, suites: ts, cases: tc}, nil
		},
		fanOutOptions{},
	)

	// Reduce: concatenate each pipeline's bundle back into the three flat slices.
	var (
		testReports []types.TestReport
		testSuites  []types.TestSuite
		testCases   []types.TestCase
	)
	for _, b := range bundles {
		testReports = append(testReports, b.report)
		testSuites = append(testSuites, b.suites...)
		testCases = append(testCases, b.cases...)
	}
	return testReports, testSuites, testCases, err
}
