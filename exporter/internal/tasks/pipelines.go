package tasks

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab/rest"
	"go.cluttr.dev/gitlab-exporter/internal/logql"
	"go.cluttr.dev/gitlab-exporter/internal/types"
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
	var (
		sections   []types.Section
		metrics    []types.Metric
		properties = make(map[int64][]types.JobLogProperty)
	)

	type result struct {
		jobId int64

		sections   []types.Section
		metrics    []types.Metric
		properties []types.JobLogProperty

		err error
	}

	var (
		wg      sync.WaitGroup
		results = make(chan result, len(jobs))
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, job := range jobs {
			if err := glab.Acquire(ctx, 1); err != nil {
				slog.Error("failed to acquire gitlab client", "error", err)
				break
			}
			wg.Add(1)
			go func() {
				defer glab.Release(1)
				defer wg.Done()

				var opt = FetchProjectJobLogDataOptions{
					Queries: opts.ProjectJobLogQueries[job.Pipeline.Project.Id],
				}

				var r result
				r.jobId = job.Id
				r.sections, r.metrics, r.properties, r.err = FetchProjectJobLogData(ctx, glab, job, opt)
				results <- r
			}()
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	var errs error
loop:
	for {
		select {
		case <-done:
			break loop
		case r := <-results:
			if r.err != nil {
				errs = errors.Join(errs, r.err)
			} else {
				sections = append(sections, r.sections...)
				metrics = append(metrics, r.metrics...)
				properties[r.jobId] = append(properties[r.jobId], r.properties...)
			}
		}
	}

	return sections, metrics, properties, errs
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
	if err != nil || log == nil {
		return nil, nil, nil, err
	}

	logData, err := rest.ParseJobLog(log)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("parse job log: %w", err)
	}

	var jobFinishedAtUnix int64
	if job.FinishedAt != nil {
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
			StartedAt:  gitlab.Ptr(time.Unix(secdat.Start, 0)),
			FinishedAt: gitlab.Ptr(time.Unix(secdat.End, 0)),
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

	_, _ = log.Seek(0, 0)
	logqlMetrics, err := queryJobLogQLMetrics(log, opts.Queries, jobRef, int64(len(logData.Metrics)))
	if err != nil {
		return sections, metrics, properties, fmt.Errorf("query job logql metrics: %w", err)
	}
	if job.FinishedAt != nil {
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
		})
	}

	return metrics, nil
}

func FetchProjectsPipelinesTestReports(ctx context.Context, glab *gitlab.Client, pipelines []types.Pipeline) ([]types.TestReport, []types.TestSuite, []types.TestCase, error) {
	var (
		testReports []types.TestReport
		testSuites  []types.TestSuite
		testCases   []types.TestCase
	)

	type result struct {
		testReport types.TestReport
		testSuites []types.TestSuite
		testCases  []types.TestCase

		err error
	}

	var (
		wg      sync.WaitGroup
		results = make(chan result, len(pipelines))
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, p := range pipelines {
			if err := glab.Acquire(ctx, 1); err != nil {
				slog.Error("failed to acquire gitlab client", "error", err)
				break
			}
			wg.Add(1)
			go func(pipeline types.Pipeline) {
				defer glab.Release(1)
				defer wg.Done()

				var (
					r   result
					err error
				)

				report, summary, err := glab.Rest.GetPipelineTestReport(ctx, pipeline.Project.Id, pipeline.Id)
				if err != nil {
					r.err = fmt.Errorf("get project pipeline test report: %w", err)
					results <- r
					return
				}

				r.testReport, r.testSuites, r.testCases, err = rest.ConvertTestReport(report, summary, pipeline)
				if err != nil {
					r.err = fmt.Errorf("convert project pipeline test report: %w", err)
				}
				results <- r
			}(p)
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	var errs error
loop:
	for {
		select {
		case <-done:
			break loop
		case r := <-results:
			if r.err != nil {
				errs = errors.Join(errs, r.err)
			} else {
				testReports = append(testReports, r.testReport)
				testSuites = append(testSuites, r.testSuites...)
				testCases = append(testCases, r.testCases...)
			}
		}
	}

	return testReports, testSuites, testCases, errs
}
