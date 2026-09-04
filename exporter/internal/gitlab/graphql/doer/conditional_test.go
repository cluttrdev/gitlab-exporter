package doer

import (
	"strings"
	"testing"
)

// conditionalTest exercises stripNodes' @include / @skip materialization.
// Like rewriteTest, assertions are substring checks: gqlparser's formatter is
// stable but its whitespace is an implementation detail.
type conditionalTest struct {
	name            string
	server          string
	vars            map[string]any
	input           string
	wantContains    []string
	wantNotContains []string
	wantPruned      []string
}

// A query in the shape the generated client actually emits: reference fields
// unconditionally, plus two mutually optional field-set fragments.
const splitQuery = `query Q($id: ID!, $_core: Boolean!, $_extra: Boolean!) {
  project(id: $id) {
    id
    ...Core @include(if: $_core)
    ...Extra @include(if: $_extra)
  }
}
fragment Core on Project {
  name
}
fragment Extra on Project {
  description
}`

func TestStripNodes_Conditionals(t *testing.T) {
	tests := []conditionalTest{
		{
			// The case this whole mechanism exists for: below 17.0 the server
			// charges for Extra even though it never resolves it, so the
			// spread, its now-orphaned definition, and both flag variables go.
			name:            "excluded spread and its orphaned fragment are removed",
			server:          "16.0",
			vars:            map[string]any{"_core": true, "_extra": false},
			input:           splitQuery,
			wantContains:    []string{"... Core", "fragment Core", "$id"},
			wantNotContains: []string{"Extra", "@include", "$_core", "$_extra"},
			wantPruned:      []string{"_core", "_extra"},
		},
		{
			name:            "included spread survives with the directive consumed",
			server:          "16.0",
			vars:            map[string]any{"_core": true, "_extra": true},
			input:           splitQuery,
			wantContains:    []string{"... Core", "... Extra", "fragment Core", "fragment Extra"},
			wantNotContains: []string{"@include", "$_core", "$_extra"},
			wantPruned:      []string{"_core", "_extra"},
		},
		{
			// 17.0+ computes complexity correctly on its own, so the rewriter
			// keeps its hands off and the server sees the query verbatim.
			name:            "modern server is left untouched",
			server:          "17.0",
			vars:            map[string]any{"_core": true, "_extra": false},
			input:           splitQuery,
			wantContains:    []string{"... Core @include(if: $_core)", "... Extra @include(if: $_extra)", "$_extra"},
			wantNotContains: []string{},
		},
		{
			name:   "field-level include is materialized too",
			server: "16.0",
			vars:   map[string]any{"_commits": false},
			input: `query Q($_commits: Boolean!) {
  mergeRequest {
    iid
    commits @include(if: $_commits) {
      nodes {
        sha
      }
    }
  }
}`,
			wantContains:    []string{"iid"},
			wantNotContains: []string{"commits", "@include", "$_commits"},
			wantPruned:      []string{"_commits"},
		},
		{
			name:   "inline fragment is materialized",
			server: "16.0",
			vars:   map[string]any{"f": false},
			input: `query Q($f: Boolean!) {
  project {
    id
    ... on Project @include(if: $f) {
      name
    }
  }
}`,
			wantContains:    []string{"id"},
			wantNotContains: []string{"name", "@include", "$f"},
			wantPruned:      []string{"f"},
		},
		{
			name:   "skip inverts include",
			server: "16.0",
			vars:   map[string]any{"f": true},
			input: `query Q($f: Boolean!) {
  project {
    id
    name @skip(if: $f)
    description @skip(if: $f)
  }
}`,
			wantContains:    []string{"id"},
			wantNotContains: []string{"name", "description", "@skip", "$f"},
		},
		{
			// Per spec, a node is rendered only when @skip is false AND
			// @include is true; either one excluding wins.
			name:   "skip and include on one node, either excluding wins",
			server: "16.0",
			vars:   map[string]any{"s": true, "i": true},
			input: `query Q($s: Boolean!, $i: Boolean!) {
  project {
    id
    name @skip(if: $s) @include(if: $i)
  }
}`,
			wantContains:    []string{"id"},
			wantNotContains: []string{"name"},
		},
		{
			name:   "skip and include on one node, both permitting keeps it",
			server: "16.0",
			vars:   map[string]any{"s": false, "i": true},
			input: `query Q($s: Boolean!, $i: Boolean!) {
  project {
    id
    name @skip(if: $s) @include(if: $i)
  }
}`,
			wantContains:    []string{"id", "name"},
			wantNotContains: []string{"@skip", "@include"},
		},
		{
			name:   "literal condition needs no variable",
			server: "16.0",
			vars:   nil,
			input: `query Q {
  project {
    id
    name @include(if: false)
    description @include(if: true)
  }
}`,
			wantContains:    []string{"id", "description"},
			wantNotContains: []string{"name", "@include"},
		},
		{
			name:   "variable default fills in for an absent request variable",
			server: "16.0",
			vars:   map[string]any{},
			input: `query Q($f: Boolean = false) {
  project {
    id
    name @include(if: $f)
  }
}`,
			wantContains:    []string{"id"},
			wantNotContains: []string{"name", "@include"},
		},
		{
			// Fail open. Pruning a node the server would have resolved
			// silently loses data, so an unresolvable condition is handed to
			// the server exactly as authored.
			name:   "unresolvable variable leaves node and directive alone",
			server: "16.0",
			vars:   map[string]any{},
			input: `query Q($f: Boolean!) {
  project {
    id
    name @include(if: $f)
  }
}`,
			wantContains: []string{"id", "name", "@include(if: $f)", "$f"},
		},
		{
			name:   "non-boolean variable value is unresolvable",
			server: "16.0",
			vars:   map[string]any{"f": "true"},
			input: `query Q($f: Boolean!) {
  project {
    id
    name @include(if: $f)
  }
}`,
			wantContains: []string{"id", "name", "@include(if: $f)"},
		},
		{
			// Fragment definitions are shared by every operation in a
			// document, so there is no single value to prune their bodies
			// against. Bail out rather than guess.
			name:   "multi-operation document is not materialized",
			server: "16.0",
			vars:   map[string]any{"f": false},
			input: `query A($f: Boolean!) {
  project {
    name @include(if: $f)
  }
}
query B($f: Boolean!) {
  group {
    name @include(if: $f)
  }
}`,
			wantContains: []string{"@include(if: $f)"},
		},
		{
			name:   "emptied selection set gets a __typename",
			server: "16.0",
			vars:   map[string]any{"f": false},
			input: `query Q($f: Boolean!) {
  project {
    description @include(if: $f)
  }
}`,
			wantContains:    []string{"__typename"},
			wantNotContains: []string{"description", "@include"},
		},
		{
			name:   "minVersion drops a node its include would have kept",
			server: "16.0",
			vars:   map[string]any{"f": true},
			input: `query Q($f: Boolean!) {
  project {
    id
    updatedAt @minVersion(v: "17.0") @include(if: $f)
  }
}`,
			wantContains:    []string{"id"},
			wantNotContains: []string{"updatedAt", "@minVersion", "@include"},
		},
		{
			name:   "include drops a node its minVersion would have kept",
			server: "16.5",
			vars:   map[string]any{"f": false},
			input: `query Q($f: Boolean!) {
  project {
    id
    name @minVersion(v: "16.0") @include(if: $f)
  }
}`,
			wantContains:    []string{"id"},
			wantNotContains: []string{"name", "@minVersion", "@include"},
		},
		{
			name:   "both permitting keeps the node and consumes both directives",
			server: "16.5",
			vars:   map[string]any{"f": true},
			input: `query Q($f: Boolean!) {
  project {
    id
    name @minVersion(v: "16.0") @include(if: $f)
  }
}`,
			wantContains:    []string{"id", "name"},
			wantNotContains: []string{"@minVersion", "@include"},
		},
		{
			// A node dropped by @minVersion takes its unconsumed directives
			// with it, including an @include the rewriter could not resolve.
			// Nothing harvests that `if:` variable except the drop path
			// itself, and leaving it declared makes the server reject the
			// request with "Variable $f is declared by Q but not used".
			name:   "spread dropped by version harvests an unresolved include variable",
			server: "16.0",
			vars:   map[string]any{},
			input: `query Q($f: Boolean!) {
  project {
    id
    ...Legacy @minVersion(v: "17.0") @include(if: $f)
  }
}
fragment Legacy on Project {
  name
}`,
			wantContains:    []string{"id"},
			wantNotContains: []string{"Legacy", "$f", "@include"},
			wantPruned:      []string{"f"},
		},
		{
			name:   "inline fragment dropped by version harvests an unresolved include variable",
			server: "16.0",
			vars:   map[string]any{},
			input: `query Q($f: Boolean!) {
  project {
    id
    ... on Project @minVersion(v: "17.0") @include(if: $f) {
      name
    }
  }
}`,
			wantContains:    []string{"id"},
			wantNotContains: []string{"$f", "@include"},
			wantPruned:      []string{"f"},
		},
		{
			// Same drop path on a server that materializes nothing: @include
			// is never consumed there, so the version drop is the only thing
			// that can harvest the variable.
			name:   "modern server still harvests the include variable of a version-dropped spread",
			server: "17.0",
			vars:   nil,
			input: `query Q($f: Boolean!) {
  project {
    id
    ...Legacy @maxVersion(v: "16.11") @include(if: $f)
  }
}
fragment Legacy on Project {
  name
}`,
			wantContains:    []string{"id"},
			wantNotContains: []string{"Legacy", "$f"},
			wantPruned:      []string{"f"},
		},
		{
			// A flag variable still referenced by a surviving node must not be
			// pruned, or the server rejects the request for an undefined
			// variable.
			name:   "shared flag variable survives while any reference remains",
			server: "16.0",
			vars:   map[string]any{"f": true},
			input: `query Q($f: Boolean!, $g: Boolean!) {
  project {
    id
    name @include(if: $f)
    description @include(if: $g)
  }
}`,
			wantContains:    []string{"id", "name", "description", "@include(if: $g)", "$g"},
			wantNotContains: []string{"$f", "@include(if: $f)"},
			wantPruned:      []string{"f"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGate(tt.server)
			if err != nil {
				t.Fatalf("NewGate(%q): %v", tt.server, err)
			}
			out, pruned, err := stripNodes(tt.input, rewriteOpts{
				gate:        g,
				vars:        tt.vars,
				materialize: g.MaterializesConditionals(),
			})
			if err != nil {
				t.Fatalf("stripNodes: %v", err)
			}
			for _, want := range tt.wantContains {
				if !strings.Contains(out, want) {
					t.Errorf("output missing %q\nfull output:\n%s", want, out)
				}
			}
			for _, banned := range tt.wantNotContains {
				if strings.Contains(out, banned) {
					t.Errorf("output unexpectedly contains %q\nfull output:\n%s", banned, out)
				}
			}
			if tt.wantPruned != nil {
				if len(pruned) != len(tt.wantPruned) {
					t.Errorf("pruned vars: got %v, want %v", keysOf(pruned), tt.wantPruned)
				}
				for _, v := range tt.wantPruned {
					if _, ok := pruned[v]; !ok {
						t.Errorf("expected variable %q to be pruned, pruned=%v", v, keysOf(pruned))
					}
				}
			}
		})
	}
}

// TestGate_MaterializesConditionals pins the version boundary. graphql-ruby
// only honored @include / @skip on fragment spreads and inline fragments from
// 2.3.0, which GitLab picked up in 17.0.
func TestGate_MaterializesConditionals(t *testing.T) {
	tests := []struct {
		server string
		want   bool
	}{
		{"15.11.13-ee", true},
		{"16.0.0-ee", true},
		{"16.11.9", true},
		{"17.0.0-ee.0", false},
		{"17.7.0", false},
		{"18.1.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.server, func(t *testing.T) {
			g, err := NewGate(tt.server)
			if err != nil {
				t.Fatalf("NewGate(%q): %v", tt.server, err)
			}
			if got := g.MaterializesConditionals(); got != tt.want {
				t.Errorf("MaterializesConditionals() = %v, want %v", got, tt.want)
			}
		})
	}

	if (*Gate)(nil).MaterializesConditionals() {
		t.Error("a nil gate must not claim to need materialization")
	}
}

// TestStripForVersion_IgnoresConditionals asserts the older entry point still
// leaves @include alone, so callers that pass no variables are unaffected.
func TestStripForVersion_IgnoresConditionals(t *testing.T) {
	g, err := NewGate("16.0")
	if err != nil {
		t.Fatal(err)
	}
	out, err := StripForVersion(splitQuery, g)
	if err != nil {
		t.Fatalf("StripForVersion: %v", err)
	}
	if !strings.Contains(out, "@include(if: $_extra)") {
		t.Errorf("StripForVersion must not materialize @include\noutput:\n%s", out)
	}
}

// TestStripForVersionWithVariables_MatchesStripForVersionOnModernServers
// guards the promise that 17.0+ behavior is byte-for-byte unchanged.
func TestStripForVersionWithVariables_MatchesStripForVersionOnModernServers(t *testing.T) {
	for _, server := range []string{"17.0", "17.7.0-ee", "18.1.0"} {
		t.Run(server, func(t *testing.T) {
			g, err := NewGate(server)
			if err != nil {
				t.Fatal(err)
			}
			plain, err := StripForVersion(splitQuery, g)
			if err != nil {
				t.Fatal(err)
			}
			withVars, err := StripForVersionWithVariables(splitQuery, g,
				map[string]any{"_core": true, "_extra": false})
			if err != nil {
				t.Fatal(err)
			}
			if plain != withVars {
				t.Errorf("output differs on a server that needs no materialization\nStripForVersion:\n%s\nStripForVersionWithVariables:\n%s",
					plain, withVars)
			}
		})
	}
}
