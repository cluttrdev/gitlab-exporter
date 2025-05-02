package junitxml

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/types"
)

type TestReport struct {
	XMLName    xml.Name    `xml:"testsuites"`
	Tests      int64       `xml:"tests,attr,omitempty"`
	Failures   int64       `xml:"failures,attr,omitempty"`
	Errors     int64       `xml:"errors,attr,omitempty"`
	Skipped    int64       `xml:"skipped,attr,omitempty"`
	Time       float64     `xml:"time,attr,omitempty"`
	Timestamp  string      `xml:"timestamp,attr,omitempty"`
	TestSuites []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	// XMLName    xml.Name   `xml:"testsuite"`
	Name       string     `xml:"name,attr,omitempty"`
	Tests      int64      `xml:"tests,attr,omitempty"`
	Failures   int64      `xml:"failures,attr,omitempty"`
	Errors     int64      `xml:"errors,attr,omitempty"`
	Skipped    int64      `xml:"skipped,attr,omitempty"`
	Time       float64    `xml:"time,attr,omitempty"`
	Timestamp  string     `xml:"timestamp,attr,omitempty"`
	File       string     `xml:"file,attr,omitempty"`
	Properties []Property `xml:"properties>property,omitempty"`
	SystemOut  *SystemOut `xml:"system-out"`
	SystemErr  *SystemErr `xml:"system-err"`
	TestCases  []TestCase `xml:"testcase"`
}

type TestCase struct {
	// XMLName   xml.Name  `xml:"testcase"`
	Name       string     `xml:"name,attr,omitempty"`
	Classname  string     `xml:"classname,attr,omitempty"`
	Tests      int64      `xml:"tests,attr,omitempty"`
	Time       float64    `xml:"time,attr,omitempty"`
	File       string     `xml:"file,attr,omitempty"`
	Line       int64      `xml:"line,attr,omitempty"`
	Failure    *Failure   `xml:"failure"`
	Error      *Error     `xml:"error"`
	Skipped    *Skipped   `xml:"skipped"`
	Properties []Property `xml:"properties>property,omitempty"`
	SystemOut  *SystemOut `xml:"system-out"`
	SystemErr  *SystemErr `xml:"system-err"`
}

type Failure struct {
	// XMLName xml.Name `xml:"failure"`
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
	Text    string `xml:",innerxml"`
}

func (f *Failure) Output() string {
	return formatTestOutput(f.Message, f.Type, f.Text)
}

type Error struct {
	// XMLName xml.Name `xml:"error"`
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
	Text    string `xml:",innerxml"`
}

func (e *Error) Output() string {
	return formatTestOutput(e.Message, e.Type, e.Text)
}

type Skipped struct {
	// XMLName xml.Name `xml:"skipped"`
	Message string `xml:"message,attr,omitempty"`
}

type Property struct {
	// XMLName xml.Name `xml:"property"`
	Name  string `xml:"name,attr,omitempty"`
	Value string `xml:"value,attr,omitempty"`
	Text  string `xml:",innerxml"`
}

type SystemOut struct {
	// XMLName xml.Name `xml:"system-out"`
	Text string `xml:",innerxml"`
}

type SystemErr struct {
	// XMLName xml.Name `xml:"system-err"`
	Text string `xml:",innerxml"`
}

func ConvertTestReport(xmlReport TestReport, job types.JobReference) (types.TestReport, []types.TestSuite, []types.TestCase) {
	testReportId := fmt.Sprintf("%d-%d", job.Pipeline.Id, job.Id)
	testReport := types.TestReport{
		Id:  testReportId,
		Job: job,

		// TotalTime:    time.Duration(xmlReport.Time * float64(time.Second)),
		// TotalCount:   xmlReport.Tests,
		// FailedCount:  xmlReport.Failures,
		// ErrorCount:   xmlReport.Errors,
		// SkippedCount: xmlReport.Skipped,
		// SuccessCount: xmlReport.Tests - (xmlReport.Failures + xmlReport.Errors + xmlReport.Skipped),
	}

	testReportRef := types.TestReportReference{
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

func convertTestSuites(xmlSuites []TestSuite, testReportRef types.TestReportReference) ([]types.TestSuite, []types.TestCase) {
	var (
		testSuites []types.TestSuite
		testCases  []types.TestCase
	)

	for i, ts := range xmlSuites {
		testSuiteId := fmt.Sprintf("%s-%d", testReportRef.Id, i+1)
		testSuite := types.TestSuite{
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

		testSuiteRef := types.TestSuiteReference{
			Id:         testSuite.Id,
			TestReport: testSuite.TestReport,
		}
		testSuiteCases := convertTestCases(ts.TestCases, testSuiteRef)

		// accumulate test case stats
		for _, tc := range testSuiteCases {
			testSuite.TotalCount++
			testSuite.TotalTime += tc.ExecutionTime
			switch tc.Status {
			case types.TestCaseStatusFailed:
				testSuite.FailedCount++
			case types.TestCaseStatusError:
				testSuite.ErrorCount++
			case types.TestCaseStatusSkipped:
				testSuite.SkippedCount++
			case types.TestCaseStatusSuccess:
				testSuite.SuccessCount++
			}
		}

		testSuites = append(testSuites, testSuite)
		testCases = append(testCases, testSuiteCases...)
	}

	return testSuites, testCases
}

func convertTestCases(xmlCases []TestCase, testSuiteRef types.TestSuiteReference) []types.TestCase {
	testCases := make([]types.TestCase, 0, len(xmlCases))

	for i, tc := range xmlCases {
		// see: https://gitlab.com/gitlab-org/gitlab/-/blob/master/lib/gitlab/ci/parsers/test/junit.rb#create_test_case

		var (
			status string
			output strings.Builder
		)

		switch {
		case tc.Failure != nil:
			status = types.TestCaseStatusFailed
			output.WriteString(tc.Failure.Output())
		case tc.Error != nil:
			status = types.TestCaseStatusError
			output.WriteString(tc.Error.Output())
		case tc.Skipped != nil:
			status = types.TestCaseStatusSkipped
			output.WriteString(tc.Skipped.Message)
		default:
			status = types.TestCaseStatusSuccess
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

		attachements := ParseTextAttachments(output.String())
		properties := convertTestProperties(append(tc.Properties, ParseTextProperties(output.String())...))

		testCaseId := fmt.Sprintf("%s-%d", testSuiteRef.Id, i+1)
		testCase := types.TestCase{
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

func convertTestProperties(xmlProps []Property) []types.TestProperty {
	if len(xmlProps) == 0 {
		return nil
	}

	properties := make([]types.TestProperty, 0, len(xmlProps))
	for _, p := range xmlProps {
		value := p.Value
		if value == "" {
			value = p.Text
		}
		properties = append(properties, types.TestProperty{Name: p.Name, Value: value})
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
