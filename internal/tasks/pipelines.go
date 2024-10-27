package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/graphql"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/rest"
	"github.com/cluttrdev/gitlab-exporter/internal/types"
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
	if err != nil {
		return nil, fmt.Errorf("get pipeline fields: %w", err)
	}

	pipelines := make([]types.Pipeline, 0, len(pipelinesFields))
	for _, pf := range pipelinesFields {
		p, err := graphql.ConvertPipeline(pf)
		if err != nil {
			return nil, fmt.Errorf("convert pipeline fields: %w", err)
		}
		pipelines = append(pipelines, p)
	}

	return pipelines, nil
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
	if err != nil {
		return nil, fmt.Errorf("get projects pipelines jobs: %w", err)
	}

	jobs := make([]types.Job, 0, len(jobsFields))
	for _, jf := range jobsFields {
		j, err := graphql.ConvertJob(jf)
		if err != nil {
			return nil, fmt.Errorf("convert job fields: %w", err)
		}
		jobs = append(jobs, j)
	}

	return jobs, nil
}

func FetchProjectsJobsLogData(ctx context.Context, glab *gitlab.Client, jobs []types.Job) ([]types.Section, []types.Metric, error) {
	var (
		sections []types.Section
		metrics  []types.Metric
	)

	type result struct {
		sections []types.Section
		metrics  []types.Metric

		err error
	}

	var (
		wg      sync.WaitGroup
		results = make(chan result, len(jobs))
	)
	for _, job := range jobs {
		if err := glab.Acquire(ctx, 1); err != nil {
			slog.Error("failed to acquire gitlab client", "error", err)
			break
		}
		wg.Add(1)
		go func() {
			defer glab.Release(1)
			defer wg.Done()

			var r result
			r.sections, r.metrics, r.err = FetchProjectJobLogData(ctx, glab, job)
			results <- r
		}()
	}

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
			}
		}
	}

	return sections, metrics, errs
}

func FetchProjectJobLogData(ctx context.Context, glab *gitlab.Client, job types.Job) ([]types.Section, []types.Metric, error) {
	var (
		sections []types.Section
		metrics  []types.Metric
	)

	logData, err := glab.Rest.GetJobLogData(ctx, job.Pipeline.Project.Id, job.Id)
	if err != nil {
		return nil, nil, err
	}

	for secnum, secdat := range logData.Sections {
		sections = append(sections, types.Section{
			Id: job.Id*1000 + int64(secnum),
			Job: types.JobReference{
				Id:       job.Id,
				Name:     job.Name,
				Pipeline: job.Pipeline,
			},

			Name:       secdat.Name,
			StartedAt:  gitlab.Ptr(time.Unix(secdat.Start, 0)),
			FinishedAt: gitlab.Ptr(time.Unix(secdat.End, 0)),
			Duration:   time.Duration((secdat.End - secdat.Start) * int64(time.Second)),
		})
	}

	for iid, m := range logData.Metrics {
		metrics = append(metrics, types.Metric{
			Id:  fmt.Sprintf("%d-%d", job.Id, iid+1),
			Iid: int64(iid + 1),
			Job: types.JobReference{
				Id:       job.Id,
				Name:     job.Name,
				Pipeline: job.Pipeline,
			},

			Name:      m.Name,
			Labels:    m.Labels,
			Value:     m.Value,
			Timestamp: m.Timestamp,
		})
	}

	return sections, metrics, nil
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
	for _, pipeline := range pipelines {
		if err := glab.Acquire(ctx, 1); err != nil {
			slog.Error("failed to acquire gitlab client", "error", err)
			break
		}
		wg.Add(1)
		go func() {
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
		}()
	}

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

	return testReports, testSuites, testCases, nil
}
