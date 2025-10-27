package messages

import (
	"go.cluttr.dev/gitlab-exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewTestReportReference(testReport types.TestReportReference) *typespb.TestReportReference {
	return &typespb.TestReportReference{
		Id:  testReport.Id,
		Job: NewJobReference(testReport.Job),
	}
}

func NewTestReport(testReport types.TestReport) *typespb.TestReport {
	return &typespb.TestReport{
		Id:  testReport.Id,
		Job: NewJobReference(testReport.Job),

		TotalTime:    testReport.TotalTime.Seconds(),
		TotalCount:   testReport.TotalCount,
		ErrorCount:   testReport.ErrorCount,
		FailedCount:  testReport.FailedCount,
		SkippedCount: testReport.SkippedCount,
		SuccessCount: testReport.SuccessCount,
	}
}

func NewTestSuiteReference(testSuite types.TestSuiteReference) *typespb.TestSuiteReference {
	return &typespb.TestSuiteReference{
		Id:         testSuite.Id,
		TestReport: NewTestReportReference(testSuite.TestReport),
	}
}

func NewTestSuite(testSuite types.TestSuite) *typespb.TestSuite {
	ts := &typespb.TestSuite{
		Id:         testSuite.Id,
		TestReport: NewTestReportReference(testSuite.TestReport),

		Name:         testSuite.Name,
		TotalTime:    testSuite.TotalTime.Seconds(),
		TotalCount:   testSuite.TotalCount,
		ErrorCount:   testSuite.ErrorCount,
		FailedCount:  testSuite.FailedCount,
		SkippedCount: testSuite.SkippedCount,
		SuccessCount: testSuite.SuccessCount,
	}

	for _, p := range testSuite.Properties {
		ts.Properties = append(ts.Properties, &typespb.TestProperty{
			Name:  p.Name,
			Value: p.Value,
		})
	}

	return ts
}

func NewTestCase(testCase types.TestCase) *typespb.TestCase {
	tc := &typespb.TestCase{
		Id:        testCase.Id,
		TestSuite: NewTestSuiteReference(testCase.TestSuite),

		Status:        testCase.Status,
		Name:          testCase.Name,
		Classname:     testCase.Classname,
		ExecutionTime: testCase.ExecutionTime.Seconds(),
		File:          testCase.File,
		StackTrace:    testCase.StackTrace,
		SystemOutput:  testCase.SystemOutput,
		AttachmentUrl: testCase.AttachmentUrl,
	}

	for _, p := range testCase.Properties {
		tc.Properties = append(tc.Properties, &typespb.TestProperty{
			Name:  p.Name,
			Value: p.Value,
		})
	}

	return tc
}
