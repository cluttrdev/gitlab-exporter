package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

// setupTestDB creates an in-memory database for testing
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	if err := RunMigrations(t.Context(), db, "gitlab_ci"); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	return db
}

func TestRecorder_RecordProjects(t *testing.T) {
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	r := &Recorder{db: db}

	req := &servicepb.RecordProjectsRequest{
		Data: []*typespb.Project{
			{
				Id:       123,
				Name:     "test-project",
				FullPath: "group/test-project",
				Namespace: &typespb.NamespaceReference{
					Id:       456,
					FullPath: "group",
				},
			},
		},
	}

	summary, err := r.RecordProjects(context.Background(), req)
	if err != nil {
		t.Fatalf("RecordProjects() error = %v", err)
	}

	if summary.RecordedCount != 1 {
		t.Errorf("RecordedCount = %d, want 1", summary.RecordedCount)
	}

	// Verify data was inserted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query count: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 row in projects table, got %d", count)
	}
}

func TestRecorder_RecordPipelines(t *testing.T) {
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	r := &Recorder{db: db}

	req := &servicepb.RecordPipelinesRequest{
		Data: []*typespb.Pipeline{
			{
				Id:  789,
				Iid: 42,
				Project: &typespb.ProjectReference{
					Id:       123,
					FullPath: "group/project",
				},
				Ref:    "main",
				Status: "success",
			},
		},
	}

	summary, err := r.RecordPipelines(context.Background(), req)
	if err != nil {
		t.Fatalf("RecordPipelines() error = %v", err)
	}

	if summary.RecordedCount != 1 {
		t.Errorf("RecordedCount = %d, want 1", summary.RecordedCount)
	}
}

func TestRecorder_RecordJobs(t *testing.T) {
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	r := &Recorder{db: db}

	req := &servicepb.RecordJobsRequest{
		Data: []*typespb.Job{
			{
				Id:   999,
				Name: "build",
				Pipeline: &typespb.PipelineReference{
					Id:  789,
					Iid: 42,
					Project: &typespb.ProjectReference{
						Id: 123,
					},
				},
				Status: "success",
			},
		},
	}

	summary, err := r.RecordJobs(context.Background(), req)
	if err != nil {
		t.Fatalf("RecordJobs() error = %v", err)
	}

	if summary.RecordedCount != 1 {
		t.Errorf("RecordedCount = %d, want 1", summary.RecordedCount)
	}
}

func TestRecorder_RecordMultipleItems(t *testing.T) {
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	r := &Recorder{db: db}

	// Test recording multiple projects at once
	req := &servicepb.RecordProjectsRequest{
		Data: []*typespb.Project{
			{
				Id:   1,
				Name: "project-1",
				Namespace: &typespb.NamespaceReference{
					Id: 10,
				},
			},
			{
				Id:   2,
				Name: "project-2",
				Namespace: &typespb.NamespaceReference{
					Id: 20,
				},
			},
			{
				Id:   3,
				Name: "project-3",
				Namespace: &typespb.NamespaceReference{
					Id: 30,
				},
			},
		},
	}

	summary, err := r.RecordProjects(context.Background(), req)
	if err != nil {
		t.Fatalf("RecordProjects() error = %v", err)
	}

	if summary.RecordedCount != 3 {
		t.Errorf("RecordedCount = %d, want 3", summary.RecordedCount)
	}

	// Verify data was inserted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query count: %v", err)
	}

	if count != 3 {
		t.Errorf("Expected 3 rows in projects table, got %d", count)
	}
}

func TestRecorder_RecordIssues(t *testing.T) {
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	r := &Recorder{db: db}

	req := &servicepb.RecordIssuesRequest{
		Data: []*typespb.Issue{
			{
				Id:    555,
				Iid:   77,
				Title: "Test Issue",
				Project: &typespb.ProjectReference{
					Id: 123,
				},
			},
		},
	}

	summary, err := r.RecordIssues(context.Background(), req)
	if err != nil {
		t.Fatalf("RecordIssues() error = %v", err)
	}

	if summary.RecordedCount != 1 {
		t.Errorf("RecordedCount = %d, want 1", summary.RecordedCount)
	}
}

func TestRecorder_RecordMergeRequests(t *testing.T) {
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	r := &Recorder{db: db}

	req := &servicepb.RecordMergeRequestsRequest{
		Data: []*typespb.MergeRequest{
			{
				Id:    666,
				Iid:   88,
				Title: "Test MR",
				Project: &typespb.ProjectReference{
					Id: 123,
				},
			},
		},
	}

	summary, err := r.RecordMergeRequests(context.Background(), req)
	if err != nil {
		t.Fatalf("RecordMergeRequests() error = %v", err)
	}

	if summary.RecordedCount != 1 {
		t.Errorf("RecordedCount = %d, want 1", summary.RecordedCount)
	}
}

func TestRecorder_RecordRunners(t *testing.T) {
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	r := &Recorder{db: db}

	req := &servicepb.RecordRunnersRequest{
		Data: []*typespb.Runner{
			{
				Id:          888,
				ShortSha:    "abc123",
				Description: "Test Runner",
			},
		},
	}

	summary, err := r.RecordRunners(context.Background(), req)
	if err != nil {
		t.Fatalf("RecordRunners() error = %v", err)
	}

	if summary.RecordedCount != 1 {
		t.Errorf("RecordedCount = %d, want 1", summary.RecordedCount)
	}
}

func TestRecorder_RecordTestReports(t *testing.T) {
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	r := &Recorder{db: db}

	req := &servicepb.RecordTestReportsRequest{
		Data: []*typespb.TestReport{
			{
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

	summary, err := r.RecordTestReports(context.Background(), req)
	if err != nil {
		t.Fatalf("RecordTestReports() error = %v", err)
	}

	if summary.RecordedCount != 1 {
		t.Errorf("RecordedCount = %d, want 1", summary.RecordedCount)
	}
}

func TestRecorder_RecordCoverageReports(t *testing.T) {
	db := setupTestDB(t)
	defer func() { _ = db.Close() }()

	r := &Recorder{db: db}

	req := &servicepb.RecordCoverageReportsRequest{
		Data: []*typespb.CoverageReport{
			{
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

	summary, err := r.RecordCoverageReports(context.Background(), req)
	if err != nil {
		t.Fatalf("RecordCoverageReports() error = %v", err)
	}

	if summary.RecordedCount != 1 {
		t.Errorf("RecordedCount = %d, want 1", summary.RecordedCount)
	}
}

func TestNumFields(t *testing.T) {
	tests := []struct {
		name      string
		input     interface{}
		wantCount int
	}{
		{
			name: "Project struct",
			input: Project{
				Id:          123,
				NamespaceId: 456,
				Data:        []byte("test"),
			},
			wantCount: 3,
		},
		{
			name: "Pipeline struct",
			input: Pipeline{
				Id:        789,
				Iid:       42,
				ProjectId: 123,
				Data:      []byte("test"),
			},
			wantCount: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count := numFields(tt.input)
			if count != tt.wantCount {
				t.Errorf("numFields() = %d, want %d", count, tt.wantCount)
			}
		})
	}
}
