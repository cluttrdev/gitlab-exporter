package rest

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/types"
)

type PipelineTestReportData struct {
	TestReport types.TestReport
	TestSuites []types.TestSuite
	TestCases  []types.TestCase
}

func (c *Client) GetPipelineTestReport(ctx context.Context, projectID int64, pipelineID int64) (*gitlab.PipelineTestReport, *PipelineTestReportSummary, error) {
	report, _, err := c.client.Pipelines.GetPipelineTestReport(int(projectID), int(pipelineID), gitlab.WithContext(ctx))
	if err != nil {
		return nil, nil, fmt.Errorf("get pipeline test report: %w", err)
	}
	summary, err := c.GetPipelineTestReportSummary(ctx, projectID, pipelineID)
	if err != nil {
		return nil, nil, fmt.Errorf("get pipeline test report summary: %w", err)
	}

	return report, summary, nil
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

	options := []gitlab.RequestOptionFunc{
		gitlab.WithContext(ctx),
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

func ConvertTestReport(report *gitlab.PipelineTestReport, summary *PipelineTestReportSummary, pipeline types.Pipeline) (types.TestReport, []types.TestSuite, []types.TestCase, error) {
	pipelineRefs := types.PipelineReference{
		Id:      pipeline.Id,
		Iid:     pipeline.Iid,
		Project: pipeline.Project,
	}

	testReportId := fmt.Sprint(pipeline.Id)

	testReport := types.TestReport{
		Id:       testReportId,
		Pipeline: pipelineRefs,

		TotalTime:    time.Duration(report.TotalTime * float64(time.Second)),
		TotalCount:   int64(report.TotalCount),
		ErrorCount:   int64(report.ErrorCount),
		FailedCount:  int64(report.FailedCount),
		SkippedCount: int64(report.SkippedCount),
		SuccessCount: int64(report.SuccessCount),
	}

	testSuites := make([]types.TestSuite, 0, len(report.TestSuites))
	testCases := []types.TestCase{}
	for _, testSuite := range report.TestSuites {
		index := slices.IndexFunc(summary.TestSuites, func(sts *PipelineTestReportSummaryTestSuite) bool {
			return testSuite.Name == sts.Name
		})
		if index < 0 {
			return types.TestReport{}, nil, nil, fmt.Errorf("cannot find test suite in summary: %s", testSuite.Name)
		}
		testSuiteSummary := summary.TestSuites[index]
		if len(testSuiteSummary.BuildIDs) == 0 {
			return types.TestReport{}, nil, nil, fmt.Errorf("test suite has no build id: %s", testSuiteSummary.Name)
		}
		testSuiteId := fmt.Sprint(testSuiteSummary.BuildIDs[0])

		testSuites = append(testSuites, types.TestSuite{
			Id:           testSuiteId,
			TestReportId: testReport.Id,
			Pipeline:     pipelineRefs,

			Name:         testSuite.Name,
			TotalTime:    time.Duration(testSuite.TotalTime * float64(time.Second)),
			TotalCount:   int64(testSuite.TotalCount),
			ErrorCount:   int64(testSuite.ErrorCount),
			FailedCount:  int64(testSuite.FailedCount),
			SkippedCount: int64(testSuite.SkippedCount),
			SuccessCount: int64(testSuite.SuccessCount),
		})

		testSuiteCases := make([]types.TestCase, 0, len(testSuite.TestCases))
		for j, testcase := range testSuite.TestCases {
			testSuiteCases = append(testSuiteCases, types.TestCase{
				Id:            fmt.Sprintf("%s-%d", testSuiteId, j+1),
				TestSuiteId:   testSuiteId,
				TestReportId:  testReportId,
				Pipeline:      pipelineRefs,
				Status:        testcase.Status,
				Name:          testcase.Name,
				Classname:     testcase.Classname,
				File:          testcase.File,
				ExecutionTime: time.Duration(testcase.ExecutionTime * float64(time.Second)),
				StackTrace:    testcase.StackTrace,
				SystemOutput:  fmt.Sprint(testcase.SystemOutput),
			})
		}
		testCases = append(testCases, testSuiteCases...)
	}

	return testReport, testSuites, testCases, nil
}
