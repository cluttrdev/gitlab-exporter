package junitxml

import "time"

func Merge(reports []TestReport) TestReport {
	var report TestReport

	var timestamp *time.Time
	for _, r := range reports {
		report.Tests += r.Tests
		report.Failures += r.Failures
		report.Errors += r.Errors
		report.Skipped += r.Skipped
		report.Time += r.Time

		t, err := time.Parse(time.RFC3339, r.Timestamp)
		if err == nil && (timestamp == nil || t.Before(*timestamp)) {
			report.Timestamp = r.Timestamp
			timestamp = &t
		}

		report.TestSuites = append(report.TestSuites, r.TestSuites...)
	}

	return report
}
