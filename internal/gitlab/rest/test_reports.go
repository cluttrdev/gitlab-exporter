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
	testReportId := fmt.Sprint(pipeline.Id)
	testReport := types.TestReport{
		Id: testReportId,
		Pipeline: types.PipelineReference{
			Id:      pipeline.Id,
			Iid:     pipeline.Iid,
			Project: pipeline.Project,
		},

		TotalTime:    time.Duration(report.TotalTime * float64(time.Second)),
		TotalCount:   int64(report.TotalCount),
		ErrorCount:   int64(report.ErrorCount),
		FailedCount:  int64(report.FailedCount),
		SkippedCount: int64(report.SkippedCount),
		SuccessCount: int64(report.SuccessCount),
	}

	testSuites := make([]types.TestSuite, 0, len(report.TestSuites))
	testCases := []types.TestCase{}
	for _, ts := range report.TestSuites {
		index := slices.IndexFunc(summary.TestSuites, func(sts *PipelineTestReportSummaryTestSuite) bool {
			return ts.Name == sts.Name
		})
		if index < 0 {
			return types.TestReport{}, nil, nil, fmt.Errorf("cannot find test suite in summary: %s", ts.Name)
		}
		testSuiteSummary := summary.TestSuites[index]
		if len(testSuiteSummary.BuildIDs) == 0 {
			return types.TestReport{}, nil, nil, fmt.Errorf("test suite has no build id: %s", testSuiteSummary.Name)
		}

		testSuiteId := fmt.Sprint(testSuiteSummary.BuildIDs[0])
		testSuite := types.TestSuite{
			Id: testSuiteId,
			TestReport: types.TestReportReference{
				Id:       testReport.Id,
				Pipeline: testReport.Pipeline,
			},

			Name:         ts.Name,
			TotalTime:    time.Duration(ts.TotalTime * float64(time.Second)),
			TotalCount:   int64(ts.TotalCount),
			ErrorCount:   int64(ts.ErrorCount),
			FailedCount:  int64(ts.FailedCount),
			SkippedCount: int64(ts.SkippedCount),
			SuccessCount: int64(ts.SuccessCount),
		}

		testSuiteCases := make([]types.TestCase, 0, len(ts.TestCases))
		for j, tc := range ts.TestCases {
			testCaseId := fmt.Sprintf("%s-%d", testSuite.Id, j+1)
			testSuiteCases = append(testSuiteCases, types.TestCase{
				Id: testCaseId,
				TestSuite: types.TestSuiteReference{
					Id:         testSuite.Id,
					TestReport: testSuite.TestReport,
				},

				Status:        tc.Status,
				Name:          tc.Name,
				Classname:     tc.Classname,
				File:          tc.File,
				ExecutionTime: time.Duration(tc.ExecutionTime * float64(time.Second)),
				StackTrace:    tc.StackTrace,
				SystemOutput:  fmt.Sprint(tc.SystemOutput),
				AttachmentUrl: tc.AttachmentURL,
			})
		}

		testSuites = append(testSuites, testSuite)
		testCases = append(testCases, testSuiteCases...)
	}

	return testReport, testSuites, testCases, nil
}
