package junitxml_test

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/cluttrdev/gitlab-exporter/internal/junitxml"
	"github.com/cluttrdev/gitlab-exporter/internal/types"
	"github.com/google/go-cmp/cmp"
)

func TestConvertTestReport(t *testing.T) {
	job := types.JobReference{
		Id:   1337,
		Name: "tests-registration",
		Pipeline: types.PipelineReference{
			Id: 42,
		},
	}

	xmlReport := junitxml.TestReport{
		XMLName:   xml.Name{Local: "testsuites"},
		Tests:     8,
		Failures:  1,
		Errors:    1,
		Skipped:   1,
		Time:      16.082687,
		Timestamp: "2021-04-02T15:48:23",
		TestSuites: []junitxml.TestSuite{
			{
				Name:      "Tests.Registration",
				Tests:     8,
				Failures:  1,
				Errors:    1,
				Skipped:   1,
				Time:      16.082687,
				Timestamp: "2021-04-02T15:48:23",
				File:      "tests/registration.code",
				Properties: []junitxml.Property{
					{Name: "version", Value: "1.774"},
					{Name: "commit", Value: "ef7bebf"},
					{Name: "browser", Value: "Google Chrome"},
					{Name: "ci", Value: "https://github.com/actions/runs/1234"},
					{Name: "config", Text: "Config line #1"},
				},
				SystemOut: junitxml.SystemOut{Text: "Data written to standard out."},
				SystemErr: junitxml.SystemErr{Text: "Data written to standard error."},
				TestCases: []junitxml.TestCase{
					{Name: "testCase1", Classname: "Tests.Registration", Time: 2.436, File: "tests/registration.code", Line: 24},
					{Name: "testCase2", Classname: "Tests.Registration", Time: 1.534, File: "tests/registration.code", Line: 62},
					{Name: "testCase3", Classname: "Tests.Registration", Time: 0.822, File: "tests/registration.code", Line: 102},

					{
						Name: "testCase4", Classname: "Tests.Registration", Time: 0, File: "tests/registration.code", Line: 164,
						Skipped: &junitxml.Skipped{Message: "Test was skipped."},
					},
					{
						Name: "testCase5", Classname: "Tests.Registration", Time: 2.902412, File: "tests/registration.code", Line: 202,
						Failure: &junitxml.Failure{Message: "Expected value did not match.", Type: "AssertionError", Text: "Failure description or stack trace"},
					},
					{
						Name: "testCase6", Classname: "Tests.Registration", Time: 3.819, File: "tests/registration.code", Line: 235,
						Error: &junitxml.Error{Message: "Division by zero.", Type: "ArithmeticError", Text: "Error description or stack trace"},
					},
					{
						Name: "testCase7", Classname: "Tests.Registration", Time: 2.944, File: "tests/registration.code", Line: 287,
						SystemOut: &junitxml.SystemOut{Text: "Data written to standard out."},
						SystemErr: &junitxml.SystemErr{Text: "Data written to standard err."},
					},
					{
						Name: "testCase8", Classname: "Tests.Registration", Time: 1.625275, File: "tests/registration.code", Line: 302,
						Properties: []junitxml.Property{
							{Name: "priority", Value: "high"},
							{Name: "language", Value: "english"},
							{Name: "author", Value: "Adrian"},
							{Name: "attachment", Value: "screenshots/dashboard.png"},
							{Name: "attachment", Value: "screenshots/users.png"},
							{Name: "description", Text: "This text describes the purpose of this test case and provides an overview of what the test does and how it works."},
						},
					},
				},
			},
		},
	}

	testReport, testSuites, testCases := junitxml.ConvertTestReport(xmlReport, job)

	wantReport := types.TestReport{
		Id: "42-1337",
		Pipeline: types.PipelineReference{
			Id: 42,
		},

		TotalTime:    time.Duration(16.082687 * float64(time.Second)),
		TotalCount:   8,
		FailedCount:  1,
		ErrorCount:   1,
		SkippedCount: 1,
		SuccessCount: 5,
	}
	if diff := cmp.Diff(wantReport, testReport); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}

	wantTestSuites := []types.TestSuite{
		{
			Id: "42-1337-1",
			TestReport: types.TestReportReference{
				Id: "42-1337",
				Pipeline: types.PipelineReference{
					Id: 42,
				},
			},

			Name:         "Tests.Registration",
			TotalTime:    time.Duration(16.082687 * float64(time.Second)),
			TotalCount:   8,
			FailedCount:  1,
			ErrorCount:   1,
			SkippedCount: 1,
			SuccessCount: 5,

			Properties: []types.TestProperty{
				{Name: "version", Value: "1.774"},
				{Name: "commit", Value: "ef7bebf"},
				{Name: "browser", Value: "Google Chrome"},
				{Name: "ci", Value: "https://github.com/actions/runs/1234"},
				{Name: "config", Value: "Config line #1"},
			},
		},
	}
	if diff := cmp.Diff(wantTestSuites, testSuites); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}

	wantTestSuiteRef := types.TestSuiteReference{
		Id: "42-1337-1",
		TestReport: types.TestReportReference{
			Id: "42-1337",
			Pipeline: types.PipelineReference{
				Id: 42,
			},
		},
	}
	wantTestCases := []types.TestCase{
		{
			Id:        "42-1337-1-1",
			TestSuite: wantTestSuiteRef,

			Name:          "testCase1",
			Classname:     "Tests.Registration",
			ExecutionTime: time.Duration(2.436 * float64(time.Second)),
			File:          "tests/registration.code",
			Status:        "success",
		},
		{
			Id:        "42-1337-1-2",
			TestSuite: wantTestSuiteRef,

			Name:          "testCase2",
			Classname:     "Tests.Registration",
			ExecutionTime: time.Duration(1.534 * float64(time.Second)),
			File:          "tests/registration.code",
			Status:        "success",
		},
		{
			Id:        "42-1337-1-3",
			TestSuite: wantTestSuiteRef,

			Name:          "testCase3",
			Classname:     "Tests.Registration",
			ExecutionTime: time.Duration(0.822 * float64(time.Second)),
			File:          "tests/registration.code",
			Status:        "success",
		},

		{
			Id:        "42-1337-1-4",
			TestSuite: wantTestSuiteRef,

			Name:          "testCase4",
			Classname:     "Tests.Registration",
			ExecutionTime: 0,
			File:          "tests/registration.code",
			Status:        "skipped",
			SystemOutput:  "Test was skipped.",
		},
		{
			Id:        "42-1337-1-5",
			TestSuite: wantTestSuiteRef,

			Name:          "testCase5",
			Classname:     "Tests.Registration",
			ExecutionTime: time.Duration(2.902412 * float64(time.Second)),
			File:          "tests/registration.code",
			Status:        "failed",
			SystemOutput:  "AssertionError: Expected value did not match.\n\nFailure description or stack trace",
		},
		{
			Id:        "42-1337-1-6",
			TestSuite: wantTestSuiteRef,

			Name:          "testCase6",
			Classname:     "Tests.Registration",
			ExecutionTime: time.Duration(3.819 * float64(time.Second)),
			File:          "tests/registration.code",
			Status:        "error",
			SystemOutput:  "ArithmeticError: Division by zero.\n\nError description or stack trace",
		},
		{
			Id:        "42-1337-1-7",
			TestSuite: wantTestSuiteRef,

			Name:          "testCase7",
			Classname:     "Tests.Registration",
			ExecutionTime: time.Duration(2.944 * float64(time.Second)),
			File:          "tests/registration.code",
			Status:        "success",
			SystemOutput:  "System Out:\n\nData written to standard out.\n\nSystem Err:\n\nData written to standard err.",
		},
		{
			Id:        "42-1337-1-8",
			TestSuite: wantTestSuiteRef,

			Name:          "testCase8",
			Classname:     "Tests.Registration",
			ExecutionTime: time.Duration(1.625275 * float64(time.Second)),
			File:          "tests/registration.code",
			Status:        "success",

			Properties: []types.TestProperty{
				{Name: "priority", Value: "high"},
				{Name: "language", Value: "english"},
				{Name: "author", Value: "Adrian"},
				{Name: "attachment", Value: "screenshots/dashboard.png"},
				{Name: "attachment", Value: "screenshots/users.png"},
				{Name: "description", Value: "This text describes the purpose of this test case and provides an overview of what the test does and how it works."},
			},
		},
	}
	if diff := cmp.Diff(wantTestCases, testCases); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}
