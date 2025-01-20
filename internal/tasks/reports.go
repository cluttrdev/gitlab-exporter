package tasks

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/graphql"
	"github.com/cluttrdev/gitlab-exporter/internal/junitxml"
	"github.com/cluttrdev/gitlab-exporter/internal/types"
)

func FetchProjectsPipelinesJunitReports(ctx context.Context, glab *gitlab.Client, projectPipelines map[string][]string) ([]types.TestReport, []types.TestSuite, []types.TestCase, error) {
	var (
		testReports []types.TestReport
		testSuites  []types.TestSuite
		testCases   []types.TestCase
	)

	type result struct {
		testReports []types.TestReport
		testSuites  []types.TestSuite
		testCases   []types.TestCase

		err error
	}

	var (
		wg      sync.WaitGroup
		results = make(chan result)
	)
	wg.Add(1)
	go func() {
		defer wg.Done()

		for projectPath, pipelineIids := range projectPipelines {
			for _, pipelineIid := range pipelineIids {
				if err := glab.Acquire(ctx, 1); err != nil {
					slog.Error("failed to acquire gitlab client", "error", err)
					continue
				}
				wg.Add(1)
				go func(projectPath string, pipelineIid string) {
					defer glab.Release(1)
					defer wg.Done()

					tr, ts, tc, err := FetchProjectPipelineJunitReports(ctx, glab, projectPath, pipelineIid)

					results <- result{
						testReports: tr,
						testSuites:  ts,
						testCases:   tc,
						err:         err,
					}
				}(projectPath, pipelineIid)
			}
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
				testReports = append(testReports, r.testReports...)
				testSuites = append(testSuites, r.testSuites...)
				testCases = append(testCases, r.testCases...)
			}
		}
	}

	return testReports, testSuites, testCases, errs
}

func FetchProjectPipelineJunitReports(ctx context.Context, glab *gitlab.Client, projectPath string, pipelineIid string) ([]types.TestReport, []types.TestSuite, []types.TestCase, error) {
	var (
		testReports []types.TestReport
		testSuites  []types.TestSuite
		testCases   []types.TestCase
	)

	artifacts, err := glab.GraphQL.GetProjectPipelineJobsArtifacts(ctx, projectPath, pipelineIid)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get project pipeline job artifacts: %w", err)
	}

	for _, artifact := range artifacts {
		if artifact.FileType == nil || *artifact.FileType != graphql.JobArtifactFileTypeJunit {
			continue
		}
		if artifact.DownloadPath == nil || *artifact.DownloadPath == "" {
			continue
		}

		downloadPath := *artifact.DownloadPath
		jobRef, err := graphql.ConvertJobReference(artifact.Job, artifact.Pipeline, artifact.Project)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("convert job reference: %w", err)
		}

		resp, err := glab.HTTP.GetPath(downloadPath)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("download report: %w", err)
		}

		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("read report: %w", err)
		}

		report, err := junitxml.Parse(reader)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("parse report: %w", err)
		}

		tr, ts, tc := junitxml.ConvertTestReport(report, jobRef)

		testReports = append(testReports, tr)
		testSuites = append(testSuites, ts...)
		testCases = append(testCases, tc...)
	}

	return testReports, testSuites, testCases, nil
}
