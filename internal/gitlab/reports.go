package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"slices"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/types"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type PipelineTestReportData struct {
	TestReport *typespb.TestReport
	TestSuites []*typespb.TestSuite
	TestCases  []*typespb.TestCase
}

func (c *Client) GetPipelineTestReport(ctx context.Context, projectID int64, pipelineID int64) (*PipelineTestReportData, error) {
	report, _, err := c.client.Pipelines.GetPipelineTestReport(int(projectID), int(pipelineID), gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("get pipeline test report: %w", err)
	}
	summary, err := c.GetPipelineTestReportSummary(ctx, projectID, pipelineID)
	if err != nil {
		return nil, fmt.Errorf("get pipeline test report summary: %w", err)
	}

	testreport, testsuites, testcases := types.ConvertTestReport(pipelineID, report)

	if err := overrideIDs(pipelineID, summary, testreport, testsuites, testcases); err != nil {
		return nil, fmt.Errorf("set test report ids: %w", err)
	}

	return &PipelineTestReportData{
		TestReport: testreport,
		TestSuites: testsuites,
		TestCases:  testcases,
	}, nil
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

func overrideIDs(pipelineID int64, summary *PipelineTestReportSummary, report *typespb.TestReport, suites []*typespb.TestSuite, cases []*typespb.TestCase) error {
	trID := fmt.Sprint(pipelineID)

	report.Id = fmt.Sprint(pipelineID)

	for _, rts := range suites {
		index := slices.IndexFunc(summary.TestSuites, func(sts *PipelineTestReportSummaryTestSuite) bool {
			return rts.Name == sts.Name
		})
		if index < 0 {
			return fmt.Errorf("cannot find test suite in summary: %s", rts.Name)
		}

		sts := summary.TestSuites[index]
		if len(sts.BuildIDs) == 0 {
			return fmt.Errorf("test suite has no build id: %s", sts.Name)
		}

		tsID := fmt.Sprint(sts.BuildIDs[0])

		tcNum := 0
		for _, tc := range cases {
			if tc.TestsuiteId == rts.Id {
				tcNum++
				tc.Id = fmt.Sprintf("%s-%d", tsID, tcNum)
				tc.TestsuiteId = tsID
				tc.TestreportId = trID
				tc.PipelineId = pipelineID
			}
		}

		rts.Id = tsID
		rts.TestreportId = trID
		rts.PipelineId = pipelineID
	}

	return nil
}
