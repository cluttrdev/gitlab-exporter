package models

import (
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

func ConvertTestReport(pipelineID int64, report *gitlab.PipelineTestReport) (*typespb.TestReport, []*typespb.TestSuite, []*typespb.TestCase) {
	testreport := &typespb.TestReport{
		Id:           testReportID(pipelineID),
		PipelineId:   pipelineID,
		TotalTime:    report.TotalTime,
		TotalCount:   int64(report.TotalCount),
		SuccessCount: int64(report.SuccessCount),
		FailedCount:  int64(report.FailedCount),
		SkippedCount: int64(report.SkippedCount),
		ErrorCount:   int64(report.ErrorCount),
	}

	testsuites := make([]*typespb.TestSuite, 0, len(report.TestSuites))
	testcases := []*typespb.TestCase{}
	for i, testsuite := range report.TestSuites {
		testsuiteID := testSuiteID(testreport.Id, i)
		testsuites = append(testsuites, &typespb.TestSuite{
			Id:           testsuiteID,
			TestreportId: testreport.Id,
			PipelineId:   pipelineID,
			Name:         testsuite.Name,
			TotalTime:    testsuite.TotalTime,
			TotalCount:   int64(testsuite.TotalCount),
			SuccessCount: int64(testsuite.SuccessCount),
			FailedCount:  int64(testsuite.FailedCount),
			SkippedCount: int64(testsuite.SkippedCount),
			ErrorCount:   int64(testsuite.ErrorCount),
		})

		cases := make([]*typespb.TestCase, 0, len(testsuite.TestCases))
		for j, testcase := range testsuite.TestCases {
			cases = append(cases, &typespb.TestCase{
				Id:             testCaseID(testsuiteID, j),
				TestsuiteId:    testsuiteID,
				TestreportId:   testreport.Id,
				PipelineId:     pipelineID,
				Status:         testcase.Status,
				Name:           testcase.Name,
				Classname:      testcase.Classname,
				File:           testcase.File,
				ExecutionTime:  testcase.ExecutionTime,
				SystemOutput:   fmt.Sprint(testcase.SystemOutput),
				StackTrace:     testcase.StackTrace,
				AttachmentUrl:  testcase.AttachmentURL,
				RecentFailures: convertTestCaseRecentFailures(testcase.RecentFailures),
			})
		}
		testcases = append(testcases, cases...)
	}

	return testreport, testsuites, testcases
}

func convertTestCaseRecentFailures(f *gitlab.RecentFailures) *typespb.TestCase_RecentFailures {
	var r typespb.TestCase_RecentFailures
	if f != nil {
		r.Count = int64(f.Count)
		r.BaseBranch = f.BaseBranch
	}
	return &r
}

func testReportID(pipelineID int64) string {
	return fmt.Sprint(pipelineID)
}

func testSuiteID(reportID string, suiteIndex int) string {
	return fmt.Sprintf("%s-%d", reportID, suiteIndex+1)
}

func testCaseID(suiteID string, caseIndex int) string {
	return fmt.Sprintf("%s-%d", suiteID, caseIndex+1)
}
