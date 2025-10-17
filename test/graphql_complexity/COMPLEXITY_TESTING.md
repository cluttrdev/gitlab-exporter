# GraphQL Query Complexity Testing

This directory contains a test utility for measuring GitLab GraphQL query
complexity.

## Overview

GitLab's GraphQL API can return complexity information by including the
`queryComplexity` field in a query:

```graphql
query {
  queryComplexity {
    score  # The actual complexity score of this query
    limit  # The maximum allowed complexity (usually 250 for authenticated requests)
  }
  # ... rest of your query
}
```

## Smart Fragment Loading

The tool includes **smart fragment detection** that recursively analyzes
queries to include only the fragments they actually use, including transitive
dependencies. This solves the "unused fragment" validation errors.

**How it works:**
1. Parses your query to find fragment references (e.g., `...ProjectReferenceFields`)
2. Recursively discovers fragment dependencies (if Fragment A uses Fragment B, both are included)
3. Builds the final query with only the necessary fragments

To enable auto-loading from `spec/` files, change `if false` to `if true` in the test:

```go
if true { // Enable auto-loading
    // Auto-loads queries from spec/queries/ and fragments from spec/fragments/
```

**Known Limitations:**
- Some queries may fail due to variable name mismatches (e.g., expecting `$projectPath` but test provides `$ids`)
- These failures don't affect the usefulness of the tool - successfully executed queries still show their complexity scores

## Usage

### Prerequisites

Set the following environment variables:

```bash
export GITLAB_URL="https://gitlab.com"           # Your GitLab instance URL
export GITLAB_TOKEN="glpat-xxxxxxxxxxxxx"        # Personal access token with API scope
export GITLAB_PROJECT_ID="gid://gitlab/Project/123456"  # A valid project ID
```

### Running the Tests

Run the manual complexity test:

```bash
go test -v -run TestQueryComplexity_Manual ./internal/gitlab/graphql
```

This will execute several queries with different configurations and report their complexity scores.

### Example Output

```
Query Complexity Measurements
=============================
Endpoint: https://gitlab.com/api/graphql
Project ID: gid://gitlab/Project/50817395

=== RUN   TestQueryComplexity_Manual/getProjectsPipelines_-_core_only
✓ getProjectsPipelines - core only:
  Score: 145 / 250
  Usage: 58.0%

=== RUN   TestQueryComplexity_Manual/getProjectsPipelines_-_relations_only
✓ getProjectsPipelines - relations only:
  Score: 89 / 250
  Usage: 35.6%

=== RUN   TestQueryComplexity_Manual/getProjectsPipelines_-_both
✓ getProjectsPipelines - both:
  Score: 234 / 250
  Usage: 93.6%
  ⚠️  WARNING: Query uses more than 80% of complexity limit!
```

### Testing Custom Queries

You can modify the test cases in `complexity_test.go` to test different queries or variable combinations.

**Important**: Your query must include `queryComplexity { score limit }` and all necessary fragment definitions.

Example with fragments:

```go
{
    name: "my custom query",
    query: `#graphql
query myQuery($ids: [ID!]) {
    queryComplexity { score limit }
    projects(ids: $ids) {
        nodes {
            ...ProjectReferenceFields
            name
        }
    }
}

fragment ProjectReferenceFields on Project {
    id
    fullPath
}`,
    variables: map[string]interface{}{
        "ids": []string{projectId},
    },
},
```

### Using Spec Files

The test includes **smart fragment loading** that automatically:
1. Loads all queries from `spec/queries/`
2. Loads all fragments from `spec/fragments/`
3. **Intelligently includes only the fragments each query needs** (including transitive dependencies)

To enable, change `if false` to `if true` in line 408 of `complexity_test.go`.

**What you get:**
- Automatic complexity measurement for all your spec queries
- Tests for query variations (core only, relations only, both) for queries with conditional fragments
- Warnings when queries approach or exceed the 200-point complexity limit

**Example output:**
```
=== RUN   TestQueryComplexity_Manual/jobs/getProjectsPipelinesJobs_(core_only)
    complexity_test.go:530: ❌ jobs/getProjectsPipelinesJobs (core only): ERROR - Query has complexity of 209, which exceeds max complexity of 200

=== RUN   TestQueryComplexity_Manual/merge_requests/getProjectsMergeRequests_(core_only)
    complexity_test.go:532: ✓ merge_requests/getProjectsMergeRequests (core only):
    complexity_test.go:533:   Score: 173 / 200
    complexity_test.go:538:   Usage: 86.5%
    complexity_test.go:541:   ⚠️  WARNING: Query uses more than 80% of complexity limit!
```

Helper functions are also available if you want to build custom loading logic:
- `extractQueriesFromFile()` - Extract queries from a .graphql file
- `extractFragmentsFromFile()` - Extract fragments from a .graphql file
- `injectQueryComplexity()` - Add complexity measurement to a query
- `collectRequiredFragments()` - Recursively collect fragment dependencies

## Interpreting Results

- **Score**: The calculated complexity of your query
- **Limit**: The maximum complexity allowed (usually 250)
- **Usage %**: Percentage of the limit used

### Thresholds

- **< 70%**: Safe - plenty of headroom
- **70-80%**: Caution - approaching limit
- **> 80%**: Warning - consider splitting the query

## Tips for Reducing Complexity

1. **Split queries**: Use the `$_core` and `$_relations` parameters to fetch data in separate requests
2. **Limit pagination**: Use `first: N` on nested collections to reduce the number of items fetched
3. **Remove unused fields**: Only request fields you actually need
4. **Filter aggressively**: Use time ranges and filters to reduce result set sizes

## Related Documentation

- [GitLab GraphQL API Limits](https://docs.gitlab.com/ee/api/graphql/#limits)
- [Query Complexity](https://docs.gitlab.com/ee/development/api_graphql_styleguide.html#max-complexity)
