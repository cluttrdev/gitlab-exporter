package gitlab

import (
	"context"
	"fmt"

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
	report, _, err := c.client.Pipelines.GetPipelineTestReport(int(projectID), int(pipelineID), _gitlab.WithContext(ctx))
	c.RUnlock()
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
