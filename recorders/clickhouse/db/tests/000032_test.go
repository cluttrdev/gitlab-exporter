package migrations_test

import (
	"context"
	"fmt"
	"testing"
)

func prepareTestData000032(ctx context.Context, client *client) error {
	// insert jobs data
	jobsQuery := `
	INSERT INTO jobs (id, pipeline_id, finished_at) VALUES
		(1001001, 1001, 1.1),
		(1001002, 1001, 1.2),
		(1002001, 1002, 2.1),
		(1003001, 1003, 0.0)
	`
	if err := client.Conn.Exec(ctx, jobsQuery); err != nil {
		return fmt.Errorf("insert jobs data: %w", err)
	}
	// insert pipelines data
	pipelinesQuery := `
	INSERT INTO pipelines (id, finished_at) VALUES
		(1001, 1.9),
		(1002, 2.9),
		(1003, 3.9)
	`
	if err := client.Conn.Exec(ctx, pipelinesQuery); err != nil {
		return fmt.Errorf("insert pipelines data: %w", err)
	}
	// insert testcases data
	testcasesQuery := `
	INSERT INTO testcases (id, testsuite_id, testreport_id, job_id, pipeline_id, project_id, status, name) VALUES
		('1-1001001-1-1', '1-1001001-1', '1-1001001', 1001001, 1001, 1, 'success', 'test_case_1_1_1_a'),
		('1-1001001-1-2', '1-1001001-1', '1-1001001', 1001001, 1001, 1, 'failed',  'test_case_1_1_1_b'),
		('1-1001002-1-1', '1-1001002-1', '1-1001002', 1001002, 1001, 1, 'skipped', 'test_case_1_1_2_a'),
		('1-1002001-1-1', '1-1002001-1', '1-1002001', 1002001, 1002, 1, 'success', 'test_case_1_2_1_a'),
		('1-1003001-1-1', '1-1003001-1', '1-1003001', 1003001, 1003, 1, 'success', 'test_case_1_3_1_a'),
		('2-2001001-1-1', '2-2001001-1', '2-2001001', 2001001, 2001, 2, 'success', 'test_case_2_1_1_a')
	`
	if err := client.Conn.Exec(ctx, testcasesQuery); err != nil {
		return fmt.Errorf("insert testcases data: %w", err)
	}

	return nil
}

func Test000032Up(t *testing.T) {
	client := testClient(t)
	ctx := t.Context()

	// Prepare
	if err := client.Migration.Migrate(31); err != nil {
		t.Fatalf("failed to migrate to version 31: %v", err)
	}

	if err := prepareTestData000032(ctx, client); err != nil {
		t.Fatalf("failed to prepare test data: %v", err)
	}

	// Migrate
	if err := client.Migration.Migrate(32); err != nil {
		t.Fatalf("failed to migrate to version 32: %v", err)
	}

	// Check
	type result struct {
		ID              string `ch:"id"`
		JobID           int64  `ch:"job_id"`
		PipelineID      int64  `ch:"pipeline_id"`
		Status          string `ch:"status"`
		Name            string `ch:"name"`
		ReportCreatedAt uint32 `ch:"report_created_at"`
	}

	var results []result
	query := "SELECT id, job_id, pipeline_id, status, name, report_created_at FROM testcases"
	if err := client.Conn.Select(ctx, &results, query); err != nil {
		t.Fatalf("failed to query testcases: %v", err)
	}

	if len(results) != 6 {
		t.Fatalf("expected 6 testcase, got: %d", len(results))
	}

	tc := results[0]

	// Verify original data is preserved
	if tc.ID != "1-1001001-1-1" {
		t.Errorf("expected id = '1-1001001-1-1', got: %s", tc.ID)
	}
	if tc.JobID != 1001001 {
		t.Errorf("expected job_id = 1001001, got: %d", tc.JobID)
	}
	if tc.PipelineID != 1001 {
		t.Errorf("expected pipeline_id = 1001, got: %d", tc.PipelineID)
	}
	if tc.Status != "success" {
		t.Errorf("expected status = 'success', got: %s", tc.Status)
	}
	if tc.Name != "test_case_1_1_1_a" {
		t.Errorf("expected name = 'test_case_1_1_1_a', got: %s", tc.Name)
	}

	// Verify report_created_at was calculated correctly from (jobs|pipelines).finished_at

	// test_case_1_1_a, test_case_1_1_b, test_case_1_2_a -> job.finished_at = 1.x -> report_created_at = 1
	if results[0].ReportCreatedAt != 1 {
		t.Errorf("expected report_created_at = 1, got: %d", results[0].ReportCreatedAt)
	}
	if results[1].ReportCreatedAt != 1 {
		t.Errorf("expected report_created_at = 1, got: %d", results[1].ReportCreatedAt)
	}
	if results[2].ReportCreatedAt != 1 {
		t.Errorf("expected report_created_at = 1, got: %d", results[2].ReportCreatedAt)
	}
	// test_case_2_1_a -> job.finished_at = 2.x -> report_created_at = 2
	if results[3].ReportCreatedAt != 2 {
		t.Errorf("expected report_created_at = 2, got: %d", results[3].ReportCreatedAt)
	}
	// test_case_3_1_a -> job.finished_at = 0.0, pipeline.finished_at = 3.x -> report_created_at = 3
	if results[4].ReportCreatedAt != 3 {
		t.Errorf("expected report_created_at = 3, got: %d", results[4].ReportCreatedAt)
	}
	// test_case with no corresponding job/pipelines, left join -> job.finished_at = 0.0 and pipeline.finished_at = 0.0 -> report_created_at = 0
	if results[5].ReportCreatedAt != 0 {
		t.Errorf("expected report_created_at = 0, got: %d", results[5].ReportCreatedAt)
	}
}
