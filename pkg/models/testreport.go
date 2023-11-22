package models

import (
	"fmt"
	"hash/adler32"

	_gitlab "github.com/xanzy/go-gitlab"
)

type PipelineTestReport struct {
	ID           int64
	PipelineID   int64
	TotalTime    float64              `json:"total_time"`
	TotalCount   int64                `json:"total_count"`
	SuccessCount int64                `json:"success_count"`
	FailedCount  int64                `json:"failed_count"`
	SkippedCount int64                `json:"skipped_count"`
	ErrorCount   int64                `json:"error_count"`
	TestSuites   []*PipelineTestSuite `json:"test_suites"`
}

type PipelineTestSuite struct {
	ID           int64
	TestReport   *TestReportReference
	Name         string              `json:"name"`
	TotalTime    float64             `json:"total_time"`
	TotalCount   int64               `json:"total_count"`
	SuccessCount int64               `json:"success_count"`
	FailedCount  int64               `json:"failed_count"`
	SkippedCount int64               `json:"skipped_count"`
	ErrorCount   int64               `json:"error_count"`
	TestCases    []*PipelineTestCase `json:"test_cases"`
}

type PipelineTestCase struct {
	ID             int64
	TestSuite      *TestSuiteReference
	TestReport     *TestReportReference
	Status         string          `json:"status"`
	Name           string          `json:"name"`
	Classname      string          `json:"classname"`
	File           string          `json:"file"`
	ExecutionTime  float64         `json:"execution_time"`
	SystemOutput   string          `json:"system_output"`
	StackTrace     string          `json:"stack_trace"`
	AttachmentURL  string          `json:"attachment_url"`
	RecentFailures *RecentFailures `json:"recent_failures"`
}

type RecentFailures struct {
	Count      int64  `json:"count"`
	BaseBranch string `json:"base_branch"`
}

type TestReportReference struct {
	ID         int64
	PipelineID int64
}

type TestSuiteReference struct {
	ID int64
}

func NewPipelineTestReport(pipelineID int64, tr *_gitlab.PipelineTestReport) *PipelineTestReport {
	report := &TestReportReference{
		ID:         hashStringID(fmt.Sprint(pipelineID)),
		PipelineID: pipelineID,
	}

	suites := []*PipelineTestSuite{}
	for _, ts := range tr.TestSuites {
		suites = append(suites, NewPipelineTestSuites(report, ts))
	}
	return &PipelineTestReport{
		ID:           report.ID,
		PipelineID:   pipelineID,
		TotalTime:    tr.TotalTime,
		TotalCount:   int64(tr.TotalCount),
		SuccessCount: int64(tr.SuccessCount),
		FailedCount:  int64(tr.FailedCount),
		SkippedCount: int64(tr.SkippedCount),
		ErrorCount:   int64(tr.ErrorCount),
		TestSuites:   suites,
	}
}

func NewPipelineTestSuites(report *TestReportReference, ts *_gitlab.PipelineTestSuites) *PipelineTestSuite {
	suite := &TestSuiteReference{
		ID: hashStringID(fmt.Sprint(report.ID) + ts.Name),
	}

	cases := []*PipelineTestCase{}
	for _, tc := range ts.TestCases {
		cases = append(cases, NewPipelineTestCases(report, suite, tc))
	}

	return &PipelineTestSuite{
		ID:           suite.ID,
		TestReport:   report,
		Name:         ts.Name,
		TotalTime:    ts.TotalTime,
		TotalCount:   int64(ts.TotalCount),
		SuccessCount: int64(ts.SuccessCount),
		FailedCount:  int64(ts.FailedCount),
		SkippedCount: int64(ts.SkippedCount),
		ErrorCount:   int64(ts.ErrorCount),
		TestCases:    cases,
	}
}

func NewPipelineTestCases(report *TestReportReference, suite *TestSuiteReference, tc *_gitlab.PipelineTestCases) *PipelineTestCase {
	return &PipelineTestCase{
		ID:             hashStringID(fmt.Sprint(suite.ID) + tc.Name),
		TestSuite:      suite,
		TestReport:     report,
		Status:         tc.Status,
		Name:           tc.Name,
		Classname:      tc.Classname,
		File:           tc.File,
		ExecutionTime:  tc.ExecutionTime,
		SystemOutput:   fmt.Sprint(tc.SystemOutput),
		StackTrace:     tc.StackTrace,
		AttachmentURL:  tc.AttachmentURL,
		RecentFailures: NewRecentFailures(tc.RecentFailures),
	}
}

func NewRecentFailures(rf *_gitlab.RecentFailures) *RecentFailures {
	if rf == nil {
		return &RecentFailures{}
	}
	return &RecentFailures{
		Count:      int64(rf.Count),
		BaseBranch: rf.BaseBranch,
	}
}

func hashStringID(s string) int64 {
	return int64(adler32.Checksum([]byte(s)))
}
