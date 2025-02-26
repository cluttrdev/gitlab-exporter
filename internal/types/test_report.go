package types

import (
	"time"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type TestReportReference struct {
	Id       string
	Pipeline PipelineReference
}

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

type TestSuiteReference struct {
	Id         string
	TestReport TestReportReference
}

type TestSuite struct {
	Id         string
	TestReport TestReportReference

	Name         string
	TotalTime    time.Duration
	TotalCount   int64
	ErrorCount   int64
	FailedCount  int64
	SkippedCount int64
	SuccessCount int64

	Properties []TestProperty
}

type TestCase struct {
	Id        string
	TestSuite TestSuiteReference

	Name          string
	Classname     string
	Status        string
	ExecutionTime time.Duration
	File          string
	StackTrace    string
	SystemOutput  string
	AttachmentUrl string

	Properties []TestProperty
}

type TestCaseStatus string

const (
	TestCaseStatusFailed  = "failed"
	TestCaseStatusError   = "error"
	TestCaseStatusSkipped = "skipped"
	TestCaseStatusSuccess = "success"
)

type TestProperty struct {
	Name  string
	Value string
}

func ConvertTestReportReference(testReport TestReportReference) *typespb.TestReportReference {
	return &typespb.TestReportReference{
		Id:       testReport.Id,
		Pipeline: ConvertPipelineReference(testReport.Pipeline),
	}
}

func ConvertTestReport(testReport TestReport) *typespb.TestReport {
	return &typespb.TestReport{
		Id:       testReport.Id,
		Pipeline: ConvertPipelineReference(testReport.Pipeline),

		TotalTime:    testReport.TotalTime.Seconds(),
		TotalCount:   testReport.TotalCount,
		ErrorCount:   testReport.ErrorCount,
		FailedCount:  testReport.FailedCount,
		SkippedCount: testReport.SkippedCount,
		SuccessCount: testReport.SuccessCount,
	}
}

func ConvertTestSuiteReference(testSuite TestSuiteReference) *typespb.TestSuiteReference {
	return &typespb.TestSuiteReference{
		Id:         testSuite.Id,
		TestReport: ConvertTestReportReference(testSuite.TestReport),
	}
}

func ConvertTestSuite(testSuite TestSuite) *typespb.TestSuite {
	ts := &typespb.TestSuite{
		Id:         testSuite.Id,
		TestReport: ConvertTestReportReference(testSuite.TestReport),

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

func ConvertTestCase(testCase TestCase) *typespb.TestCase {
	tc := &typespb.TestCase{
		Id:        testCase.Id,
		TestSuite: ConvertTestSuiteReference(testCase.TestSuite),

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
