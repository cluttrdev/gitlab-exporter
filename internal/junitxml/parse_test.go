package junitxml_test

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.cluttr.dev/gitlab-exporter/internal/junitxml"
)

const (
	dataBasic = `
        <?xml version="1.0" encoding="UTF-8"?>
        <testsuites time="15.682687">
            <testsuite name="Tests.Registration" time="6.605871">
                <testcase name="testCase1" classname="Tests.Registration" time="2.113871" />
                <testcase name="testCase2" classname="Tests.Registration" time="1.051" />
                <testcase name="testCase3" classname="Tests.Registration" time="3.441" />
            </testsuite>
            <testsuite name="Tests.Authentication" time="9.076816">
                <testcase name="testCase7" classname="Tests.Authentication" time="2.508" />
                <testcase name="testCase8" classname="Tests.Authentication" time="1.230816" />
                <testcase name="testCase9" classname="Tests.Authentication" time="0.982">
                    <failure message="Assertion error message" type="AssertionError">Call stack printed here</failure>            
                </testcase>
            </testsuite>
        </testsuites>
    `
)

func TestParse(t *testing.T) {
	report, err := junitxml.Parse(strings.NewReader(dataBasic))
	if err != nil {
		t.Errorf("error parsing data: %v", err)
	}

	want := junitxml.TestReport{
		XMLName: xml.Name{Local: "testsuites"},
		Time:    15.682687,
		TestSuites: []junitxml.TestSuite{
			{
				Name: "Tests.Registration",
				Time: 6.605871,
				TestCases: []junitxml.TestCase{
					{Name: "testCase1", Classname: "Tests.Registration", Time: 2.113871},
					{Name: "testCase2", Classname: "Tests.Registration", Time: 1.051},
					{Name: "testCase3", Classname: "Tests.Registration", Time: 3.441},
				},
			},
			{
				Name: "Tests.Authentication",
				Time: 9.076816,
				TestCases: []junitxml.TestCase{
					{Name: "testCase7", Classname: "Tests.Authentication", Time: 2.508},
					{Name: "testCase8", Classname: "Tests.Authentication", Time: 1.230816},
					{Name: "testCase9", Classname: "Tests.Authentication", Time: 0.982, Failure: &junitxml.Failure{
						Message: "Assertion error message", Type: "AssertionError", Text: `Call stack printed here`,
					}},
				},
			},
		},
	}

	if diff := cmp.Diff(want, report); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestParseFull(t *testing.T) {
	data := `
    <?xml version="1.0" encoding="UTF-8"?>

    <testsuites name="Test run" tests="8" failures="1" errors="1" skipped="1" 
        assertions="20" time="16.082687" timestamp="2021-04-02T15:48:23">

        <testsuite name="Tests.Registration" tests="8" failures="1" errors="1" skipped="1" assertions="20" time="16.082687" timestamp="2021-04-02T15:48:23" file="tests/registration.code">

            <properties>
                <property name="version" value="1.774" />
                <property name="commit" value="ef7bebf" />
                <property name="browser" value="Google Chrome" />
                <property name="ci" value="https://github.com/actions/runs/1234" />
                <property name="config">Config line #1</property>
            </properties>

            <system-out>Data written to standard out.</system-out>

            <system-err>Data written to standard error.</system-err>

            <testcase name="testCase1" classname="Tests.Registration" assertions="2" time="2.436" file="tests/registration.code" line="24" />
            <testcase name="testCase2" classname="Tests.Registration" assertions="6" time="1.534" file="tests/registration.code" line="62" />
            <testcase name="testCase3" classname="Tests.Registration" assertions="3" time="0.822" file="tests/registration.code" line="102" />
            
            <testcase name="testCase4" classname="Tests.Registration" assertions="0" time="0" file="tests/registration.code" line="164">
                <skipped message="Test was skipped." />
            </testcase>

            <testcase name="testCase5" classname="Tests.Registration" assertions="2" time="2.902412" file="tests/registration.code" line="202">
                <failure message="Expected value did not match." type="AssertionError">Failure description or stack trace</failure>
            </testcase>

            <testcase name="testCase6" classname="Tests.Registration" assertions="0" time="3.819" file="tests/registration.code" line="235">
                <error message="Division by zero." type="ArithmeticError">Error description or stack trace</error>
            </testcase>

            <testcase name="testCase7" classname="Tests.Registration" assertions="3" time="2.944" file="tests/registration.code" line="287">
                <system-out>Data written to standard out.</system-out>

                <system-err>Data written to standard error.</system-err>
            </testcase>

            <testcase name="testCase8" classname="Tests.Registration" assertions="4" time="1.625275" file="tests/registration.code" line="302">
                <properties>
                    <property name="priority" value="high" />
                    <property name="language" value="english" />
                    <property name="author" value="Adrian" />
                    <property name="attachment" value="screenshots/dashboard.png" />
                    <property name="attachment" value="screenshots/users.png" />
                    <property name="description">This text describes the purpose of this test case and provides an overview of what the test does and how it works.</property>
                </properties>
            </testcase>
        </testsuite>
    </testsuites>
    `

	report, err := junitxml.Parse(strings.NewReader(data))
	if err != nil {
		t.Errorf("error parsing data: %v", err)
	}

	want := junitxml.TestReport{
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
				SystemOut: &junitxml.SystemOut{Text: "Data written to standard out."},
				SystemErr: &junitxml.SystemErr{Text: "Data written to standard error."},
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
						SystemErr: &junitxml.SystemErr{Text: "Data written to standard error."},
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

	if diff := cmp.Diff(want, report); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestParseMany(t *testing.T) {
	data := `
    <testsuites time="2.113871">
        <testsuite name="Tests.Registration" time="2.113871">
            <testcase name="testCase1" classname="Tests.Registration" time="2.113871" />
        </testsuite>
    </testsuites>
    <testsuites time="2.508">
        <testsuite name="Tests.Authentication" time="2.508">
            <testcase name="testCase1" classname="Tests.Authentication" time="2.508" />
        </testsuite>
    </testsuites>
    `
	reports, err := junitxml.ParseMany(strings.NewReader(data))
	if err != nil {
		t.Errorf("error parsing data: %v", err)
	}

	want := []junitxml.TestReport{
		{
			XMLName: xml.Name{Local: "testsuites"},
			Time:    2.113871,
			TestSuites: []junitxml.TestSuite{
				{
					Name: "Tests.Registration",
					Time: 2.113871,
					TestCases: []junitxml.TestCase{
						{Name: "testCase1", Classname: "Tests.Registration", Time: 2.113871},
					},
				},
			},
		},
		{
			XMLName: xml.Name{Local: "testsuites"},
			Time:    2.508,
			TestSuites: []junitxml.TestSuite{
				{
					Name: "Tests.Authentication",
					Time: 2.508,
					TestCases: []junitxml.TestCase{
						{Name: "testCase1", Classname: "Tests.Authentication", Time: 2.508},
					},
				},
			},
		},
	}

	if diff := cmp.Diff(want, reports); diff != "" {
		t.Errorf("Mismatch (-want +got):\n%s", diff)
	}
}

func TestParseTextAttachments(t *testing.T) {
	data := `
    Output line #1
    Output line #2
    [[ATTACHMENT|screenshots/dashboard.png]]
    [[ATTACHMENT|screenshots/users.png]]
    `

	attachments := junitxml.ParseTextAttachments(data)

	want := []string{
		"screenshots/dashboard.png",
		"screenshots/users.png",
	}

	if diff := cmp.Diff(want, attachments); diff != "" {
		t.Errorf("Config mismatch (-want +got):\n%s", diff)
	}
}

func TestParseTextProperties(t *testing.T) {
	data := `
    Output line #1
    Output line #2

    [[PROPERTY|author=Adrian]]
    [[PROPERTY|language=english]]

    [[PROPERTY|browser-log]]
    Log line #1
    Log line #2
    Log line #3
    [[/PROPERTY]]
    `

	properties := junitxml.ParseTextProperties(data)

	want := []junitxml.Property{
		{Name: "author", Value: "Adrian"},
		{Name: "language", Value: "english"},
		// {Name: "browser-log", Text: `Log line #1\Log line #1\nnLog line #1\n`},
	}

	if diff := cmp.Diff(want, properties); diff != "" {
		t.Errorf("Config mismatch (-want +got):\n%s", diff)
	}
}
