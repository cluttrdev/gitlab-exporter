package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

type TestReport struct {
	XMLName    xml.Name    `xml:"testsuites"`
	Tests      int64       `xml:"tests,attr"`
	Failures   int64       `xml:"failures,attr"`
	Errors     int64       `xml:"errors,attr"`
	Skipped    int64       `xml:"skipped,attr"`
	Time       float64     `xml:"time,attr"`
	Timestamp  string      `xml:"timestamp,attr"`
	TestSuites []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	XMLName    xml.Name   `xml:"testsuite"`
	Name       string     `xml:"name,attr"`
	Tests      int64      `xml:"tests,attr"`
	Failures   int64      `xml:"failures,attr"`
	Errors     int64      `xml:"errors,attr"`
	Skipped    int64      `xml:"skipped,attr"`
	Time       float64    `xml:"time,attr"`
	Timestamp  string     `xml:"timestamp,attr"`
	File       string     `xml:"file,attr"`
	Properties []Property `xml:"properties>property"`
	TestCases  []TestCase `xml:"testcase"`
}

type TestCase struct {
	XMLName   xml.Name `xml:"testcase"`
	Name      string   `xml:"name,attr"`
	Classname string   `xml:"classname,attr"`
	Tests     int64    `xml:"tests,attr"`
	Failures  int64    `xml:"failures,attr"`
	Errors    int64    `xml:"errors,attr"`
	Skipped   int64    `xml:"skipped,attr"`
	Time      float64  `xml:"time,attr"`
	Timestamp string   `xml:"timestamp,attr"`
	File      string   `xml:"file,attr"`
	Line      int64    `xml:"line,attr"`
}

type Property struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

func logEmbeddedMetrics(report TestReport) ([]string, error) {
	re := regexp.MustCompile(`[^a-zA-Z_]+`)

	var metrics []string
	for _, testsuite := range report.TestSuites {
		t, err := time.Parse(time.RFC3339, testsuite.Timestamp)
		if err != nil {
			fmt.Println(err)
			t = time.Now().UTC()
		}

		labelpairs := []string{
			fmt.Sprintf("name=\"%s\"", testsuite.Name),
		}
		for _, p := range testsuite.Properties {
			name := re.ReplaceAllString(p.Name, "_")
			labelpairs = append(labelpairs, fmt.Sprintf("%s=\"%s\"", name, p.Value))
		}
		labels := strings.Join(labelpairs, ",")

		metrics = append(metrics,
			fmt.Sprintf("METRIC_junit_testsuite_tests{%s} %d %d", labels, testsuite.Tests, t.UnixMilli()),
			fmt.Sprintf("METRIC_junit_testsuite_failures{%s} %d %d", labels, testsuite.Failures, t.UnixMilli()),
			fmt.Sprintf("METRIC_junit_testsuite_errors{%s} %d %d", labels, testsuite.Errors, t.UnixMilli()),
			fmt.Sprintf("METRIC_junit_testsuite_skipped{%s} %d %d", labels, testsuite.Skipped, t.UnixMilli()),
		)
	}

	return metrics, nil
}

func execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("not enough arguments")
	}

	data, err := os.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	var report TestReport
	if err := xml.Unmarshal(data, &report); err != nil {
		return fmt.Errorf("unmarshal data: %w", err)
	}

	metrics, err := logEmbeddedMetrics(report)
	if err != nil {
		return fmt.Errorf("parse data: %w", err)
	}

	for _, m := range metrics {
		fmt.Println(m)
	}
	return nil
}

func main() {
	if err := execute(os.Args[1:]); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
