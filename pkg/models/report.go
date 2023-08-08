package models

import (
	"fmt"

	gogitlab "github.com/xanzy/go-gitlab"
)

type PipelineTestReport struct {
	Pipeline struct {
		ID int64
	}
	TotalTime    float64               `json:"total_time"`
	TotalCount   int                   `json:"total_count"`
	SuccessCount int                   `json:"success_count"`
	FailedCount  int                   `json:"failed_count"`
	SkippedCount int                   `json:"skipped_count"`
	ErrorCount   int                   `json:"error_count"`
	TestSuites   []*PipelineTestSuites `json:"test_suites"`
}

type PipelineTestSuites struct {
	Name         string               `json:"name"`
	TotalTime    float64              `json:"total_time"`
	TotalCount   int                  `json:"total_count"`
	SuccessCount int                  `json:"success_count"`
	FailedCount  int                  `json:"failed_count"`
	SkippedCount int                  `json:"skipped_count"`
	ErrorCount   int                  `json:"error_count"`
	TestCases    []*PipelineTestCases `json:"test_cases"`
}

type PipelineTestCases struct {
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
	Count      int    `json:"count"`
	BaseBranch string `json:"base_branch"`
}

func NewPipelineTestReport(pipelineID int64, tr *gogitlab.PipelineTestReport) *PipelineTestReport {
	testSuites := []*PipelineTestSuites{}
	for _, ts := range tr.TestSuites {
		testSuites = append(testSuites, NewPipelineTestSuites(ts))
	}
	return &PipelineTestReport{
		Pipeline:     struct{ ID int64 }{pipelineID},
		TotalTime:    tr.TotalTime,
		TotalCount:   tr.TotalCount,
		SuccessCount: tr.SuccessCount,
		FailedCount:  tr.FailedCount,
		SkippedCount: tr.SkippedCount,
		ErrorCount:   tr.ErrorCount,
		TestSuites:   testSuites,
	}
}

func NewPipelineTestSuites(ts *gogitlab.PipelineTestSuites) *PipelineTestSuites {
	testCases := []*PipelineTestCases{}
	for _, tc := range ts.TestCases {
		testCases = append(testCases, NewPipelineTestCases(tc))
	}
	return &PipelineTestSuites{
		Name:         ts.Name,
		TotalTime:    ts.TotalTime,
		TotalCount:   ts.TotalCount,
		SuccessCount: ts.SuccessCount,
		FailedCount:  ts.FailedCount,
		SkippedCount: ts.SkippedCount,
		ErrorCount:   ts.ErrorCount,
		TestCases:    testCases,
	}
}

func NewPipelineTestCases(tc *gogitlab.PipelineTestCases) *PipelineTestCases {
	return &PipelineTestCases{
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

func NewRecentFailures(rf *gogitlab.RecentFailures) *RecentFailures {
	return &RecentFailures{
		Count:      rf.Count,
		BaseBranch: rf.BaseBranch,
	}
}
