package gitlab

import (
	"context"
	"fmt"
	"net/http"

	_gitlab "github.com/xanzy/go-gitlab"

	pb "github.com/cluttrdev/gitlab-exporter/grpc/exporterpb"
	"github.com/cluttrdev/gitlab-exporter/internal/models"
)

type PipelineTestReportData struct {
	TestReports []*pb.TestReport
	TestSuites  []*pb.TestSuite
	TestCases   []*pb.TestCase
}

func (c *Client) GetPipelineTestReport(ctx context.Context, projectID int64, pipelineID int64) (*PipelineTestReportData, error) {
	c.RLock()
	defer c.RUnlock()
	report, _, err := c.client.Pipelines.GetPipelineTestReport(int(projectID), int(pipelineID), _gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error getting pipeline test report: %w", err)
	}

	testreport, testsuites, testcases := models.ConvertTestReport(pipelineID, report)

	return &PipelineTestReportData{
		TestReports: []*pb.TestReport{testreport},
		TestSuites:  testsuites,
		TestCases:   testcases,
	}, nil
}

func (c *Client) GetPipelineHierarchyTestReports(ctx context.Context, ph *models.PipelineHierarchy) (*PipelineTestReportData, error) {
	var results PipelineTestReportData

	result, err := c.GetPipelineTestReport(ctx, ph.Pipeline.ProjectId, ph.Pipeline.Id)
	if err != nil {
		return nil, err
	}

	results.TestReports = append(result.TestReports, result.TestReports...)
	results.TestSuites = append(result.TestSuites, result.TestSuites...)
	results.TestCases = append(result.TestCases, result.TestCases...)

	for _, dph := range ph.DownstreamPipelines {
		rs, err := c.GetPipelineHierarchyTestReports(ctx, dph)
		if err != nil {
			return nil, err
		}
		results.TestReports = append(results.TestReports, rs.TestReports...)
		results.TestSuites = append(results.TestSuites, rs.TestSuites...)
		results.TestCases = append(results.TestCases, rs.TestCases...)
	}

	return &results, nil
}

type PipelineTestReportSummary struct {
	Total      *PipelineTestReportSummaryTotal       `json:"total"`
	TestSuites []*PipelineTestReportSummaryTestSuite `json:"test_suites"`
}

type PipelineTestReportSummaryTotal struct {
	Time       float64 `json:"time"`
	Count      int     `json:"count"`
	Success    int     `json:"success"`
	Failed     int     `json:"failed"`
	Skipped    int     `json:"skipped"`
	Error      int     `json:"error"`
	SuiteError string  `json:"suite_error"`
}

type PipelineTestReportSummaryTestSuite struct {
	Name         string  `json:"name"`
	TotalTime    float64 `json:"total_time"`
	TotalCount   int     `json:"total_count"`
	SuccessCount int     `json:"success_count"`
	FailedCount  int     `json:"failed_count"`
	SkippedCount int     `json:"skipped_count"`
	ErrorCount   int     `json:"error_count"`
	BuildIDs     []int   `json:"build_ids"`
	SuiteError   string  `json:"suite_error"`
}

func (c *Client) GetPipelineTestReportSummary(ctx context.Context, projectID int64, pipelineID int64) (*PipelineTestReportSummary, error) {
	u := fmt.Sprintf("projects/%d/pipelines/%d/test_report_summary", int(projectID), int(pipelineID))

	options := []_gitlab.RequestOptionFunc{
		_gitlab.WithContext(ctx),
	}

	req, err := c.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, err
	}

	p := new(PipelineTestReportSummary)
	_, err = c.client.Do(req, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
