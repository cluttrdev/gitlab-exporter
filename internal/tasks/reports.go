package tasks

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"go.cluttr.dev/junitxml"

	"go.cluttr.dev/gitlab-exporter/internal/cobertura"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/internal/types"
)

// ############################################################################
// # Test Reports (JUnit)
// ############################################################################

func FetchProjectsPipelinesJunitReports(ctx context.Context, glab *gitlab.Client, projectPipelines map[string][]string, projectArtifactPaths map[string][]string) ([]types.TestReport, []types.TestSuite, []types.TestCase, error) {
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
			artifactPaths := projectArtifactPaths[projectPath]
			for _, pipelineIid := range pipelineIids {
				if err := glab.Acquire(ctx, 1); err != nil {
					slog.Error("failed to acquire gitlab client", "error", err)
					continue
				}
				wg.Add(1)
				go func(projectPath string, pipelineIid string, artifactPaths []string) {
					defer glab.Release(1)
					defer wg.Done()

					tr, ts, tc, err := FetchProjectPipelineJunitReports(ctx, glab, projectPath, pipelineIid, artifactPaths)

					results <- result{
						testReports: tr,
						testSuites:  ts,
						testCases:   tc,
						err:         err,
					}
				}(projectPath, pipelineIid, artifactPaths)
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

func FetchProjectPipelineJunitReports(ctx context.Context, glab *gitlab.Client, projectPath string, pipelineIid string, artifactPaths []string) ([]types.TestReport, []types.TestSuite, []types.TestCase, error) {
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

		jobRef, err := graphql.ConvertJobReference(artifact.Job, artifact.Pipeline, artifact.Project)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("convert job reference: %w", err)
		}

		var reports []junitxml.TestReport
		if len(artifactPaths) > 0 {
			reports, err = fetchProjectJobJunitReportsAPI(ctx, glab, projectPath, jobRef.Id, artifactPaths)
		} else if artifact.DownloadPath != nil {
			reports, err = fetchProjectJobJunitReportHTTP(ctx, glab, *artifact.DownloadPath)
		} else {
			continue
		}
		if err != nil {
			return nil, nil, nil, err
		}

		for _, report := range reports {
			tr, ts, tc := types.ConvertTestReport(report, jobRef)

			testReports = append(testReports, tr)
			testSuites = append(testSuites, ts...)
			testCases = append(testCases, tc...)
		}
	}

	return testReports, testSuites, testCases, nil
}

func fetchProjectJobJunitReportsAPI(ctx context.Context, glab *gitlab.Client, projectPath string, jobId int64, artifactPaths []string) ([]junitxml.TestReport, error) {
	var reports []junitxml.TestReport

	for _, path := range artifactPaths {
		reader, err := glab.Rest.GetProjectJobArtifact(ctx, projectPath, jobId, path)
		if errors.Is(err, gitlab.ErrNotFound) {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("download file: %w", err)
		}

		rs, err := junitxml.ParseMany(reader)
		if err != nil {
			return nil, fmt.Errorf("parse file: %w", err)
		}

		reports = append(reports, rs...)
	}

	return reports, nil
}

func fetchProjectJobJunitReportHTTP(ctx context.Context, glab *gitlab.Client, downloadPath string) ([]junitxml.TestReport, error) {
	resp, err := glab.HTTP.GetPath(downloadPath)
	if err != nil {
		return nil, fmt.Errorf("download report: %w", err)
	}

	// junit.xml.gz
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read report: %w", err)
	}

	reports, err := junitxml.ParseMany(reader)
	if err != nil {
		return nil, fmt.Errorf("parse report: %w", err)
	}

	return reports, nil
}

// ############################################################################
// # Coverage Reports (Cobertura)
// ############################################################################

func FetchProjectsPipelinesCoberturaReports(ctx context.Context, glab *gitlab.Client, projectPipelines map[string][]string, projectArtifactPaths map[string][]string) ([]types.CoverageReport, []types.CoveragePackage, []types.CoverageClass, []types.CoverageMethod, error) {
	var (
		covReports  []types.CoverageReport
		covPackages []types.CoveragePackage
		covClasses  []types.CoverageClass
		covMethods  []types.CoverageMethod
	)

	type result struct {
		covReports  []types.CoverageReport
		covPackages []types.CoveragePackage
		covClasses  []types.CoverageClass
		covMethods  []types.CoverageMethod

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
			artifactPaths := projectArtifactPaths[projectPath]
			for _, pipelineIid := range pipelineIids {
				if err := glab.Acquire(ctx, 1); err != nil {
					slog.Error("failed to acquire gitlab client", "error", err)
					continue
				}
				wg.Add(1)
				go func(projectPath string, pipelineIid string, artifactPaths []string) {
					defer glab.Release(1)
					defer wg.Done()

					cr, cp, cc, cm, err := FetchProjectPipelineCoberturaReports(ctx, glab, projectPath, pipelineIid, artifactPaths)

					results <- result{
						covReports:  cr,
						covPackages: cp,
						covClasses:  cc,
						covMethods:  cm,
						err:         err,
					}
				}(projectPath, pipelineIid, artifactPaths)
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
				covReports = append(covReports, r.covReports...)
				covPackages = append(covPackages, r.covPackages...)
				covClasses = append(covClasses, r.covClasses...)
				covMethods = append(covMethods, r.covMethods...)
			}
		}
	}

	return covReports, covPackages, covClasses, covMethods, errs
}

func FetchProjectPipelineCoberturaReports(ctx context.Context, glab *gitlab.Client, projectPath string, pipelineIid string, artifactPaths []string) ([]types.CoverageReport, []types.CoveragePackage, []types.CoverageClass, []types.CoverageMethod, error) {
	var (
		covReports  []types.CoverageReport
		covPackages []types.CoveragePackage
		covClasses  []types.CoverageClass
		covMethods  []types.CoverageMethod
	)

	artifacts, err := glab.GraphQL.GetProjectPipelineJobsArtifacts(ctx, projectPath, pipelineIid)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get project pipeline job artifacts: %w", err)
	}

	for _, artifact := range artifacts {
		if artifact.FileType == nil || *artifact.FileType != graphql.JobArtifactFileTypeCobertura {
			continue
		}

		jobRef, err := graphql.ConvertJobReference(artifact.Job, artifact.Pipeline, artifact.Project)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("convert job reference: %w", err)
		}

		var report cobertura.CoverageReport
		if len(artifactPaths) > 0 {
			reportCounter := 0
			for _, path := range artifactPaths {
				report, err = fetchProjectJobCoberturaReportAPI(ctx, glab, projectPath, jobRef.Id, path)
				if errors.Is(err, gitlab.ErrNotFound) {
					continue
				} else if err != nil {
					return nil, nil, nil, nil, err
				}
				cr, cp, cc, cm := cobertura.ConvertCoverageReport(reportCounter, report, jobRef)
				reportCounter++

				covReports = append(covReports, cr)
				covPackages = append(covPackages, cp...)
				covClasses = append(covClasses, cc...)
				covMethods = append(covMethods, cm...)
			}
		} else if artifact.DownloadPath != nil {
			report, err = fetchProjectJobCoberturaReportHTTP(ctx, glab, *artifact.DownloadPath)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			cr, cp, cc, cm := cobertura.ConvertCoverageReport(0, report, jobRef)
			covReports = append(covReports, cr)
			covPackages = append(covPackages, cp...)
			covClasses = append(covClasses, cc...)
			covMethods = append(covMethods, cm...)
		} else {
			continue
		}
	}

	return covReports, covPackages, covClasses, covMethods, nil
}

func fetchProjectJobCoberturaReportAPI(ctx context.Context, glab *gitlab.Client, projectPath string, jobId int64, artifactPath string) (cobertura.CoverageReport, error) {
	reader, err := glab.Rest.GetProjectJobArtifact(ctx, projectPath, jobId, artifactPath)
	if err != nil {
		return cobertura.CoverageReport{}, fmt.Errorf("download file: %w", err)
	}

	report, err := cobertura.Parse(reader)
	if err != nil {
		return cobertura.CoverageReport{}, fmt.Errorf("parse file: %w", err)
	}

	return report, nil
}

func fetchProjectJobCoberturaReportHTTP(ctx context.Context, glab *gitlab.Client, downloadPath string) (cobertura.CoverageReport, error) {
	resp, err := glab.HTTP.GetPath(downloadPath)
	if err != nil {
		return cobertura.CoverageReport{}, fmt.Errorf("download report: %w", err)
	}

	// cobertura-coverage.xml.gz
	reader, err := gzip.NewReader(resp.Body)
	if err != nil {
		return cobertura.CoverageReport{}, fmt.Errorf("read report: %w", err)
	}

	report, err := cobertura.Parse(reader)
	if err != nil {
		return cobertura.CoverageReport{}, fmt.Errorf("parse report: %w", err)
	}

	return report, nil
}
