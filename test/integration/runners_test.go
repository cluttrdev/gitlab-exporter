package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/tasks"
)

func TestRunnersPipeline(t *testing.T) {
	mux, glab := setupGitLab(t)
	exp, rec := setupExporter(t)

	// Mock the GraphQL endpoint for runners (GraphQL client uses GET with query params)
	mux.HandleFunc("/api/graphql/", func(w http.ResponseWriter, r *http.Request) {
		// GraphQL client uses GET requests
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Return mock runners response
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"runners": map[string]interface{}{
					"nodes": []map[string]interface{}{
						{
							"id":          "gid://gitlab/Ci::Runner/100",
							"shortSha":    "abc12345",
							"description": "Integration Test Runner",
							"runnerType":  "INSTANCE_TYPE",
							"tagList":     []string{"docker", "linux", "test"},
							"status":      "ONLINE",
							"active":      true,
							"locked":      false,
							"paused":      false,
							"accessLevel": "NOT_PROTECTED",
							"runUntagged": true,
							"contactedAt": time.Now().Format(time.RFC3339),
							"createdAt":   time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
							"createdBy": map[string]interface{}{
								"id":       "gid://gitlab/User/1",
								"username": "test-admin",
								"name":     "Test Admin",
							},
						},
						{
							"id":          "gid://gitlab/Ci::Runner/200",
							"shortSha":    "def67890",
							"description": "Project Runner",
							"runnerType":  "PROJECT_TYPE",
							"tagList":     []string{"kubernetes"},
							"status":      "STALE",
							"active":      true,
							"locked":      true,
							"paused":      false,
							"accessLevel": "REF_PROTECTED",
							"runUntagged": false,
							"contactedAt": time.Now().Add(-8 * 24 * time.Hour).Format(time.RFC3339),
							"createdAt":   time.Now().Add(-30 * 24 * time.Hour).Format(time.RFC3339),
							"createdBy":   nil,
						},
						{
							"id":          "gid://gitlab/Ci::Runner/300",
							"shortSha":    "ghi11111",
							"description": "Never Contacted Runner",
							"runnerType":  "GROUP_TYPE",
							"tagList":     []string{},
							"status":      "NEVER_CONTACTED",
							"active":      false,
							"locked":      false,
							"paused":      true,
							"accessLevel": "NOT_PROTECTED",
							"runUntagged": false,
							"contactedAt": nil,
							"createdAt":   time.Now().Add(-1 * time.Hour).Format(time.RFC3339),
							"createdBy": map[string]interface{}{
								"id":       "gid://gitlab/User/2",
								"username": "project-owner",
								"name":     "Project Owner",
							},
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
						"endCursor":   nil,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	})

	// Test the full pipeline: Fetch -> Export
	ctx := context.Background()

	// Fetch runners
	fetchedAt := time.Now().UTC()
	runners, err := tasks.FetchRunners(ctx, glab)
	if err != nil {
		t.Fatalf("FetchRunners failed: %v", err)
	}

	// Verify we got all runners
	if len(runners) != 3 {
		t.Fatalf("expected 3 runners, got %d", len(runners))
	}

	// Verify first runner details
	if runners[0].Id != 100 {
		t.Errorf("expected runner ID 100, got %d", runners[0].Id)
	}
	if runners[0].ShortSha != "abc12345" {
		t.Errorf("expected ShortSha 'abc12345', got %s", runners[0].ShortSha)
	}
	if runners[0].Description != "Integration Test Runner" {
		t.Errorf("expected description 'Integration Test Runner', got %s", runners[0].Description)
	}
	if len(runners[0].TagList) != 3 {
		t.Errorf("expected 3 tags, got %d", len(runners[0].TagList))
	}
	if runners[0].CreatedBy.Username != "test-admin" {
		t.Errorf("expected creator username 'test-admin', got %s", runners[0].CreatedBy.Username)
	}

	// Verify second runner (with different settings)
	if runners[1].Id != 200 {
		t.Errorf("expected runner ID 200, got %d", runners[1].Id)
	}
	if !runners[1].Locked {
		t.Error("expected runner 200 to be locked")
	}
	if runners[1].CreatedBy.Id != 0 {
		t.Error("expected runner 200 to have no creator (should be zero)")
	}

	// Verify third runner (never contacted)
	if runners[2].Id != 300 {
		t.Errorf("expected runner ID 300, got %d", runners[2].Id)
	}
	if runners[2].ContactedAt != nil {
		t.Error("expected runner 300 to have nil ContactedAt")
	}
	if !runners[2].Paused {
		t.Error("expected runner 300 to be paused")
	}

	// Export runners
	if err := exp.ExportRunners(ctx, runners, fetchedAt); err != nil {
		t.Fatalf("ExportRunners failed: %v", err)
	}

	// Wait a bit for async export to complete
	time.Sleep(100 * time.Millisecond)

	// Verify the runners were stored in the datastore
	ds := rec.Datastore()
	recordedRunners := ds.ListRunners()

	if len(recordedRunners) != 3 {
		t.Fatalf("expected 3 runners in datastore, got %d", len(recordedRunners))
	}

	// Verify runner 100
	runner100 := ds.GetRunner(100)
	if runner100 == nil {
		t.Fatal("runner 100 not found in datastore")
	}
	if runner100.ShortSha != "abc12345" {
		t.Errorf("runner 100: expected ShortSha 'abc12345', got %s", runner100.ShortSha)
	}
	if runner100.Description != "Integration Test Runner" {
		t.Errorf("runner 100: expected description 'Integration Test Runner', got %s", runner100.Description)
	}
	if len(runner100.TagList) != 3 {
		t.Errorf("runner 100: expected 3 tags, got %d", len(runner100.TagList))
	}
	if !runner100.Flags.Active {
		t.Error("runner 100: expected to be active")
	}
	if runner100.CreatedBy == nil {
		t.Error("runner 100: expected to have creator")
	} else if runner100.CreatedBy.Username != "test-admin" {
		t.Errorf("runner 100: expected creator 'test-admin', got %s", runner100.CreatedBy.Username)
	}

	// Verify runner 200
	runner200 := ds.GetRunner(200)
	if runner200 == nil {
		t.Fatal("runner 200 not found in datastore")
	}
	if !runner200.Flags.Locked {
		t.Error("runner 200: expected to be locked")
	}
	if runner200.CreatedBy.GetId() != 0 {
		t.Error("runner 200: expected no creator, i.e. user id 0")
	}

	// Verify runner 300
	runner300 := ds.GetRunner(300)
	if runner300 == nil {
		t.Fatal("runner 300 not found in datastore")
	}
	if !runner300.Flags.Paused {
		t.Error("runner 300: expected to be paused")
	}
	if !runner300.Timestamps.ContactedAt.AsTime().IsZero() {
		t.Errorf("runner 300: expected ContactedAt to be zero-value of time.Time")
	}
}

func TestRunnersPipeline_EmptyRunners(t *testing.T) {
	mux, glab := setupGitLab(t)
	exp, rec := setupExporter(t)

	// Mock the GraphQL endpoint returning no runners
	mux.HandleFunc("/api/graphql/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"runners": map[string]interface{}{
					"nodes": []map[string]interface{}{},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
						"endCursor":   nil,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	ctx := context.Background()

	fetchedAt := time.Now().UTC()
	runners, err := tasks.FetchRunners(ctx, glab)
	if err != nil {
		t.Fatalf("FetchRunners failed: %v", err)
	}

	if len(runners) != 0 {
		t.Errorf("expected 0 runners, got %d", len(runners))
	}

	// Export should succeed with empty list
	if err := exp.ExportRunners(ctx, runners, fetchedAt); err != nil {
		t.Fatalf("ExportRunners failed with empty list: %v", err)
	}

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Verify no runners in datastore
	ds := rec.Datastore()
	if len(ds.ListRunners()) != 0 {
		t.Errorf("expected 0 runners in datastore, got %d", len(ds.ListRunners()))
	}
}

func TestRunnersPipeline_Pagination(t *testing.T) {
	mux, glab := setupGitLab(t)
	exp, rec := setupExporter(t)

	page := 0
	mux.HandleFunc("/api/graphql/", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Variables map[string]interface{} `json:"variables"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		var response map[string]interface{}

		if page == 0 {
			// First page with 2 runners
			response = map[string]interface{}{
				"data": map[string]interface{}{
					"runners": map[string]interface{}{
						"nodes": []map[string]interface{}{
							{
								"id":          "gid://gitlab/Ci::Runner/1",
								"shortSha":    "page1runner1",
								"description": "Page 1 Runner 1",
								"runnerType":  "INSTANCE_TYPE",
								"tagList":     []string{},
								"status":      "ONLINE",
								"active":      true,
								"locked":      false,
								"paused":      false,
								"accessLevel": "NOT_PROTECTED",
								"runUntagged": false,
								"contactedAt": nil,
								"createdAt":   nil,
								"createdBy":   nil,
							},
							{
								"id":          "gid://gitlab/Ci::Runner/2",
								"shortSha":    "page1runner2",
								"description": "Page 1 Runner 2",
								"runnerType":  "GROUP_TYPE",
								"tagList":     []string{},
								"status":      "OFFLINE",
								"active":      false,
								"locked":      false,
								"paused":      false,
								"accessLevel": "NOT_PROTECTED",
								"runUntagged": false,
								"contactedAt": nil,
								"createdAt":   nil,
								"createdBy":   nil,
							},
						},
						"pageInfo": map[string]interface{}{
							"hasNextPage": true,
							"endCursor":   "cursor_page_1",
						},
					},
				},
			}
			page++
		} else {
			// Second page with 1 runner
			response = map[string]interface{}{
				"data": map[string]interface{}{
					"runners": map[string]interface{}{
						"nodes": []map[string]interface{}{
							{
								"id":          "gid://gitlab/Ci::Runner/3",
								"shortSha":    "page2runner1",
								"description": "Page 2 Runner 1",
								"runnerType":  "PROJECT_TYPE",
								"tagList":     []string{},
								"status":      "STALE",
								"active":      true,
								"locked":      false,
								"paused":      false,
								"accessLevel": "REF_PROTECTED",
								"runUntagged": false,
								"contactedAt": nil,
								"createdAt":   nil,
								"createdBy":   nil,
							},
						},
						"pageInfo": map[string]interface{}{
							"hasNextPage": false,
							"endCursor":   nil,
						},
					},
				},
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	ctx := context.Background()

	fetchedAt := time.Now().UTC()
	runners, err := tasks.FetchRunners(ctx, glab)
	if err != nil {
		t.Fatalf("FetchRunners failed: %v", err)
	}

	// Should fetch all runners across pages
	if len(runners) != 3 {
		t.Fatalf("expected 3 runners from pagination, got %d", len(runners))
	}

	// Export all runners
	if err := exp.ExportRunners(ctx, runners, fetchedAt); err != nil {
		t.Fatalf("ExportRunners failed: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Verify all runners were recorded in datastore
	ds := rec.Datastore()
	recordedRunners := ds.ListRunners()

	if len(recordedRunners) != 3 {
		t.Errorf("expected 3 runners in datastore, got %d", len(recordedRunners))
	}
}
