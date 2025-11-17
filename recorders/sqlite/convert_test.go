package sqlite

import (
	"encoding/json"
	"testing"

	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func TestConvertProject(t *testing.T) {
	msg := &typespb.Project{
		Id: 123,
		Namespace: &typespb.NamespaceReference{
			Id:       456,
			FullPath: "group/subgroup",
		},
		Name:     "test-project",
		FullPath: "group/subgroup/test-project",
	}

	result, err := ConvertProject(msg)
	if err != nil {
		t.Fatalf("ConvertProject() error = %v", err)
	}

	if result.Id != 123 {
		t.Errorf("Id = %d, want 123", result.Id)
	}
	if result.NamespaceId != 456 {
		t.Errorf("NamespaceId = %d, want 456", result.NamespaceId)
	}
	if len(result.Data) == 0 {
		t.Error("Data should not be empty")
	}

	// Verify data can be unmarshaled back
	var unmarshaled typespb.Project
	if err := json.Unmarshal(result.Data, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal Data: %v", err)
	}
}

func TestConvertPipeline(t *testing.T) {
	msg := &typespb.Pipeline{
		Id:  789,
		Iid: 42,
		Project: &typespb.ProjectReference{
			Id:       123,
			FullPath: "group/project",
		},
		Ref:    "main",
		Status: "success",
	}

	result, err := ConvertPipeline(msg)
	if err != nil {
		t.Fatalf("ConvertPipeline() error = %v", err)
	}

	if result.Id != 789 {
		t.Errorf("Id = %d, want 789", result.Id)
	}
	if result.Iid != 42 {
		t.Errorf("Iid = %d, want 42", result.Iid)
	}
	if result.ProjectId != 123 {
		t.Errorf("ProjectId = %d, want 123", result.ProjectId)
	}
}

func TestConvertJob(t *testing.T) {
	msg := &typespb.Job{
		Id:   999,
		Name: "build",
		Pipeline: &typespb.PipelineReference{
			Id:  789,
			Iid: 42,
			Project: &typespb.ProjectReference{
				Id:       123,
				FullPath: "group/project",
			},
		},
		Status: "success",
	}

	result, err := ConvertJob(msg)
	if err != nil {
		t.Fatalf("ConvertJob() error = %v", err)
	}

	if result.Id != 999 {
		t.Errorf("Id = %d, want 999", result.Id)
	}
	if result.PipelineId != 789 {
		t.Errorf("PipelineId = %d, want 789", result.PipelineId)
	}
	if result.ProjectId != 123 {
		t.Errorf("ProjectId = %d, want 123", result.ProjectId)
	}
}

func TestConvertSection(t *testing.T) {
	msg := &typespb.Section{
		Id:   111,
		Name: "test_section",
		Job: &typespb.JobReference{
			Id:   999,
			Name: "build",
			Pipeline: &typespb.PipelineReference{
				Id:  789,
				Iid: 42,
				Project: &typespb.ProjectReference{
					Id: 123,
				},
			},
		},
	}

	result, err := ConvertSection(msg)
	if err != nil {
		t.Fatalf("ConvertSection() error = %v", err)
	}

	if result.Id != 111 {
		t.Errorf("Id = %d, want 111", result.Id)
	}
	if result.JobId != 999 {
		t.Errorf("JobId = %d, want 999", result.JobId)
	}
	if result.PipelineId != 789 {
		t.Errorf("PipelineId = %d, want 789", result.PipelineId)
	}
	if result.ProjectId != 123 {
		t.Errorf("ProjectId = %d, want 123", result.ProjectId)
	}
}

func TestConvertMetric(t *testing.T) {
	msg := &typespb.Metric{
		Id:   []byte("metric-id-123"),
		Iid:  42,
		Name: "test_metric",
		Job: &typespb.JobReference{
			Id: 999,
			Pipeline: &typespb.PipelineReference{
				Id: 789,
				Project: &typespb.ProjectReference{
					Id: 123,
				},
			},
		},
		Value: 3.14,
	}

	result, err := ConvertMetric(msg)
	if err != nil {
		t.Fatalf("ConvertMetric() error = %v", err)
	}

	if result.Id != "metric-id-123" {
		t.Errorf("Id = %s, want metric-id-123", result.Id)
	}
	if result.Iid != 42 {
		t.Errorf("Iid = %d, want 42", result.Iid)
	}
	if result.JobId != 999 {
		t.Errorf("JobId = %d, want 999", result.JobId)
	}
}

func TestConvertIssue(t *testing.T) {
	msg := &typespb.Issue{
		Id:    555,
		Iid:   77,
		Title: "Test Issue",
		Project: &typespb.ProjectReference{
			Id: 123,
		},
		State: typespb.IssueState_ISSUE_STATE_OPENED,
	}

	result, err := ConvertIssue(msg)
	if err != nil {
		t.Fatalf("ConvertIssue() error = %v", err)
	}

	if result.Id != 555 {
		t.Errorf("Id = %d, want 555", result.Id)
	}
	if result.Iid != 77 {
		t.Errorf("Iid = %d, want 77", result.Iid)
	}
	if result.ProjectId != 123 {
		t.Errorf("ProjectId = %d, want 123", result.ProjectId)
	}
}

func TestConvertMergeRequest(t *testing.T) {
	msg := &typespb.MergeRequest{
		Id:    666,
		Iid:   88,
		Title: "Test MR",
		Project: &typespb.ProjectReference{
			Id: 123,
		},
		State: "opened",
	}

	result, err := ConvertMergeRequest(msg)
	if err != nil {
		t.Fatalf("ConvertMergeRequest() error = %v", err)
	}

	if result.Id != 666 {
		t.Errorf("Id = %d, want 666", result.Id)
	}
	if result.Iid != 88 {
		t.Errorf("Iid = %d, want 88", result.Iid)
	}
	if result.ProjectId != 123 {
		t.Errorf("ProjectId = %d, want 123", result.ProjectId)
	}
}

func TestConvertMergeRequestNoteEvent(t *testing.T) {
	msg := &typespb.MergeRequestNoteEvent{
		Id: 777,
		MergeRequest: &typespb.MergeRequestReference{
			Id:  666,
			Iid: 88,
			Project: &typespb.ProjectReference{
				Id: 123,
			},
		},
		Type: "note",
	}

	result, err := ConvertMergeRequestNoteEvent(msg)
	if err != nil {
		t.Fatalf("ConvertMergeRequestNoteEvent() error = %v", err)
	}

	if result.Id != 777 {
		t.Errorf("Id = %d, want 777", result.Id)
	}
	if result.MergeRequestId != 666 {
		t.Errorf("MergeRequestId = %d, want 666", result.MergeRequestId)
	}
	if result.MergeRequestIid != 88 {
		t.Errorf("MergeRequestIid = %d, want 88", result.MergeRequestIid)
	}
	if result.MergeRequestProjectId != 123 {
		t.Errorf("MergeRequestProjectId = %d, want 123", result.MergeRequestProjectId)
	}
}

func TestConvertRunner(t *testing.T) {
	msg := &typespb.Runner{
		Id:          888,
		ShortSha:    "abc123",
		Description: "Test Runner",
	}

	result, err := ConvertRunner(msg)
	if err != nil {
		t.Fatalf("ConvertRunner() error = %v", err)
	}

	if result.Id != 888 {
		t.Errorf("Id = %d, want 888", result.Id)
	}
	if len(result.Data) == 0 {
		t.Error("Data should not be empty")
	}
}

func TestConvertDeployment(t *testing.T) {
	msg := &typespb.Deployment{
		Id:  333,
		Iid: 22,
		Environment: &typespb.EnvironmentReference{
			Id:   444,
			Name: "production",
		},
		Job: &typespb.JobReference{
			Id: 999,
			Pipeline: &typespb.PipelineReference{
				Id: 789,
				Project: &typespb.ProjectReference{
					Id: 123,
				},
			},
		},
	}

	result, err := ConvertDeployment(msg)
	if err != nil {
		t.Fatalf("ConvertDeployment() error = %v", err)
	}

	if result.Id != 333 {
		t.Errorf("Id = %d, want 333", result.Id)
	}
	if result.Iid != 22 {
		t.Errorf("Iid = %d, want 22", result.Iid)
	}
	if result.EnvironmentId != 444 {
		t.Errorf("EnvironmentId = %d, want 444", result.EnvironmentId)
	}
	if result.JobId != 999 {
		t.Errorf("JobId = %d, want 999", result.JobId)
	}
}

func TestConvertTestReport(t *testing.T) {
	msg := &typespb.TestReport{
		Id: "test-report-1",
		Job: &typespb.JobReference{
			Id: 999,
			Pipeline: &typespb.PipelineReference{
				Id: 789,
				Project: &typespb.ProjectReference{
					Id: 123,
				},
			},
		},
		TotalCount:   100,
		SuccessCount: 95,
		FailedCount:  5,
	}

	result, err := ConvertTestReport(msg)
	if err != nil {
		t.Fatalf("ConvertTestReport() error = %v", err)
	}

	if result.Id != "test-report-1" {
		t.Errorf("Id = %s, want test-report-1", result.Id)
	}
	if result.JobId != 999 {
		t.Errorf("JobId = %d, want 999", result.JobId)
	}
	if result.PipelineId != 789 {
		t.Errorf("PipelineId = %d, want 789", result.PipelineId)
	}
	if result.ProjectId != 123 {
		t.Errorf("ProjectId = %d, want 123", result.ProjectId)
	}
}

func TestConvertTestSuite(t *testing.T) {
	msg := &typespb.TestSuite{
		Id:   "test-suite-1",
		Name: "Unit Tests",
		TestReport: &typespb.TestReportReference{
			Id: "test-report-1",
			Job: &typespb.JobReference{
				Id: 999,
				Pipeline: &typespb.PipelineReference{
					Id: 789,
					Project: &typespb.ProjectReference{
						Id: 123,
					},
				},
			},
		},
	}

	result, err := ConvertTestSuite(msg)
	if err != nil {
		t.Fatalf("ConvertTestSuite() error = %v", err)
	}

	if result.Id != "test-suite-1" {
		t.Errorf("Id = %s, want test-suite-1", result.Id)
	}
	if result.TestReportId != "test-report-1" {
		t.Errorf("ReportId = %s, want test-report-1", result.TestReportId)
	}
	if result.JobId != 999 {
		t.Errorf("JobId = %d, want 999", result.JobId)
	}
}

func TestConvertTestCase(t *testing.T) {
	msg := &typespb.TestCase{
		Id:        "test-case-1",
		Name:      "test_addition",
		Classname: "MathTests",
		Status:    "success",
		TestSuite: &typespb.TestSuiteReference{
			Id: "test-suite-1",
			TestReport: &typespb.TestReportReference{
				Id: "test-report-1",
				Job: &typespb.JobReference{
					Id: 999,
					Pipeline: &typespb.PipelineReference{
						Id: 789,
						Project: &typespb.ProjectReference{
							Id: 123,
						},
					},
				},
			},
		},
	}

	result, err := ConvertTestCase(msg)
	if err != nil {
		t.Fatalf("ConvertTestCase() error = %v", err)
	}

	if result.Id != "test-case-1" {
		t.Errorf("Id = %s, want test-case-1", result.Id)
	}
	if result.TestSuiteId != "test-suite-1" {
		t.Errorf("SuiteId = %s, want test-suite-1", result.TestSuiteId)
	}
	if result.TestReportId != "test-report-1" {
		t.Errorf("ReportId = %s, want test-report-1", result.TestReportId)
	}
}

func TestConvertCoverageReport(t *testing.T) {
	msg := &typespb.CoverageReport{
		Id: "coverage-report-1",
		Job: &typespb.JobReference{
			Id: 999,
			Pipeline: &typespb.PipelineReference{
				Id: 789,
				Project: &typespb.ProjectReference{
					Id: 123,
				},
			},
		},
	}

	result, err := ConvertCoverageReport(msg)
	if err != nil {
		t.Fatalf("ConvertCoverageReport() error = %v", err)
	}

	if result.Id != "coverage-report-1" {
		t.Errorf("Id = %s, want coverage-report-1", result.Id)
	}
	if result.JobId != 999 {
		t.Errorf("JobId = %d, want 999", result.JobId)
	}
}

func TestConvertCoveragePackage(t *testing.T) {
	msg := &typespb.CoveragePackage{
		Id:   "package-1",
		Name: "com.example",
		Report: &typespb.CoverageReportReference{
			Id: "coverage-report-1",
			Job: &typespb.JobReference{
				Id: 999,
				Pipeline: &typespb.PipelineReference{
					Id: 789,
					Project: &typespb.ProjectReference{
						Id: 123,
					},
				},
			},
		},
	}

	result, err := ConvertCoveragePackage(msg)
	if err != nil {
		t.Fatalf("ConvertCoveragePackage() error = %v", err)
	}

	if result.Id != "package-1" {
		t.Errorf("Id = %s, want package-1", result.Id)
	}
	if result.ReportId != "coverage-report-1" {
		t.Errorf("ReportId = %s, want coverage-report-1", result.ReportId)
	}
}

func TestConvertCoverageClass(t *testing.T) {
	msg := &typespb.CoverageClass{
		Id:   "class-1",
		Name: "Calculator",
		Package: &typespb.CoveragePackageReference{
			Id:   "package-1",
			Name: "com.example",
			Report: &typespb.CoverageReportReference{
				Id: "coverage-report-1",
				Job: &typespb.JobReference{
					Id: 999,
					Pipeline: &typespb.PipelineReference{
						Id: 789,
						Project: &typespb.ProjectReference{
							Id: 123,
						},
					},
				},
			},
		},
	}

	result, err := ConvertCoverageClass(msg)
	if err != nil {
		t.Fatalf("ConvertCoverageClass() error = %v", err)
	}

	if result.Id != "class-1" {
		t.Errorf("Id = %s, want class-1", result.Id)
	}
	if result.PackageId != "package-1" {
		t.Errorf("PackageId = %s, want package-1", result.PackageId)
	}
	if result.ReportId != "coverage-report-1" {
		t.Errorf("ReportId = %s, want coverage-report-1", result.ReportId)
	}
}

func TestConvertCoverageMethod(t *testing.T) {
	msg := &typespb.CoverageMethod{
		Id: "method-1",
		Class: &typespb.CoverageClassReference{
			Id:   "class-1",
			Name: "Calculator",
			Package: &typespb.CoveragePackageReference{
				Id:   "package-1",
				Name: "com.example",
				Report: &typespb.CoverageReportReference{
					Id: "coverage-report-1",
					Job: &typespb.JobReference{
						Id: 999,
						Pipeline: &typespb.PipelineReference{
							Id: 789,
							Project: &typespb.ProjectReference{
								Id: 123,
							},
						},
					},
				},
			},
		},
	}

	result, err := ConvertCoverageMethod(msg)
	if err != nil {
		t.Fatalf("ConvertCoverageMethod() error = %v", err)
	}

	if result.Id != "method-1" {
		t.Errorf("Id = %s, want method-1", result.Id)
	}
	if result.ClassId != "class-1" {
		t.Errorf("ClassId = %s, want class-1", result.ClassId)
	}
	if result.PackageId != "package-1" {
		t.Errorf("PackageId = %s, want package-1", result.PackageId)
	}
	if result.ReportId != "coverage-report-1" {
		t.Errorf("ReportId = %s, want coverage-report-1", result.ReportId)
	}
}
