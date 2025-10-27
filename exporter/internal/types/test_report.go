package types

import (
	"fmt"
	"strings"
	"time"

	"go.cluttr.dev/junitxml"
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

func ConvertTestReport(xmlReport junitxml.TestReport, job JobReference) (TestReport, []TestSuite, []TestCase) {
	testReportId := fmt.Sprintf("%d-%d", job.Pipeline.Id, job.Id)
	testReport := TestReport{
		Id:  testReportId,
		Job: job,

		// TotalTime:    time.Duration(xmlReport.Time * float64(time.Second)),
		// TotalCount:   xmlReport.Tests,
		// FailedCount:  xmlReport.Failures,
		// ErrorCount:   xmlReport.Errors,
		// SkippedCount: xmlReport.Skipped,
		// SuccessCount: xmlReport.Tests - (xmlReport.Failures + xmlReport.Errors + xmlReport.Skipped),
	}

	testReportRef := TestReportReference{
		Id:  testReport.Id,
		Job: testReport.Job,
	}
	testSuites, testCases := convertTestSuites(xmlReport.TestSuites, testReportRef)

	// accumulate test suite stats
	for _, ts := range testSuites {
		testReport.TotalCount += ts.TotalCount
		testReport.TotalTime += ts.TotalTime
		testReport.FailedCount += ts.FailedCount
		testReport.ErrorCount += ts.ErrorCount
		testReport.SkippedCount += ts.SkippedCount
		testReport.SuccessCount += ts.SuccessCount
	}

	return testReport, testSuites, testCases
}

func convertTestSuites(xmlSuites []junitxml.TestSuite, testReportRef TestReportReference) ([]TestSuite, []TestCase) {
	var (
		testSuites []TestSuite
		testCases  []TestCase
	)

	for i, ts := range xmlSuites {
		testSuiteId := fmt.Sprintf("%s-%d", testReportRef.Id, i+1)
		testSuite := TestSuite{
			Id:         testSuiteId,
			TestReport: testReportRef,

			Name: ts.Name,
			// TotalTime:    time.Duration(ts.Time * float64(time.Second)),
			// TotalCount:   ts.Tests,
			// FailedCount:  ts.Failures,
			// ErrorCount:   ts.Errors,
			// SkippedCount: ts.Skipped,
			// SuccessCount: ts.Tests - (ts.Failures + ts.Errors + ts.Skipped),

			Properties: convertTestProperties(ts.Properties),
		}

		testSuiteRef := TestSuiteReference{
			Id:         testSuite.Id,
			TestReport: testSuite.TestReport,
		}
		testSuiteCases := convertTestCases(ts.TestCases, testSuiteRef)

		// accumulate test case stats
		for _, tc := range testSuiteCases {
			testSuite.TotalCount++
			testSuite.TotalTime += tc.ExecutionTime
			switch tc.Status {
			case TestCaseStatusFailed:
				testSuite.FailedCount++
			case TestCaseStatusError:
				testSuite.ErrorCount++
			case TestCaseStatusSkipped:
				testSuite.SkippedCount++
			case TestCaseStatusSuccess:
				testSuite.SuccessCount++
			}
		}

		testSuites = append(testSuites, testSuite)
		testCases = append(testCases, testSuiteCases...)
	}

	return testSuites, testCases
}

func convertTestCases(xmlCases []junitxml.TestCase, testSuiteRef TestSuiteReference) []TestCase {
	testCases := make([]TestCase, 0, len(xmlCases))

	for i, tc := range xmlCases {
		// see: https://gitlab.com/gitlab-org/gitlab/-/blob/master/lib/gitlab/ci/parsers/test/junit.rb#create_test_case

		var (
			status string
			output strings.Builder
		)

		switch {
		case tc.Failure != nil:
			status = TestCaseStatusFailed
			output.WriteString(formatTestOutput(tc.Failure.Message, tc.Failure.Type, tc.Failure.Text))
		case tc.Error != nil:
			status = TestCaseStatusError
			output.WriteString(formatTestOutput(tc.Error.Message, tc.Error.Type, tc.Error.Text))
		case tc.Skipped != nil:
			status = TestCaseStatusSkipped
			output.WriteString(tc.Skipped.Message)
		default:
			status = TestCaseStatusSuccess
		}

		if tc.SystemOut != nil {
			if output.Len() > 0 {
				output.WriteString("\n\n")
			}
			output.WriteString("System Out:\n\n")
			output.WriteString(tc.SystemOut.Text)
		}
		if tc.SystemErr != nil {
			if output.Len() > 0 {
				output.WriteString("\n\n")
			}
			output.WriteString("System Err:\n\n")
			output.WriteString(tc.SystemErr.Text)
		}

		attachements := junitxml.ParseTextAttachments(output.String())
		properties := convertTestProperties(append(tc.Properties, junitxml.ParseTextProperties(output.String())...))

		testCaseId := fmt.Sprintf("%s-%d", testSuiteRef.Id, i+1)
		testCase := TestCase{
			Id:        testCaseId,
			TestSuite: testSuiteRef,

			Name:          tc.Name,
			Classname:     tc.Classname,
			ExecutionTime: time.Duration(tc.Time * float64(time.Second)),
			File:          tc.File,
			StackTrace:    "",

			Status:        status,
			SystemOutput:  output.String(),
			AttachmentUrl: strings.Join(attachements, "\n"),

			Properties: properties,
		}

		testCases = append(testCases, testCase)
	}

	return testCases
}

func convertTestProperties(xmlProps []junitxml.Property) []TestProperty {
	if len(xmlProps) == 0 {
		return nil
	}

	properties := make([]TestProperty, 0, len(xmlProps))
	for _, p := range xmlProps {
		value := p.Value
		if value == "" {
			value = p.Text
		}
		properties = append(properties, TestProperty{Name: p.Name, Value: value})
	}

	return properties
}

func formatTestOutput(message string, typ string, text string) string {
	var output string
	if typ != "" {
		output += typ
	}
	if message != "" {
		if len(output) > 0 {
			output += ": "
		}
		output += message
	}
	if text != "" {
		if len(output) > 0 {
			output += "\n\n"
		}
		output += text
	}
	return output
}
