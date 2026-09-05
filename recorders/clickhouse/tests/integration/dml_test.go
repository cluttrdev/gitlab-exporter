package integration_tests

import (
	"context"
	"fmt"
	"testing"

	v1 "go.opentelemetry.io/proto/otlp/common/v1"
	resourcev1 "go.opentelemetry.io/proto/otlp/resource/v1"
	tracev1 "go.opentelemetry.io/proto/otlp/trace/v1"

	"go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/recorders/clickhouse/internal/clickhouse"
)

func TestIntegration_InsertPipelines(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Error(err)
	}

	data := []*typespb.Pipeline{
		{
			Id:  1082136862,
			Iid: 12,
			Project: &typespb.ProjectReference{
				Id: 50817395,
			},
			Status: "success", Source: "push", Ref: "main",
			Sha: "e860ecdc74aee9aab22f6336a705bab05634c0c3",
			Timestamps: &typespb.PipelineTimestamps{
				CreatedAt:  &timestamppb.Timestamp{Seconds: 1700690657, Nanos: 951000000},
				UpdatedAt:  &timestamppb.Timestamp{Seconds: 1700690933, Nanos: 886000000},
				StartedAt:  &timestamppb.Timestamp{Seconds: 1700690659, Nanos: 10000000},
				FinishedAt: &timestamppb.Timestamp{Seconds: 1700690933, Nanos: 875000000},
			},
			Duration:       &durationpb.Duration{Seconds: 273},
			QueuedDuration: &durationpb.Duration{Seconds: 1},
		},
	}

	n, err := clickhouse.InsertPipelines(client, context.Background(), data)
	if err != nil {
		t.Error(err)
	}

	if n != len(data) {
		t.Errorf("Inserted %d pipelines, expected: %d", n, len(data))
	}
}

func TestIntegration_InsertJobs(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Error(err)
	}

	data := []*typespb.Job{
		{
			Pipeline: &typespb.PipelineReference{
				Id: 1082136862,
				Project: &typespb.ProjectReference{
					Id: 50817395,
				},
			},
			Id: 5599404160, Name: "test", Ref: "main", Stage: "test", Status: "success",
			Timestamps: &typespb.JobTimestamps{
				CreatedAt:  &timestamppb.Timestamp{Seconds: 1700690657, Nanos: 999000000},
				StartedAt:  &timestamppb.Timestamp{Seconds: 1700690851, Nanos: 366000000},
				FinishedAt: &timestamppb.Timestamp{Seconds: 1700690933, Nanos: 765000000},
			},
			Duration:       &durationpb.Duration{Seconds: 82, Nanos: 399463000},
			QueuedDuration: &durationpb.Duration{Nanos: 359749000},
		},
		{
			Pipeline: nil,
			Id:       42,
		},
	}

	n, err := clickhouse.InsertJobs(client, context.Background(), data)
	if err == nil {
		t.Errorf("Expected error due to job without pipeline, got `nil`")
	} else if err.Error() != "job without pipeline: 42" {
		t.Errorf("Unexpected error: %v", err)
	}

	if n != len(data)-1 {
		t.Errorf("Inserted %d jobs, expected: %d", n, len(data)-1)
	}
}

func TestIntegration_InsertSections(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Error(err)
	}

	data := []*typespb.Section{
		{
			Id:   5599404160001,
			Name: "script",
			Job: &typespb.JobReference{
				Id: 5599404160,
				Pipeline: &typespb.PipelineReference{
					Id: 1082136862,
					Project: &typespb.ProjectReference{
						Id: 50817395,
					},
				},
			},
			StartedAt:  &timestamppb.Timestamp{Seconds: 1700690851, Nanos: 366000000},
			FinishedAt: &timestamppb.Timestamp{Seconds: 1700690933, Nanos: 765000000},
			Duration:   &durationpb.Duration{Seconds: 82, Nanos: 399463000},
		},
	}

	n, err := clickhouse.InsertSections(client, context.Background(), data)
	if err != nil {
		t.Error(err)
	}

	if n != len(data) {
		t.Errorf("Inserted %d sections, expected: %d", n, len(data))
	}
}

func TestIntegration_InsertTestCases(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Error(err)
	}

	testSuiteRef := &typespb.TestSuiteReference{
		Id: "6252785472",
		TestReport: &typespb.TestReportReference{
			Id: "1190130970",
			Job: &typespb.JobReference{
				Id:   0,
				Name: "",
				Pipeline: &typespb.PipelineReference{
					Id: 1190130970,
				},
			},
		},
	}

	data := []*typespb.TestCase{
		{Id: "6252785472-1", TestSuite: testSuiteRef},
		{Id: "6252785472-2", TestSuite: testSuiteRef},
		{Id: "6252785472-3", TestSuite: testSuiteRef},
		{Id: "6252785472-4", TestSuite: testSuiteRef},
		{Id: "6252785472-5", TestSuite: testSuiteRef},
		{Id: "6252785472-6", TestSuite: testSuiteRef},
		{Id: "6252785472-7", TestSuite: testSuiteRef},
		{Id: "6252785472-8", TestSuite: testSuiteRef},
		{Id: "6252785472-9", TestSuite: testSuiteRef},
		{Id: "6252785472-10", TestSuite: testSuiteRef},
	}

	n, err := clickhouse.InsertTestCases(client, context.Background(), data)
	if err != nil {
		t.Error(err)
	}

	if n != 10 {
		t.Errorf("Inserted %d testcases, expected: %d", n, 10)
	}
}

func TestIntegration_InsertTraces(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Error(err)
	}

	// var ts int64 = (4294967295 + 1) * 1e9 // uint32 overflow
	var ts int64 = (4294967295 + 0) * 1e9

	data := []*typespb.Trace{
		{
			Data: &tracev1.TracesData{
				ResourceSpans: []*tracev1.ResourceSpans{
					{
						Resource: &resourcev1.Resource{},
						ScopeSpans: []*tracev1.ScopeSpans{
							{
								Scope: &v1.InstrumentationScope{},
								Spans: []*tracev1.Span{
									{
										StartTimeUnixNano: uint64(ts),
										EndTimeUnixNano:   uint64(ts) + 1,
										Status:            &tracev1.Status{},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	n, err := clickhouse.InsertTraces(client, context.Background(), data)
	if err != nil {
		t.Error(err)
	}

	if n != len(data) {
		t.Errorf("Inserted %d traces, expected: %d", n, len(data))
	}
}

func TestIntegration_InsertRunners(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Error(err)
	}

	metadata := &servicepb.RecordRequestMetadata{
		FetchedAt: &timestamppb.Timestamp{Seconds: 1700690657, Nanos: 500000000},
	}

	data := []*typespb.Runner{
		{
			Id:          1,
			ShortSha:    "a1b2c3d4",
			Description: "Test Runner 1",
			RunnerType:  typespb.RunnerType_RUNNER_TYPE_INSTANCE,
			TagList:     []string{"docker", "linux"},
			Status:      typespb.RunnerStatus_RUNNER_STATUS_ONLINE,
			Flags: &typespb.RunnerFlags{
				Locked:       false,
				Paused:       false,
				RunProtected: false,
				RunUntagged:  true,
			},
			Timestamps: &typespb.RunnerTimestamps{
				CreatedAt:   &timestamppb.Timestamp{Seconds: 1700690000, Nanos: 0},
				ContactedAt: &timestamppb.Timestamp{Seconds: 1700690600, Nanos: 0},
			},
			CreatedBy: &typespb.UserReference{
				Id:       100,
				Username: "admin",
				Name:     "Admin User",
			},
		},
		{
			Id:          2,
			ShortSha:    "e5f6g7h8",
			Description: "Test Runner 2",
			RunnerType:  typespb.RunnerType_RUNNER_TYPE_PROJECT,
			TagList:     []string{"kubernetes", "prod"},
			Status:      typespb.RunnerStatus_RUNNER_STATUS_OFFLINE,
			Flags: &typespb.RunnerFlags{
				Locked:       true,
				Paused:       false,
				RunProtected: true,
				RunUntagged:  false,
			},
			Timestamps: &typespb.RunnerTimestamps{
				CreatedAt:   &timestamppb.Timestamp{Seconds: 1700680000, Nanos: 0},
				ContactedAt: &timestamppb.Timestamp{Seconds: 1700685000, Nanos: 0},
			},
			CreatedBy: &typespb.UserReference{
				Id:       200,
				Username: "devops",
				Name:     "DevOps Team",
			},
		},
	}

	n, err := clickhouse.InsertRunners(client, context.Background(), data, metadata)
	if err != nil {
		t.Error(err)
	}

	if n != len(data) {
		t.Errorf("Inserted %d runners, expected: %d", n, len(data))
	}
}

func TestIntegration_InsertRunners_PopulatesBothTables(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Error(err)
	}

	metadata := &servicepb.RecordRequestMetadata{
		FetchedAt: &timestamppb.Timestamp{Seconds: 1700690700, Nanos: 0},
	}

	data := []*typespb.Runner{
		{
			Id:          3,
			ShortSha:    "i9j0k1l2",
			Description: "Test Runner 3",
			RunnerType:  typespb.RunnerType_RUNNER_TYPE_GROUP,
			TagList:     []string{"shared"},
			Status:      typespb.RunnerStatus_RUNNER_STATUS_ONLINE,
			Flags: &typespb.RunnerFlags{
				Locked:       false,
				Paused:       false,
				RunProtected: false,
				RunUntagged:  true,
			},
			Timestamps: &typespb.RunnerTimestamps{
				CreatedAt:   &timestamppb.Timestamp{Seconds: 1700690000, Nanos: 0},
				ContactedAt: &timestamppb.Timestamp{Seconds: 1700690650, Nanos: 0},
			},
			CreatedBy: &typespb.UserReference{
				Id:       300,
				Username: "group-admin",
				Name:     "Group Admin",
			},
		},
	}

	n, err := clickhouse.InsertRunners(client, context.Background(), data, metadata)
	if err != nil {
		t.Error(err)
	}

	if n != len(data) {
		t.Errorf("Inserted %d runners, expected: %d", n, len(data))
	}

	// Wait for async insert to complete
	ctx := context.Background()
	if err := client.Exec(ctx, "SYSTEM FLUSH LOGS"); err != nil {
		t.Errorf("Failed to flush logs: %v", err)
	}

	// Check runners table (deduplicated)
	var countRunnersResult []struct {
		Count uint64 `ch:"count()"`
	}
	if err := client.Select(ctx, &countRunnersResult, "SELECT count(*) FROM runners WHERE id = 3"); err != nil {
		t.Errorf("Failed to query runners table: %v", err)
	}
	if len(countRunnersResult) == 0 || countRunnersResult[0].Count != 1 {
		count := uint64(0)
		if len(countRunnersResult) > 0 {
			count = countRunnersResult[0].Count
		}
		t.Errorf("Expected 1 row in runners table, got: %d", count)
	}

	// Check runners_raw table (all records)
	var countRunnersRawResult []struct {
		Count uint64 `ch:"count()"`
	}
	if err := client.Select(ctx, &countRunnersRawResult, "SELECT count(*) FROM _runners_raw WHERE id = 3"); err != nil {
		t.Errorf("Failed to query runners_raw table: %v", err)
	}
	if len(countRunnersRawResult) == 0 || countRunnersRawResult[0].Count != 1 {
		count := uint64(0)
		if len(countRunnersRawResult) > 0 {
			count = countRunnersRawResult[0].Count
		}
		t.Errorf("Expected 1 row in runners_raw table, got: %d", count)
	}
}

func TestIntegration_InsertRunners_Deduplication(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Error(err)
	}

	// Insert runner with first fetch
	metadata1 := &servicepb.RecordRequestMetadata{
		FetchedAt: &timestamppb.Timestamp{Seconds: 1700690800, Nanos: 0},
	}
	runner := &typespb.Runner{
		Id:          4,
		ShortSha:    "m3n4o5p6",
		Description: "Test Runner 4 - Version 1",
		RunnerType:  typespb.RunnerType_RUNNER_TYPE_INSTANCE,
		TagList:     []string{"test"},
		Status:      typespb.RunnerStatus_RUNNER_STATUS_ONLINE,
		Flags: &typespb.RunnerFlags{
			Locked:       false,
			Paused:       false,
			RunProtected: false,
			RunUntagged:  true,
		},
		Timestamps: &typespb.RunnerTimestamps{
			CreatedAt:   &timestamppb.Timestamp{Seconds: 1700690000, Nanos: 0},
			ContactedAt: &timestamppb.Timestamp{Seconds: 1700690750, Nanos: 0},
		},
		CreatedBy: &typespb.UserReference{
			Id:       400,
			Username: "test-user",
			Name:     "Test User",
		},
	}

	_, err = clickhouse.InsertRunners(client, context.Background(), []*typespb.Runner{runner}, metadata1)
	if err != nil {
		t.Error(err)
	}

	// Insert same runner with updated data and later fetch timestamp
	metadata2 := &servicepb.RecordRequestMetadata{
		FetchedAt: &timestamppb.Timestamp{Seconds: 1700690900, Nanos: 0},
	}
	runner.Description = "Test Runner 4 - Version 2"
	runner.Status = typespb.RunnerStatus_RUNNER_STATUS_OFFLINE

	_, err = clickhouse.InsertRunners(client, context.Background(), []*typespb.Runner{runner}, metadata2)
	if err != nil {
		t.Error(err)
	}

	// Wait for async insert and deduplication to complete
	ctx := context.Background()
	if err := client.Exec(ctx, "SYSTEM FLUSH LOGS"); err != nil {
		t.Errorf("Failed to flush logs: %v", err)
	}

	// Check runners_raw table - should have 2 records
	var countRunnersRawResult []struct {
		Count uint64 `ch:"count()"`
	}
	if err := client.Select(ctx, &countRunnersRawResult, "SELECT count(*) FROM _runners_raw WHERE id = 4"); err != nil {
		t.Errorf("Failed to query runners_raw table: %v", err)
	}
	if len(countRunnersRawResult) == 0 || countRunnersRawResult[0].Count != 2 {
		count := uint64(0)
		if len(countRunnersRawResult) > 0 {
			count = countRunnersRawResult[0].Count
		}
		t.Errorf("Expected 2 rows in runners_raw table (event log), got: %d", count)
	}

	// Check runners table - should have 1 record with latest data
	type runnerResult struct {
		Description string  `ch:"description"`
		Status      string  `ch:"status"`
		FetchedAt   float64 `ch:"_fetched_at"`
	}
	var results []runnerResult
	if err := client.Select(ctx, &results, "SELECT description, status, _fetched_at FROM runners FINAL WHERE id = 4"); err != nil {
		t.Errorf("Failed to query runners table: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("Expected 1 row in runners table after deduplication, got: %d", len(results))
	}

	// Verify we got the latest version
	if results[0].Description != "Test Runner 4 - Version 2" {
		t.Errorf("Expected description 'Test Runner 4 - Version 2', got: %s", results[0].Description)
	}
	if results[0].Status != "offline" {
		t.Errorf("Expected status 'offline', got: %s", results[0].Status)
	}
	expectedFetchedAt := float64(1700690900)
	if results[0].FetchedAt != expectedFetchedAt {
		t.Errorf("Expected _fetched_at %f, got: %f", expectedFetchedAt, results[0].FetchedAt)
	}
}

// ts is a small helper to build a protobuf timestamp from whole seconds.
func ts(seconds int64) *timestamppb.Timestamp {
	return &timestamppb.Timestamp{Seconds: seconds}
}

// assertDedup verifies that the deduplicating materialized view for a table kept
// exactly one row for the given id after out-of-order inserts, and that the
// surviving row carries the expected (latest) version value.
//
// It exercises the block-substitution assumption behind the bounded-JOIN views
// (migration 000035): the `<table>_in` reference inside the MV subquery must
// resolve to the inserted block, otherwise deduplication silently stops working.
func assertDedup(t *testing.T, client *clickhouse.Client, table, versionCol string, id int64, wantVersion float64) {
	t.Helper()

	ctx := context.Background()
	if err := client.Exec(ctx, "SYSTEM FLUSH LOGS"); err != nil {
		t.Fatalf("[%s] failed to flush logs: %v", table, err)
	}

	// FINAL collapses the ReplacingMergeTree so we assert on the deduplicated
	// state: exactly one row must survive with the latest version.
	var results []struct {
		Version float64 `ch:"v"`
	}
	query := fmt.Sprintf("SELECT %s AS v FROM %s FINAL WHERE id = %d", versionCol, table, id)
	if err := client.Select(ctx, &results, query); err != nil {
		t.Fatalf("[%s] failed to query table: %v", table, err)
	}

	if len(results) != 1 {
		t.Fatalf("[%s] expected 1 row after deduplication, got: %d", table, len(results))
	}
	if results[0].Version != wantVersion {
		t.Errorf("[%s] expected %s=%f, got: %f", table, versionCol, wantVersion, results[0].Version)
	}
}

func TestIntegration_InsertPipelines_Deduplication(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Fatal(err)
	}

	const id int64 = 900000001
	base := int64(2000000000)

	mk := func(updatedAt int64) []*typespb.Pipeline {
		return []*typespb.Pipeline{{
			Id:      id,
			Project: &typespb.ProjectReference{Id: 1},
			Timestamps: &typespb.PipelineTimestamps{
				CreatedAt: ts(base),
				UpdatedAt: ts(updatedAt),
			},
		}}
	}

	if _, err := clickhouse.InsertPipelines(client, context.Background(), mk(base)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertPipelines(client, context.Background(), mk(base-10)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertPipelines(client, context.Background(), mk(base+10)); err != nil {
		t.Fatal(err)
	}

	assertDedup(t, client, "pipelines", "updated_at", id, float64(base+10))
}

func TestIntegration_InsertProjects_Deduplication(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Fatal(err)
	}

	const id int64 = 900000002
	base := int64(2000000000)

	mk := func(lastActivityAt int64) []*typespb.Project {
		return []*typespb.Project{{
			Id:   id,
			Name: "dedup-project",
			Timestamps: &typespb.ProjectTimestamps{
				CreatedAt:      ts(base),
				UpdatedAt:      ts(base),
				LastActivityAt: ts(lastActivityAt),
			},
		}}
	}

	if _, err := clickhouse.InsertProjects(client, context.Background(), mk(base)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertProjects(client, context.Background(), mk(base-10)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertProjects(client, context.Background(), mk(base+10)); err != nil {
		t.Fatal(err)
	}

	assertDedup(t, client, "projects", "last_activity_at", id, float64(base+10))
}

func TestIntegration_InsertMergeRequests_Deduplication(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Fatal(err)
	}

	const id int64 = 900000003
	base := int64(2000000000)

	mk := func(updatedAt int64) []*typespb.MergeRequest {
		return []*typespb.MergeRequest{{
			Id:      id,
			Project: &typespb.ProjectReference{Id: 1},
			Timestamps: &typespb.MergeRequestTimestamps{
				CreatedAt: ts(base),
				UpdatedAt: ts(updatedAt),
			},
		}}
	}

	if _, err := clickhouse.InsertMergeRequests(client, context.Background(), mk(base)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertMergeRequests(client, context.Background(), mk(base-10)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertMergeRequests(client, context.Background(), mk(base+10)); err != nil {
		t.Fatal(err)
	}

	assertDedup(t, client, "mergerequests", "updated_at", id, float64(base+10))
}

func TestIntegration_InsertMergeRequestNoteEvents_Deduplication(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Fatal(err)
	}

	const id int64 = 900000004
	base := int64(2000000000)

	mk := func(updatedAt int64) []*typespb.MergeRequestNoteEvent {
		return []*typespb.MergeRequestNoteEvent{{
			Id: id,
			MergeRequest: &typespb.MergeRequestReference{
				Id:      1,
				Iid:     1,
				Project: &typespb.ProjectReference{Id: 1},
			},
			CreatedAt: ts(base),
			UpdatedAt: ts(updatedAt),
			Type:      "DiscussionNote",
		}}
	}

	if _, err := clickhouse.InsertMergeRequestNoteEvents(client, context.Background(), mk(base)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertMergeRequestNoteEvents(client, context.Background(), mk(base-10)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertMergeRequestNoteEvents(client, context.Background(), mk(base+10)); err != nil {
		t.Fatal(err)
	}

	assertDedup(t, client, "mergerequest_noteevents", "updated_at", id, float64(base+10))
}

func TestIntegration_InsertIssues_Deduplication(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Fatal(err)
	}

	// issues is ORDER BY (project_id, id); the bounded view filters on both key
	// columns, so this covers the two-column filter variant.
	const id int64 = 900000005
	base := int64(2000000000)

	mk := func(updatedAt int64) []*typespb.Issue {
		return []*typespb.Issue{{
			Id:      id,
			Iid:     1,
			Project: &typespb.ProjectReference{Id: 42},
			Timestamps: &typespb.IssueTimestamps{
				CreatedAt: ts(base),
				UpdatedAt: ts(updatedAt),
			},
			Title: "dedup-issue",
		}}
	}

	if _, err := clickhouse.InsertIssues(client, context.Background(), mk(base)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertIssues(client, context.Background(), mk(base-10)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertIssues(client, context.Background(), mk(base+10)); err != nil {
		t.Fatal(err)
	}

	assertDedup(t, client, "issues", "updated_at", id, float64(base+10))
}

func TestIntegration_InsertDeployments_Deduplication(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Fatal(err)
	}

	// deployments is ORDER BY (project_id, id); covers the two-column filter.
	const id int64 = 900000006
	base := int64(2000000000)

	mk := func(updatedAt int64) []*typespb.Deployment {
		return []*typespb.Deployment{{
			Id:  id,
			Iid: 1,
			Environment: &typespb.EnvironmentReference{
				Id:      1,
				Name:    "production",
				Project: &typespb.ProjectReference{Id: 42},
			},
			Timestamps: &typespb.DeploymentTimestamps{
				CreatedAt: ts(base),
				UpdatedAt: ts(updatedAt),
			},
		}}
	}

	if _, err := clickhouse.InsertDeployments(client, context.Background(), mk(base)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertDeployments(client, context.Background(), mk(base-10)); err != nil {
		t.Fatal(err)
	}
	if _, err := clickhouse.InsertDeployments(client, context.Background(), mk(base+10)); err != nil {
		t.Fatal(err)
	}

	assertDedup(t, client, "deployments", "updated_at", id, float64(base+10))
}

func TestIntegration_InsertRunners_BoundedDeduplication(t *testing.T) {
	client, err := GetTestClient(testSet)
	if err != nil {
		t.Fatal(err)
	}

	const id int64 = 900000007
	base := int64(2000000000)

	mk := func(fetchedAt int64) ([]*typespb.Runner, *servicepb.RecordRequestMetadata) {
		return []*typespb.Runner{{
				Id:          id,
				Description: "dedup-runner",
			}}, &servicepb.RecordRequestMetadata{
				FetchedAt: ts(fetchedAt),
			}
	}

	r1, m1 := mk(base)
	if _, err := clickhouse.InsertRunners(client, context.Background(), r1, m1); err != nil {
		t.Fatal(err)
	}
	r2, m2 := mk(base - 10)
	if _, err := clickhouse.InsertRunners(client, context.Background(), r2, m2); err != nil {
		t.Fatal(err)
	}
	r3, m3 := mk(base + 10)
	if _, err := clickhouse.InsertRunners(client, context.Background(), r3, m3); err != nil {
		t.Fatal(err)
	}

	assertDedup(t, client, "runners", "_fetched_at", id, float64(base+10))
}
