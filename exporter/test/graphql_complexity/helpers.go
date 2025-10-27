package graphql_complexity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

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
