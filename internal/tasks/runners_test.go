package tasks

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/internal/types"
)

func TestFetchRunners_Success(t *testing.T) {
	// Create a test server that mocks the GitLab GraphQL API
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// GraphQL client uses GET requests
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		// Return mock runners response
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"runners": map[string]interface{}{
					"nodes": []map[string]interface{}{
						{
							"id":          "gid://gitlab/Ci::Runner/123",
							"shortSha":    "abc12345",
							"description": "Production Runner",
							"runnerType":  "INSTANCE_TYPE",
							"tagList":     []string{"docker", "linux"},
							"status":      "ONLINE",
							"active":      true,
							"locked":      false,
							"paused":      false,
							"accessLevel": "NOT_PROTECTED",
							"runUntagged": true,
							"contactedAt": time.Now().Format(time.RFC3339),
							"createdAt":   time.Now().Format(time.RFC3339),
							"createdBy": map[string]interface{}{
								"id":       "gid://gitlab/User/456",
								"username": "admin",
								"name":     "Admin User",
							},
						},
						{
							"id":          "gid://gitlab/Ci::Runner/456",
							"shortSha":    "def67890",
							"description": "Development Runner",
							"runnerType":  "PROJECT_TYPE",
							"tagList":     []string{"kubernetes"},
							"status":      "OFFLINE",
							"active":      false,
							"locked":      true,
							"paused":      true,
							"accessLevel": "REF_PROTECTED",
							"runUntagged": false,
							"contactedAt": nil,
							"createdAt":   time.Now().Format(time.RFC3339),
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

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer srv.Close()

	// Create GitLab client with test server URL
	client, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL: srv.URL,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Test FetchRunners
	ctx := context.Background()
	runners, err := FetchRunners(ctx, client)
	if err != nil {
		t.Fatalf("FetchRunners failed: %v", err)
	}

	// Verify results
	if len(runners) != 2 {
		t.Errorf("expected 2 runners, got %d", len(runners))
	}

	// Check first runner
	if runners[0].Id != 123 {
		t.Errorf("expected runner ID 123, got %d", runners[0].Id)
	}
	if runners[0].ShortSha != "abc12345" {
		t.Errorf("expected ShortSha 'abc12345', got %s", runners[0].ShortSha)
	}
	if runners[0].Description != "Production Runner" {
		t.Errorf("expected Description 'Production Runner', got %s", runners[0].Description)
	}
	if runners[0].RunnerType != types.RunnerTypeInstance {
		t.Errorf("expected RunnerType INSTANCE, got %s", runners[0].RunnerType)
	}
	if runners[0].Status != types.RunnerStatusOnline {
		t.Errorf("expected Status ONLINE, got %s", runners[0].Status)
	}
	if !runners[0].Active {
		t.Error("expected Active to be true")
	}
	if len(runners[0].TagList) != 2 {
		t.Errorf("expected 2 tags, got %d", len(runners[0].TagList))
	}

	// Check second runner
	if runners[1].Id != 456 {
		t.Errorf("expected runner ID 456, got %d", runners[1].Id)
	}
	if runners[1].RunnerType != types.RunnerTypeProject {
		t.Errorf("expected RunnerType PROJECT, got %s", runners[1].RunnerType)
	}
	if runners[1].Status != types.RunnerStatusOffline {
		t.Errorf("expected Status OFFLINE, got %s", runners[1].Status)
	}
	if runners[1].Active {
		t.Error("expected Active to be false")
	}
	if !runners[1].Locked {
		t.Error("expected Locked to be true")
	}
	if !runners[1].Paused {
		t.Error("expected Paused to be true")
	}
}

func TestFetchRunners_EmptyList(t *testing.T) {
	// Create a test server that returns no runners
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	}))
	defer srv.Close()

	client, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL: srv.URL,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	runners, err := FetchRunners(ctx, client)
	if err != nil {
		t.Fatalf("FetchRunners failed: %v", err)
	}

	if len(runners) != 0 {
		t.Errorf("expected 0 runners, got %d", len(runners))
	}
}

func TestFetchRunners_WithPagination(t *testing.T) {
	page := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Variables map[string]interface{} `json:"variables"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		var response map[string]interface{}

		if page == 0 {
			// First page
			response = map[string]interface{}{
				"data": map[string]interface{}{
					"runners": map[string]interface{}{
						"nodes": []map[string]interface{}{
							{
								"id":          "gid://gitlab/Ci::Runner/1",
								"shortSha":    "page1",
								"description": "Runner 1",
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
						},
						"pageInfo": map[string]interface{}{
							"hasNextPage": true,
							"endCursor":   "cursor123",
						},
					},
				},
			}
			page++
		} else {
			// Second page
			response = map[string]interface{}{
				"data": map[string]interface{}{
					"runners": map[string]interface{}{
						"nodes": []map[string]interface{}{
							{
								"id":          "gid://gitlab/Ci::Runner/2",
								"shortSha":    "page2",
								"description": "Runner 2",
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
							"hasNextPage": false,
							"endCursor":   nil,
						},
					},
				},
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer srv.Close()

	client, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL: srv.URL,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	runners, err := FetchRunners(ctx, client)
	if err != nil {
		t.Fatalf("FetchRunners failed: %v", err)
	}

	// Should have fetched both pages
	if len(runners) != 2 {
		t.Errorf("expected 2 runners from pagination, got %d", len(runners))
	}

	if runners[0].ShortSha != "page1" {
		t.Errorf("expected first runner from page 1, got %s", runners[0].ShortSha)
	}
	if runners[1].ShortSha != "page2" {
		t.Errorf("expected second runner from page 2, got %s", runners[1].ShortSha)
	}
}

func TestFetchRunners_ContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)

		response := map[string]interface{}{
			"data": map[string]interface{}{
				"runners": map[string]interface{}{
					"nodes": []map[string]interface{}{},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer srv.Close()

	client, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL: srv.URL,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	// Create a context that will be cancelled immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = FetchRunners(ctx, client)
	if err == nil {
		t.Error("expected error from cancelled context")
	}
	// Check if the error is or wraps context.Canceled
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled error, got: %v", err)
	}
}

func TestFetchRunners_InvalidRunnerId(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"data": map[string]interface{}{
				"runners": map[string]interface{}{
					"nodes": []map[string]interface{}{
						{
							"id":          "invalid-id-format",
							"shortSha":    "abc123",
							"description": "Test",
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
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
						"endCursor":   nil,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer srv.Close()

	client, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL: srv.URL,
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	ctx := context.Background()
	runners, err := FetchRunners(ctx, client)

	// Should not fail, but should log error and skip invalid runner
	if err != nil {
		t.Errorf("FetchRunners should not fail on invalid runner: %v", err)
	}

	// Invalid runner should be skipped
	if len(runners) != 0 {
		t.Errorf("expected 0 runners (invalid one skipped), got %d", len(runners))
	}
}
