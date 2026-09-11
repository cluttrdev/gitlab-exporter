package tasks

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"go.cluttr.dev/junitxml"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/cobertura"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/version"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/types"
)

// projectPipelineUnit is one unit of fan-out work: a single pipeline of a single
// project, plus that project's configured artifact paths.
// Flattening the two input maps into a slice lets fanOut iterate a flat []T
// instead of a nested project→pipelines loop.
type projectPipelineUnit struct {
	projectPath   string
	pipelineIid   string
	artifactPaths []string
}

// flattenProjectPipelines expands the projectPath→pipelineIids map into one entry
// per (project, pipeline), attaching each project's artifact paths.
func flattenProjectPipelines(projectPipelines, projectArtifactPaths map[string][]string) []projectPipelineUnit {
	var out []projectPipelineUnit
	for projectPath, pipelineIids := range projectPipelines {
		artifactPaths := projectArtifactPaths[projectPath]
		for _, pipelineIid := range pipelineIids {
			out = append(out, projectPipelineUnit{
				projectPath:   projectPath,
				pipelineIid:   pipelineIid,
				artifactPaths: artifactPaths,
			})
		}
	}
	return out
}

// parseGzipped decompresses r and runs parse over the decompressed stream.
func parseGzipped[T any](r io.Reader, parse func(io.Reader) (T, error)) (T, error) {
	var zero T

	gz, err := gzip.NewReader(r)
	if err != nil {
		return zero, fmt.Errorf("read report: %w", err)
	}

	v, err := parse(gz)
	if err != nil {
		return zero, fmt.Errorf("parse report: %w", err)
	}
	return v, nil
}

// ############################################################################
// # Test Reports (JUnit)
// ############################################################################

func FetchProjectsPipelinesJunitReports(ctx context.Context, glab *gitlab.Client, projectPipelines map[string][]string, projectArtifactPaths map[string][]string) ([]types.TestReport, []types.TestSuite, []types.TestCase, error) {
	type junitBundle struct {
		reports []types.TestReport
		suites  []types.TestSuite
		cases   []types.TestCase
	}

	units := flattenProjectPipelines(projectPipelines, projectArtifactPaths)

	bundles, err := fanOut(ctx, glab, units,
		func(ctx context.Context, glab *gitlab.Client, u projectPipelineUnit) (junitBundle, error) {
			tr, ts, tc, err := FetchProjectPipelineJunitReports(ctx, glab, u.projectPath, u.pipelineIid, u.artifactPaths)
			return junitBundle{reports: tr, suites: ts, cases: tc}, err
		},
		fanOutOptions{},
	)

	var (
		testReports []types.TestReport
		testSuites  []types.TestSuite
		testCases   []types.TestCase
	)
	for _, b := range bundles {
		testReports = append(testReports, b.reports...)
		testSuites = append(testSuites, b.suites...)
		testCases = append(testCases, b.cases...)
	}
	return testReports, testSuites, testCases, err
}

func FetchProjectPipelineJunitReports(ctx context.Context, glab *gitlab.Client, projectPath string, pipelineIid string, artifactPaths []string) ([]types.TestReport, []types.TestSuite, []types.TestCase, error) {
	var (
		testReports []types.TestReport
		testSuites  []types.TestSuite
		testCases   []types.TestCase
	)

	artifacts, err := glab.GraphQL.GetProjectPipelineJobsArtifacts(ctx, projectPath, pipelineIid)
	if errors.Is(err, context.Canceled) {
		return nil, nil, nil, err
	} else if err != nil {
		slog.Error("error getting project pipeline job artifacts", "error", err)
	}

	for _, artifact := range artifacts {
		if artifact.FileType == nil || *artifact.FileType != graphql.JobArtifactFileTypeJunit {
			continue
		}

		jobRef, err := graphql.ConvertJobReference(artifact.Job, artifact.Pipeline, artifact.Project)
		if err != nil {
			slog.Error("error converting job reference",
				slog.String("jobId", *artifact.Job.Id),
				slog.String("error", err.Error()),
			)
			continue
		}
		jobFinishedAt := artifact.JobFinishedAt

		var reports []junitxml.TestReport
		switch {
		case len(artifactPaths) > 0:
			// If explicit artifact archive paths are configured, use them.
			reports, err = fetchProjectJobJunitReportAPIFiles(ctx, glab, projectPath, jobRef.Id, artifactPaths)
		case version.Compare(glab.Version, version.GitLabVersion{Major: 19, Minor: 4}) >= 0:
			// If the gitlab server version is at least 19.4 (even pre-release on gitlab.com)
			// we can make use of the job artifacts API `file_type` attribute.
			reports, err = fetchProjectJobJunitReportAPI(ctx, glab, projectPath, jobRef.Id)
		case artifact.DownloadPath != nil:
			// Lastly, try downloading the report via session-authed http route.
			reports, err = fetchProjectJobJunitReportHTTP(ctx, glab, *artifact.DownloadPath)
		default:
			continue
		}
		if errors.Is(err, context.Canceled) {
			return nil, nil, nil, err
		} else if err != nil {
			var downloadPath string
			if artifact.DownloadPath != nil {
				downloadPath = *artifact.DownloadPath
			}
			slog.Error("error fetching junit report",
				slog.String("projectPath", projectPath),
				slog.Int64("jobId", jobRef.Id),
				slog.String("downloadPath", downloadPath),
				slog.String("error", err.Error()),
			)
			continue
		}

		for _, report := range reports {
			tr, ts, tc := types.ConvertTestReport(report, jobRef)
			// set testcase's report_created_at to job.finished_at
			for i := range len(tc) {
				tc[i].ReportCreatedAt = jobFinishedAt
			}

			testReports = append(testReports, tr)
			testSuites = append(testSuites, ts...)
			testCases = append(testCases, tc...)
		}
	}

	return testReports, testSuites, testCases, nil
}

func fetchProjectJobJunitReportAPI(ctx context.Context, glab *gitlab.Client, projectPath string, jobId int64) ([]junitxml.TestReport, error) {
	in, err := glab.Rest.GetProjectJobArtifacts(ctx, projectPath, int(jobId), "junit")
	if err != nil {
		return nil, fmt.Errorf("get job artifacts: junit: %w", err)
	}
	return parseGzipped(in, junitxml.ParseMany) // junit.xml.gz
}

func fetchProjectJobJunitReportAPIFiles(ctx context.Context, glab *gitlab.Client, projectPath string, jobId int64, artifactPaths []string) ([]junitxml.TestReport, error) {
	var reports []junitxml.TestReport

	for _, path := range artifactPaths {
		reader, err := glab.Rest.DownloadProjectJobArtifactsFile(ctx, projectPath, jobId, path)
		if errors.Is(err, context.Canceled) {
			return nil, err
		} else if errors.Is(err, gitlab.ErrNotFound) {
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
	return parseGzipped(resp.Body, junitxml.ParseMany) // junit.xml.gz
}

// ############################################################################
// # Coverage Reports (Cobertura)
// ############################################################################

func FetchProjectsPipelinesCoberturaReports(ctx context.Context, glab *gitlab.Client, projectPipelines map[string][]string, projectArtifactPaths map[string][]string) ([]types.CoverageReport, []types.CoveragePackage, []types.CoverageClass, []types.CoverageMethod, error) {
	type coberturaBundle struct {
		reports  []types.CoverageReport
		packages []types.CoveragePackage
		classes  []types.CoverageClass
		methods  []types.CoverageMethod
	}

	units := flattenProjectPipelines(projectPipelines, projectArtifactPaths)

	// No options: success-only, as the original. FetchProjectPipelineCoberturaReports
	// only returns errors with nil data, so nothing partial is lost on error.
	bundles, err := fanOut(ctx, glab, units,
		func(ctx context.Context, glab *gitlab.Client, u projectPipelineUnit) (coberturaBundle, error) {
			cr, cp, cc, cm, err := FetchProjectPipelineCoberturaReports(ctx, glab, u.projectPath, u.pipelineIid, u.artifactPaths)
			return coberturaBundle{reports: cr, packages: cp, classes: cc, methods: cm}, err
		},
		fanOutOptions{},
	)

	var (
		covReports  []types.CoverageReport
		covPackages []types.CoveragePackage
		covClasses  []types.CoverageClass
		covMethods  []types.CoverageMethod
	)
	for _, b := range bundles {
		covReports = append(covReports, b.reports...)
		covPackages = append(covPackages, b.packages...)
		covClasses = append(covClasses, b.classes...)
		covMethods = append(covMethods, b.methods...)
	}
	return covReports, covPackages, covClasses, covMethods, err
}

func FetchProjectPipelineCoberturaReports(ctx context.Context, glab *gitlab.Client, projectPath string, pipelineIid string, artifactPaths []string) ([]types.CoverageReport, []types.CoveragePackage, []types.CoverageClass, []types.CoverageMethod, error) {
	var (
		covReports  []types.CoverageReport
		covPackages []types.CoveragePackage
		covClasses  []types.CoverageClass
		covMethods  []types.CoverageMethod
	)

	appendReport := func(index int, report cobertura.CoverageReport, jobRef types.JobReference) {
		cr, cp, cc, cm := cobertura.ConvertCoverageReport(index, report, jobRef)
		covReports = append(covReports, cr)
		covPackages = append(covPackages, cp...)
		covClasses = append(covClasses, cc...)
		covMethods = append(covMethods, cm...)
	}

	artifacts, err := glab.GraphQL.GetProjectPipelineJobsArtifacts(ctx, projectPath, pipelineIid)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("get project pipeline job artifacts: %w", err)
	}

artifactsLoop:
	for _, artifact := range artifacts {
		if artifact.FileType == nil || *artifact.FileType != graphql.JobArtifactFileTypeCobertura {
			continue
		}

		jobRef, err := graphql.ConvertJobReference(artifact.Job, artifact.Pipeline, artifact.Project)
		if err != nil {
			slog.Error("error converting job reference",
				slog.String("jobId", *artifact.Job.Id),
				slog.String("error", err.Error()),
			)
			continue
		}

		var report cobertura.CoverageReport
		switch {
		case len(artifactPaths) > 0:
			reportCounter := 0
			for _, path := range artifactPaths {
				report, err = fetchProjectJobCoberturaReportAPIFiles(ctx, glab, projectPath, jobRef.Id, path)
				if errors.Is(err, context.Canceled) {
					return nil, nil, nil, nil, err
				} else if err != nil {
					if !errors.Is(err, gitlab.ErrNotFound) { // 404 is fine
						slog.Error("error fetching cobertura report",
							slog.String("downloadPath", path),
							slog.String("error", err.Error()),
						)
					}
					continue artifactsLoop
				}
				appendReport(reportCounter, report, jobRef)
				reportCounter++
			}
		case version.Compare(glab.Version, version.GitLabVersion{Major: 19, Minor: 4}) >= 0:
			report, err = fetchProjectJobCoberturaReportAPI(ctx, glab, projectPath, jobRef.Id)
			if errors.Is(err, context.Canceled) {
				return nil, nil, nil, nil, err
			} else if err != nil {
				slog.Error("error fetching cobertura report",
					slog.String("projectPath", projectPath),
					slog.Int64("jobId", jobRef.Id),
					slog.String("error", err.Error()),
				)
				continue artifactsLoop
			}
			appendReport(0, report, jobRef)
		case artifact.DownloadPath != nil:
			report, err = fetchProjectJobCoberturaReportHTTP(ctx, glab, *artifact.DownloadPath)
			if errors.Is(err, context.Canceled) {
				return nil, nil, nil, nil, err
			} else if err != nil {
				slog.Error("error fetching cobertura report",
					slog.String("downloadPath", *artifact.DownloadPath),
					slog.String("error", err.Error()),
				)
				continue artifactsLoop
			}
			appendReport(0, report, jobRef)
		default:
			continue
		}
	}

	return covReports, covPackages, covClasses, covMethods, nil
}

func fetchProjectJobCoberturaReportAPI(ctx context.Context, glab *gitlab.Client, projectPath string, jobId int64) (cobertura.CoverageReport, error) {
	in, err := glab.Rest.GetProjectJobArtifacts(ctx, projectPath, int(jobId), "cobertura")
	if err != nil {
		return cobertura.CoverageReport{}, fmt.Errorf("get job artifacts: cobertura: %w", err)
	}
	return parseGzipped(in, cobertura.Parse) // cobertura-coverage.xml.gz
}

func fetchProjectJobCoberturaReportAPIFiles(ctx context.Context, glab *gitlab.Client, projectPath string, jobId int64, artifactPath string) (cobertura.CoverageReport, error) {
	reader, err := glab.Rest.DownloadProjectJobArtifactsFile(ctx, projectPath, jobId, artifactPath)
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
	return parseGzipped(resp.Body, cobertura.Parse) // cobertura-coverage.xml.gz
}
