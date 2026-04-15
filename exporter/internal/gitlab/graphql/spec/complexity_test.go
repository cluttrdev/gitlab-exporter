package spec_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// ComplexityResult contains the complexity metrics returned by GitLab
type ComplexityResult struct {
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

const (
	defaultProjectId       = "gid://gitlab/Project/50817395"
	defaultProjectPath     = "gitlab-exporter/gitlab-exporter"
	defaultPipelineIid     = "1"
	defaultMergeRequestIid = "1"
)

// TestQueryComplexity is a test for measuring query complexity
// Skipped by default, run with: go test -v -run TestQueryComplexity
// Required environment variables:
//
//	GITLAB_URL: GitLab instance URL (e.g., https://gitlab.com)
//	GITLAB_TOKEN: Personal access token with API access
//	GITLAB_PROJECT_ID: A project ID to test with (e.g., gid://gitlab/Project/50817395)
//	GITLAB_PROJECT_PATH: Full path to a project to test with (e.g., gitlab-exporter/gitlab-exporter)
func TestQueryComplexity(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping manual test in short mode")
	}

	endpoint := os.Getenv("GITLAB_URL")
	if endpoint == "" {
		endpoint = "https://gitlab.com"
	}
	token := os.Getenv("GITLAB_TOKEN")
	if token == "" {
		t.Log("GITLAB_TOKEN not set, keep in mind that unauthenticated requests may have lower complexity limits.")
	}
	projectId := os.Getenv("GITLAB_PROJECT_ID")
	if projectId == "" {
		projectId = defaultProjectId
	}
	projectPath := os.Getenv("GITLAB_PROJECT_PATH")
	if projectPath == "" {
		projectPath = defaultProjectPath
	}
	pipelineIid := os.Getenv("GITLAB_PIPELINE_IID")
	if pipelineIid == "" {
		pipelineIid = defaultPipelineIid
	}
	mergeRequestIid := os.Getenv("GITLAB_MERGE_REQUEST_IID")
	if mergeRequestIid == "" {
		mergeRequestIid = defaultMergeRequestIid
	}

	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}
	endpoint += "api/graphql"

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Common variables used by test queries
	commonVars := map[string]interface{}{
		"projectPath":     projectPath,
		"ids":             []string{projectId},
		"pipelineIid":     pipelineIid,
		"mergeRequestIid": mergeRequestIid,
		"updatedAfter":    time.Now().Add(-7 * 24 * time.Hour).Format(time.RFC3339),
		"updatedBefore":   time.Now().Format(time.RFC3339),
	}

	// Define test cases with manually constructed queries
	// These are simplified versions that include only necessary fragments
	type testCase struct {
		name      string
		query     string
		variables map[string]interface{}
	}

	var tests []testCase

	// Load all fragment definitions from spec/fragments
	fragmentsDir := filepath.Join("fragments")
	fragmentFiles, err := filepath.Glob(filepath.Join(fragmentsDir, "*.graphql"))
	if err != nil {
		t.Fatalf("get fragment files: %v", err)
	}

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
	queriesDir := filepath.Join("queries")
	queryFiles, err := filepath.Glob(filepath.Join(queriesDir, "*.graphql"))
	if err != nil {
		t.Fatalf("get query files: %v", err)
	}

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

	t.Log("Query Complexity Measurements")
	t.Log("=============================")
	t.Logf("Endpoint: %s", endpoint)
	commonVarsDump, _ := json.MarshalIndent(commonVars, "", "\t")
	t.Logf("Variables:\n%v", string(commonVarsDump))
	t.Log("")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Measure complexity
			result := measureQueryComplexity(ctx, endpoint, token, tt.query, tt.variables)

			// Report results
			if result.Error != nil {
				switch {
				case strings.Contains(result.Error.Error(), "exceeds max complexity"):
					t.Logf("❌: %v", result.Error)
					return
				case strings.Contains(result.Error.Error(), "provided invalid value"):
					t.Logf("❌: %v", result.Error)
					return
				default:
					// t.Fatalf("❌: ERROR - %v", result.Error)
					t.Logf("❌: ERROR - %v", result.Error)
					return
				}
			}

			var percentage float64
			if result.Limit > 0 {
				percentage = float64(result.Score) / float64(result.Limit) * 100
			}

			resultIcon := "✓"
			if percentage > 80 {
				resultIcon = "!"
			}
			var buf bytes.Buffer
			fmt.Fprintf(&buf, "%s : Score %d", resultIcon, result.Score)
			if percentage > 0 {
				fmt.Fprintf(&buf, " / %d (%.1f%%)", result.Limit, percentage)
			}
			t.Log(buf.String())
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

// ================================ Helpers ================================

// measureQueryComplexity executes a GraphQL query against GitLab and extracts complexity from response
// The query must include: queryComplexity { score limit }
func measureQueryComplexity(ctx context.Context, endpoint string, token string, query string, variables map[string]interface{}) ComplexityResult {
	result := ComplexityResult{
		Variables: variables,
	}

	// Prepare GraphQL request body
	reqBody := map[string]interface{}{
		"query":     query,
		"variables": variables,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		result.Error = fmt.Errorf("marshal request: %w", err)
		return result
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		result.Error = fmt.Errorf("create request: %w", err)
		return result
	}

	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	// Execute request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		result.Error = fmt.Errorf("execute request: %w", err)
		return result
	}
	defer resp.Body.Close()

	// Read and parse response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = fmt.Errorf("read response: %w", err)
		return result
	}

	if resp.StatusCode >= 400 {
		result.Error = fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		return result
	}

	// Parse GraphQL response
	var gqlResp graphqlResponse
	if err := json.Unmarshal(body, &gqlResp); err != nil {
		result.Error = fmt.Errorf("unmarshal response: %w", err)
		return result
	}

	// Check for GraphQL errors
	if len(gqlResp.Errors) > 0 {
		result.Error = fmt.Errorf("GraphQL error: %s", gqlResp.Errors[0].Message)
		return result
	}

	// Extract complexity from response
	if gqlResp.Data.QueryComplexity != nil {
		result.Score = gqlResp.Data.QueryComplexity.Score
		result.Limit = gqlResp.Data.QueryComplexity.Limit
	} else {
		result.Error = fmt.Errorf("queryComplexity not found in response - ensure query includes 'queryComplexity { score limit }'")
	}

	return result
}

// extractQueriesFromFile reads a .graphql file and extracts all query definitions
// Returns a map of query name -> query string
func extractQueriesFromFile(filePath string) (map[string]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	queries := make(map[string]string)

	// Split by "query " to find individual queries
	// This is a simple approach - a proper GraphQL parser would be better
	parts := strings.Split(string(content), "\nquery ")

	for i, part := range parts {
		if i == 0 {
			// First part before any "query" keyword
			if strings.HasPrefix(strings.TrimSpace(part), "query ") {
				part = strings.TrimPrefix(strings.TrimSpace(part), "query ")
			} else {
				continue
			}
		}

		// Extract query name (first word)
		nameMatch := regexp.MustCompile(`^(\w+)`).FindStringSubmatch(part)
		if len(nameMatch) < 2 {
			continue
		}
		queryName := nameMatch[1]

		// Reconstruct the full query
		queryText := "query " + part

		// Find the end of this query (next "query" or end of file)
		// Count braces to find where the query ends
		braceCount := 0
		inQuery := false
		endIdx := 0

		for idx, char := range queryText {
			if char == '{' {
				braceCount++
				inQuery = true
			} else if char == '}' {
				braceCount--
				if inQuery && braceCount == 0 {
					endIdx = idx + 1
					break
				}
			}
		}

		if endIdx > 0 {
			queryText = queryText[:endIdx]
		}

		queries[queryName] = strings.TrimSpace(queryText)
	}

	return queries, nil
}

// extractBooleanParameters finds all boolean parameters in a GraphQL query
// Returns a list of parameter names that are used as conditional flags (start with underscore)
// Example: "$_core: Boolean = true" -> ["_core"]
func extractBooleanParameters(queryText string) []string {
	var params []string

	// Match patterns like: $_paramName: Boolean
	// We look for parameters that start with underscore (our convention for conditional flags)
	re := regexp.MustCompile(`\$(_\w+)\s*:\s*Boolean`)
	matches := re.FindAllStringSubmatch(queryText, -1)

	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			paramName := match[1]
			if !seen[paramName] {
				params = append(params, paramName)
				seen[paramName] = true
			}
		}
	}

	return params
}

// injectQueryComplexity injects "queryComplexity { score limit }" into a GraphQL query
// It adds it as the first field in the query selection set
func injectQueryComplexity(query string) string {
	// Pattern to match: query NAME(...params) {
	// We want to inject right after the opening brace
	re := regexp.MustCompile(`(?m)(query\s+\w+[^{]*\{\s*)`)

	return re.ReplaceAllString(query, "${1}\n\tqueryComplexity { score limit }\n\t")
}

// extractFragmentsFromFile reads a .graphql file and extracts all fragment definitions
// Returns a map of fragment name -> fragment definition
func extractFragmentsFromFile(filePath string) (map[string]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	fragments := make(map[string]string)

	// Split by "fragment " to find individual fragments
	parts := strings.Split(string(content), "\nfragment ")

	for i, part := range parts {
		if i == 0 {
			// First part before any "fragment" keyword
			if strings.HasPrefix(strings.TrimSpace(part), "fragment ") {
				part = strings.TrimPrefix(strings.TrimSpace(part), "fragment ")
			} else {
				continue
			}
		}

		// Extract fragment name (first word)
		nameMatch := regexp.MustCompile(`^(\w+)`).FindStringSubmatch(part)
		if len(nameMatch) < 2 {
			continue
		}
		fragmentName := nameMatch[1]

		// Reconstruct the full fragment
		fragmentText := "fragment " + part

		// Find the end of this fragment (next "fragment" or end of file)
		// Count braces to find where the fragment ends
		braceCount := 0
		inFragment := false
		endIdx := 0

		for idx, char := range fragmentText {
			if char == '{' {
				braceCount++
				inFragment = true
			} else if char == '}' {
				braceCount--
				if inFragment && braceCount == 0 {
					endIdx = idx + 1
					break
				}
			}
		}

		if endIdx > 0 {
			fragmentText = fragmentText[:endIdx]
		}

		fragments[fragmentName] = strings.TrimSpace(fragmentText)
	}

	return fragments, nil
}

// mergeVars merges multiple variable maps, with later maps overriding earlier ones
func mergeVars(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// collectRequiredFragments recursively collects all fragments needed by a query
// including transitive dependencies
func collectRequiredFragments(queryText string, allFragments map[string]string) map[string]string {
	required := make(map[string]string)
	visited := make(map[string]bool)

	var collectDeps func(text string)
	collectDeps = func(text string) {
		// Find all fragment references in this text
		refs := extractFragmentReferences(text)

		for fragName := range refs {
			// Skip if already processed
			if visited[fragName] {
				continue
			}
			visited[fragName] = true

			// Get the fragment definition
			fragDef, exists := allFragments[fragName]
			if !exists {
				// Fragment not found, skip (it might be inline or missing)
				continue
			}

			// Add this fragment to required list
			required[fragName] = fragDef

			// Recursively collect dependencies of this fragment
			collectDeps(fragDef)
		}
	}

	// Start with the query text
	collectDeps(queryText)

	return required
}

// extractFragmentReferences finds all fragment references in a GraphQL query/fragment
// Returns a set of fragment names (e.g., "...ProjectReferenceFields" -> "ProjectReferenceFields")
func extractFragmentReferences(graphqlText string) map[string]bool {
	references := make(map[string]bool)

	// Match fragment spreads: ...FragmentName
	re := regexp.MustCompile(`\.\.\.(\w+)`)
	matches := re.FindAllStringSubmatch(graphqlText, -1)

	for _, match := range matches {
		if len(match) > 1 {
			references[match[1]] = true
		}
	}

	return references
}

// generateFlagCombinations generates all boolean combinations for the given flag parameters
// Returns a map of variation name -> variable map
// Example: ["_core", "_relations"] ->
//
//	{"core only": {_core: true, _relations: false},
//	 "relations only": {_core: false, _relations: true},
//	 "both": {_core: true, _relations: true}}
func generateFlagCombinations(flags []string) map[string]map[string]interface{} {
	if len(flags) == 0 {
		return nil
	}

	combinations := make(map[string]map[string]interface{})

	// Generate all 2^n combinations
	numCombinations := 1 << len(flags) // 2^n

	for i := 0; i < numCombinations; i++ {
		vars := make(map[string]interface{})
		var nameParts []string

		// Build the combination based on bit pattern
		allFalse := true
		for j, flag := range flags {
			bit := (i >> j) & 1
			value := bit == 1
			vars[flag] = value

			if value {
				allFalse = false
				// Create readable name part (remove leading underscore)
				nameParts = append(nameParts, strings.TrimPrefix(flag, "_"))
			}
		}

		// Skip the "all false" combination as it's usually not meaningful
		if allFalse {
			continue
		}

		// Create a readable name for this combination
		var name string
		if len(nameParts) == len(flags) {
			name = "all"
		} else if len(nameParts) == 1 {
			name = nameParts[0] + " only"
		} else {
			name = strings.Join(nameParts, "+")
		}

		combinations[name] = vars
	}

	return combinations
}
