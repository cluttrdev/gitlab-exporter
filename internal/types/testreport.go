package types

import (
	"time"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type TestReport struct {
	Id       string
	Pipeline PipelineReference

	TotalTime    time.Duration
	TotalCount   int64
	ErrorCount   int64
	FailedCount  int64
	SkippedCount int64
	SuccessCount int64
}

type TestSuite struct {
	Id           string
	TestReportId string
	Pipeline     PipelineReference

	Name         string
	TotalTime    time.Duration
	TotalCount   int64
	ErrorCount   int64
	FailedCount  int64
	SkippedCount int64
	SuccessCount int64
}

type TestCase struct {
	Id           string
	TestSuiteId  string
	TestReportId string
	Pipeline     PipelineReference

	Name          string
	Classname     string
	Status        string
	ExecutionTime time.Duration
	File          string
	StackTrace    string
	SystemOutput  string
}

func ConvertTestReport(testreport TestReport) *typespb.TestReport {
	return &typespb.TestReport{
		Id:         testreport.Id,
		PipelineId: testreport.Pipeline.Id,

		TotalTime:    testreport.TotalTime.Seconds(),
		TotalCount:   testreport.TotalCount,
		ErrorCount:   testreport.ErrorCount,
		FailedCount:  testreport.FailedCount,
		SkippedCount: testreport.SkippedCount,
		SuccessCount: testreport.SuccessCount,
	}
}

func ConvertTestSuite(testsuite TestSuite) *typespb.TestSuite {
	return &typespb.TestSuite{
		Id:           testsuite.Id,
		TestreportId: testsuite.TestReportId,
		PipelineId:   testsuite.Pipeline.Id,

		Name:         testsuite.Name,
		TotalTime:    testsuite.TotalTime.Seconds(),
		TotalCount:   testsuite.TotalCount,
		ErrorCount:   testsuite.ErrorCount,
		FailedCount:  testsuite.FailedCount,
		SkippedCount: testsuite.SkippedCount,
		SuccessCount: testsuite.SuccessCount,
	}
}

func ConvertTestCase(testcase TestCase) *typespb.TestCase {
	return &typespb.TestCase{
		Id:           testcase.Id,
		TestsuiteId:  testcase.TestSuiteId,
		TestreportId: testcase.TestReportId,
		PipelineId:   testcase.Pipeline.Id,

		Name:          testcase.Name,
		Classname:     testcase.Classname,
		ExecutionTime: testcase.ExecutionTime.Seconds(),
		File:          testcase.File,
		StackTrace:    testcase.StackTrace,
		SystemOutput:  testcase.SystemOutput,

		RecentFailures: &typespb.TestCase_RecentFailures{},
	}
}
