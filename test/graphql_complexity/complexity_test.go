package graphql_complexity

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// ComplexityResult contains the complexity metrics returned by GitLab
type ComplexityResult struct {
	QueryName string
	Score     int
	Limit     int
	Variables map[string]interface{}
	Error     error
}

// graphqlResponse represents the structure of a GraphQL response
type graphqlResponse struct {
	Data struct {
		QueryComplexity *struct {
			Score int `json:"score"`
			Limit int `json:"limit"`
		} `json:"queryComplexity"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// TestQueryComplexity_Manual is a manual test for measuring query complexity
// Skip by default, run with: go test -v -run TestQueryComplexity_Manual
// Required environment variables:
//
//	GITLAB_URL: GitLab instance URL (e.g., https://gitlab.com)
//	GITLAB_TOKEN: Personal access token with API access
//	GITLAB_PROJECT_ID: A project ID to test with (e.g., gid://gitlab/Project/50817395)
func TestQueryComplexity_Manual(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping manual test in short mode")
	}

	endpoint := os.Getenv("GITLAB_URL")
	if endpoint == "" {
		endpoint = "https://gitlab.com"
	}
	projectId := os.Getenv("GITLAB_PROJECT_ID")
	if projectId == "" {
		projectId = "gid://gitlab/Project/50817395" // Default to gitlab-org/gitlab project
	}
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		t.Log("GITLAB_TOKEN not set, keep in mind that unauthenticated requests may have lower complexity limits.")
	}

	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}
	endpoint += "api/graphql"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Common variables used by test queries
	commonVars := map[string]interface{}{
		"ids":           []string{projectId},
		"updatedAfter":  time.Now().Add(-7 * 24 * time.Hour).Format(time.RFC3339),
		"updatedBefore": time.Now().Format(time.RFC3339),
	}

	// Define test cases with manually constructed queries
	// These are simplified versions that include only necessary fragments
	type testCase struct {
		name      string
		query     string
		variables map[string]interface{}
	}
	tests := []testCase{
		{
			name: "simple project query",
			query: `#graphql
query getProjects($ids: [ID!]) {
	queryComplexity { score limit }
	projects(ids: $ids) {
		nodes {
			...ProjectReferenceFields
			name
			description
		}
	}
}

fragment ProjectReferenceFields on Project {
	id
	fullPath
}`,
			variables: commonVars,
		},
		{
			name: "projects with pipeline counts",
			query: `#graphql
query getProjects($ids: [ID!], $updatedAfter: Time, $updatedBefore: Time) {
	queryComplexity { score limit }
	projects(ids: $ids) {
		nodes {
			...ProjectReferenceFields
			pipelines(scope: FINISHED, updatedAfter: $updatedAfter, updatedBefore: $updatedBefore, first: 10) {
				count
			}
		}
	}
}

fragment ProjectReferenceFields on Project {
	id
	fullPath
}`,
			variables: commonVars,
		},
	}

	// Optional: Uncomment this section to auto-load queries from spec/ files
	// Now with smart fragment detection - only includes fragments that are actually used!
	if false { // Change to 'true' to enable auto-loading
		// Load all fragment definitions from spec/fragments
		fragmentsDir := filepath.Join("spec", "fragments")
		fragmentFiles, err := filepath.Glob(filepath.Join(fragmentsDir, "*.graphql"))
		if err == nil {
			allFragments := make(map[string]string)
			for _, fragFile := range fragmentFiles {
				frags, err := extractFragmentsFromFile(fragFile)
				if err != nil {
					t.Logf("Warning: Failed to load fragments from %s: %v", fragFile, err)
					continue
				}
				// Merge fragments into master map
				for name, def := range frags {
					allFragments[name] = def
				}
			}

			// Load queries from spec/queries
			queriesDir := filepath.Join("spec", "queries")
			queryFiles, err := filepath.Glob(filepath.Join(queriesDir, "*.graphql"))
			if err == nil {
				for _, queryFile := range queryFiles {
					queries, err := extractQueriesFromFile(queryFile)
					if err != nil {
						t.Logf("Warning: Failed to load queries from %s: %v", queryFile, err)
						continue
					}

					fileName := filepath.Base(queryFile)
					fileNameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))

					for queryName, queryText := range queries {
						// Inject queryComplexity field
						queryWithComplexity := injectQueryComplexity(queryText)

						// Smart fragment collection - only include what's needed!
						requiredFragments := collectRequiredFragments(queryWithComplexity, allFragments)

						// Build fragment definitions string
						var fragmentDefs strings.Builder
						for _, fragDef := range requiredFragments {
							fragmentDefs.WriteString("\n")
							fragmentDefs.WriteString(fragDef)
						}

						fullQuery := queryWithComplexity + fragmentDefs.String()

						// Dynamically detect boolean flag parameters
						boolFlags := extractBooleanParameters(queryText)

						if len(boolFlags) > 0 {
							// Generate all combinations for boolean flags
							combinations := generateFlagCombinations(boolFlags)

							for variantName, flagVars := range combinations {
								tests = append(tests, testCase{
									name:      fmt.Sprintf("%s/%s (%s)", fileNameWithoutExt, queryName, variantName),
									query:     fullQuery,
									variables: mergeVars(commonVars, flagVars),
								})
							}
						} else {
							// No boolean flags, just add the query as-is
							tests = append(tests, testCase{
								name:      fmt.Sprintf("%s/%s", fileNameWithoutExt, queryName),
								query:     fullQuery,
								variables: commonVars,
							})
						}
					}
				}
			}
		}
	}

	t.Log("Query Complexity Measurements")
	t.Log("=============================")
	t.Logf("Endpoint: %s", endpoint)
	t.Logf("Project ID: %s", projectId)
	t.Log("")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Measure complexity
			result := measureQueryComplexity(ctx, endpoint, token, tt.query, tt.variables)
			result.QueryName = tt.name

			// Report results
			if result.Error != nil {
				t.Logf("❌ %s: ERROR - %v", result.QueryName, result.Error)
			} else {
				t.Logf("✓ %s:", result.QueryName)
				t.Logf("  Score: %d / %d", result.Score, result.Limit)

				// Check if near limit
				if result.Limit > 0 {
					percentage := float64(result.Score) / float64(result.Limit) * 100
					t.Logf("  Usage: %.1f%%", percentage)

					if percentage > 80 {
						t.Logf("  ⚠️  WARNING: Query uses more than 80%% of complexity limit!")
					}
				}
			}
			t.Log("")
		})
	}
}

// TestInjectQueryComplexity tests the queryComplexity injection
func TestInjectQueryComplexity(t *testing.T) {
	input := `query getProjects($ids: [ID!]) {
	projects(ids: $ids) {
		nodes {
			id
		}
	}
}`

	result := injectQueryComplexity(input)

	if !strings.Contains(result, "queryComplexity { score limit }") {
		t.Errorf("Expected queryComplexity to be injected, got: %s", result)
	}

	// Verify it's injected after the opening brace
	queryIdx := strings.Index(result, "query getProjects")
	complexityIdx := strings.Index(result, "queryComplexity")
	projectsIdx := strings.Index(result, "projects(")

	if complexityIdx < queryIdx || complexityIdx > projectsIdx {
		t.Errorf("queryComplexity not injected in the correct position")
	}
}

// TestExtractFragmentReferences tests fragment reference extraction
func TestExtractFragmentReferences(t *testing.T) {
	query := `
query getProjects {
	projects {
		nodes {
			...ProjectReferenceFields
			name
			pipelines {
				...PipelineReferenceFields
			}
		}
	}
}
`
	refs := extractFragmentReferences(query)

	if !refs["ProjectReferenceFields"] {
		t.Error("Expected ProjectReferenceFields to be found")
	}
	if !refs["PipelineReferenceFields"] {
		t.Error("Expected PipelineReferenceFields to be found")
	}
	if len(refs) != 2 {
		t.Errorf("Expected 2 references, got %d", len(refs))
	}
}

// TestExtractFragmentsFromFile tests fragment extraction from a file
func TestExtractFragmentsFromFile(t *testing.T) {
	// Create a temporary test file
	tempFile := filepath.Join(t.TempDir(), "fragments.graphql")
	testContent := `fragment ProjectReferenceFields on Project {
	id
	fullPath
}

fragment PipelineReferenceFields on Pipeline {
	id
	iid
	...ProjectReferenceFields
}
`
	if err := os.WriteFile(tempFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	fragments, err := extractFragmentsFromFile(tempFile)
	if err != nil {
		t.Fatalf("extractFragmentsFromFile failed: %v", err)
	}

	if len(fragments) != 2 {
		t.Errorf("Expected 2 fragments, got %d", len(fragments))
	}

	if _, ok := fragments["ProjectReferenceFields"]; !ok {
		t.Error("Expected ProjectReferenceFields fragment to be extracted")
	}

	if _, ok := fragments["PipelineReferenceFields"]; !ok {
		t.Error("Expected PipelineReferenceFields fragment to be extracted")
	}

	// Verify fragment content
	if frag, ok := fragments["ProjectReferenceFields"]; ok {
		if !strings.Contains(frag, "id") || !strings.Contains(frag, "fullPath") {
			t.Errorf("Fragment content incorrect: %s", frag)
		}
	}
}

// TestCollectRequiredFragments tests recursive fragment dependency collection
func TestCollectRequiredFragments(t *testing.T) {
	// Define a query that uses fragments
	query := `
query getProjects {
	projects {
		nodes {
			...ProjectFields
		}
	}
}
`

	// Define fragments with dependencies
	allFragments := map[string]string{
		"ProjectFields": `fragment ProjectFields on Project {
	...ProjectReferenceFields
	name
	description
}`,
		"ProjectReferenceFields": `fragment ProjectReferenceFields on Project {
	id
	fullPath
}`,
		"UnusedFragment": `fragment UnusedFragment on User {
	id
	username
}`,
	}

	required := collectRequiredFragments(query, allFragments)

	// Should include ProjectFields and ProjectReferenceFields, but NOT UnusedFragment
	if len(required) != 2 {
		t.Errorf("Expected 2 required fragments, got %d: %v", len(required), required)
	}

	if _, ok := required["ProjectFields"]; !ok {
		t.Error("Expected ProjectFields to be required")
	}

	if _, ok := required["ProjectReferenceFields"]; !ok {
		t.Error("Expected ProjectReferenceFields to be required (transitive dependency)")
	}

	if _, ok := required["UnusedFragment"]; ok {
		t.Error("UnusedFragment should not be required")
	}
}

// TestExtractBooleanParameters tests boolean parameter extraction
func TestExtractBooleanParameters(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		expected []string
	}{
		{
			name: "query with core and relations flags",
			query: `query getPipelines($ids: [ID!], $_core: Boolean = true, $_relations: Boolean = false) {
				pipelines {
					...PipelineFieldsCore @include(if: $_core)
					...PipelineFieldsRelations @include(if: $_relations)
				}
			}`,
			expected: []string{"_core", "_relations"},
		},
		{
			name: "query with core and extra flags",
			query: `query getJobs($ids: [ID!], $_core: Boolean = true, $_extra: Boolean = false) {
				jobs {
					...JobFieldsCore @include(if: $_core)
					...JobFieldsExtra @include(if: $_extra)
				}
			}`,
			expected: []string{"_core", "_extra"},
		},
		{
			name:     "query without boolean flags",
			query:    `query getProjects($ids: [ID!]) { projects(ids: $ids) { id } }`,
			expected: []string{},
		},
		{
			name: "query with non-underscore boolean (should be ignored)",
			query: `query test($includeArchived: Boolean) {
				projects(includeArchived: $includeArchived) { id }
			}`,
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractBooleanParameters(tt.query)

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d parameters, got %d: %v", len(tt.expected), len(result), result)
				return
			}

			// Check that all expected params are present
			resultMap := make(map[string]bool)
			for _, p := range result {
				resultMap[p] = true
			}

			for _, expected := range tt.expected {
				if !resultMap[expected] {
					t.Errorf("Expected parameter %s not found in result: %v", expected, result)
				}
			}
		})
	}
}

// TestGenerateFlagCombinations tests flag combination generation
func TestGenerateFlagCombinations(t *testing.T) {
	tests := []struct {
		name          string
		flags         []string
		expectedCount int
		checkNames    []string
	}{
		{
			name:          "two flags (core, relations)",
			flags:         []string{"_core", "_relations"},
			expectedCount: 3, // core only, relations only, both (all false skipped)
			checkNames:    []string{"core only", "relations only", "all"},
		},
		{
			name:          "two flags (core, extra)",
			flags:         []string{"_core", "_extra"},
			expectedCount: 3,
			checkNames:    []string{"core only", "extra only", "all"},
		},
		{
			name:          "three flags",
			flags:         []string{"_core", "_relations", "_extra"},
			expectedCount: 7, // 2^3 - 1 (all false skipped)
			checkNames:    []string{"core only", "relations only", "extra only", "all"},
		},
		{
			name:          "no flags",
			flags:         []string{},
			expectedCount: 0,
			checkNames:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateFlagCombinations(tt.flags)

			if len(result) != tt.expectedCount {
				t.Errorf("Expected %d combinations, got %d", tt.expectedCount, len(result))
			}

			// Check that expected names exist
			for _, name := range tt.checkNames {
				if _, ok := result[name]; !ok {
					t.Errorf("Expected combination name '%s' not found. Got: %v", name, result)
				}
			}

			// Verify that each combination has the right flags
			for name, vars := range result {
				if len(vars) != len(tt.flags) {
					t.Errorf("Combination '%s' has %d flags, expected %d", name, len(vars), len(tt.flags))
				}

				// Check that at least one flag is true (all-false should be skipped)
				hasTrue := false
				for _, v := range vars {
					if v == true {
						hasTrue = true
						break
					}
				}
				if !hasTrue {
					t.Errorf("Combination '%s' has all flags false, should be skipped", name)
				}
			}
		})
	}
}

// TestExtractQueriesFromFile tests query extraction from a file
func TestExtractQueriesFromFile(t *testing.T) {
	// Create a temporary test file
	tempFile := filepath.Join(t.TempDir(), "test.graphql")
	testContent := `query getProjects($ids: [ID!]) {
	projects(ids: $ids) {
		nodes {
			id
		}
	}
}

query getProject($id: ID!) {
	project(id: $id) {
		id
		name
	}
}
`
	if err := os.WriteFile(tempFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	queries, err := extractQueriesFromFile(tempFile)
	if err != nil {
		t.Fatalf("extractQueriesFromFile failed: %v", err)
	}

	if len(queries) != 2 {
		t.Errorf("Expected 2 queries, got %d", len(queries))
	}

	if _, ok := queries["getProjects"]; !ok {
		t.Error("Expected getProjects query to be extracted")
	}

	if _, ok := queries["getProject"]; !ok {
		t.Error("Expected getProject query to be extracted")
	}

	// Verify query content
	if query, ok := queries["getProjects"]; ok {
		if !strings.Contains(query, "projects(ids: $ids)") {
			t.Errorf("Query content incorrect: %s", query)
		}
	}
}
