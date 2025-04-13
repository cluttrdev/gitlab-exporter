package types

import (
	"time"
)

type TestReportReference struct {
	Id  string
	Job JobReference
}

type TestReport struct {
	Id  string
	Job JobReference

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
