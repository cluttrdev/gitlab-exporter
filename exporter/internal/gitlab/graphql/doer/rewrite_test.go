package doer

import (
	"strings"
	"testing"
)

// rewriteTest exercises stripVersionedNodes against a fixed input query at a
// specified server version. The wantContains / wantNotContains assertions are
// substring checks against the formatted output - we don't pin the exact
// formatter output because gqlparser's formatter is stable but its choice of
// whitespace is an implementation detail.
type rewriteTest struct {
	name            string
	server          string
	input           string
	wantContains    []string
	wantNotContains []string
	wantPruned      []string
	wantErr         bool
}

func TestStripVersionedNodes(t *testing.T) {
	tests := []rewriteTest{
		{
			name:   "passing minVersion is kept",
			server: "17.0",
			input: `query Q {
  project {
    updatedAt @minVersion(v: "17.0")
    name
  }
}`,
			wantContains:    []string{"updatedAt", "name"},
			wantNotContains: []string{"@minVersion", "__typename"},
		},
		{
			name:   "failing minVersion is removed",
			server: "16.10",
			input: `query Q {
  project {
    updatedAt @minVersion(v: "17.0")
    name
  }
}`,
			wantContains:    []string{"name"},
			wantNotContains: []string{"updatedAt", "@minVersion"},
		},
		{
			name:   "exceeded maxVersion is removed",
			server: "17.5",
			input: `query Q {
  project {
    legacyField @maxVersion(v: "17.0")
    name
  }
}`,
			wantContains:    []string{"name"},
			wantNotContains: []string{"legacyField", "@maxVersion"},
		},
		{
			name:   "in range with both bounds is kept",
			server: "17.2",
			input: `query Q {
  project {
    middleField @minVersion(v: "17.0") @maxVersion(v: "18.0")
    name
  }
}`,
			wantContains:    []string{"middleField", "name"},
			wantNotContains: []string{"@minVersion", "@maxVersion"},
		},
		{
			name:   "below min with both bounds is removed",
			server: "16.0",
			input: `query Q {
  project {
    middleField @minVersion(v: "17.0") @maxVersion(v: "18.0")
    name
  }
}`,
			wantContains:    []string{"name"},
			wantNotContains: []string{"middleField"},
		},
		{
			name:   "above max with both bounds is removed",
			server: "19.0",
			input: `query Q {
  project {
    middleField @minVersion(v: "17.0") @maxVersion(v: "18.0")
    name
  }
}`,
			wantContains:    []string{"name"},
			wantNotContains: []string{"middleField"},
		},
		{
			name:   "fragment spread failing version is removed",
			server: "16.0",
			input: `query Q {
  project {
    ...NewProjectFields @minVersion(v: "17.0")
    name
  }
}
fragment NewProjectFields on Project {
  id
}`,
			// The definition goes too: nothing spreads it any more, and
			// GraphQL rejects a document declaring an unused fragment.
			wantContains:    []string{"name"},
			wantNotContains: []string{"NewProjectFields", "@minVersion"},
		},
		{
			name:   "fragment spread passing version is kept",
			server: "17.0",
			input: `query Q {
  project {
    ...NewProjectFields @minVersion(v: "17.0")
    name
  }
}
fragment NewProjectFields on Project {
  id
}`,
			wantContains:    []string{"NewProjectFields", "name"},
			wantNotContains: []string{"@minVersion"},
		},
		{
			name:   "inline fragment failing version is removed",
			server: "16.0",
			input: `query Q {
  project {
    ... on Project @minVersion(v: "17.0") {
      newField
    }
    name
  }
}`,
			wantContains:    []string{"name"},
			wantNotContains: []string{"newField"},
		},
		{
			name:   "inline fragment passing version is kept",
			server: "17.0",
			input: `query Q {
  project {
    ... on Project @minVersion(v: "17.0") {
      newField
    }
    name
  }
}`,
			wantContains:    []string{"newField", "name"},
			wantNotContains: []string{"@minVersion"},
		},
		{
			name:   "empty parent selection injects __typename",
			server: "16.0",
			input: `query Q {
  project {
    parent {
      onlyChild @minVersion(v: "17.0")
    }
  }
}`,
			wantContains:    []string{"__typename"},
			wantNotContains: []string{"onlyChild"},
		},
		{
			name:   "fully filtered fragment injects __typename",
			server: "16.0",
			input: `query Q {
  ...Frag
}
fragment Frag on Query {
  gated @minVersion(v: "17.0")
}`,
			wantContains:    []string{"__typename"},
			wantNotContains: []string{"gated"},
		},
		{
			name:   "co-existing @include is preserved when version passes",
			server: "17.0",
			input: `query Q($flag: Boolean!) {
  project {
    field @include(if: $flag) @minVersion(v: "17.0")
  }
}`,
			wantContains:    []string{"@include", "$flag", "field"},
			wantNotContains: []string{"@minVersion"},
		},
		{
			name:   "co-existing @include is removed when version fails",
			server: "16.0",
			input: `query Q($flag: Boolean!) {
  project {
    field @include(if: $flag) @minVersion(v: "17.0")
    other
  }
}`,
			wantContains:    []string{"other"},
			wantNotContains: []string{"field", "@include"},
			// $flag was only referenced by the dropped field's @include.
			wantPruned: []string{"flag"},
		},
		{
			name:   "field dropped, variable used only there is pruned",
			server: "16.0",
			input: `query Q($cursor: String, $alive: Boolean!) {
  project {
    gated(after: $cursor) @minVersion(v: "17.0")
    alwaysHere(flag: $alive)
  }
}`,
			wantPruned: []string{"cursor"},
		},
		{
			name:   "variable used in surviving and dropped argument is retained",
			server: "16.0",
			input: `query Q($cursor: String) {
  project {
    gated(after: $cursor) @minVersion(v: "17.0")
    alsoHere(after: $cursor)
  }
}`,
			wantPruned: []string{},
		},
		{
			name:   "no version markers means no rewriting",
			server: "17.0",
			input: `query Q {
  project {
    name
  }
}`,
			wantContains:    []string{"name"},
			wantNotContains: []string{"__typename"},
			wantPruned:      []string{},
		},
		{
			name:   "malformed version string in marker errors out",
			server: "17.0",
			input: `query Q {
  project {
    field @minVersion(v: "not-a-version")
  }
}`,
			wantErr: true,
		},
		{
			name:   "duplicate @minVersion errors out",
			server: "17.0",
			input: `query Q {
  project {
    field @minVersion(v: "17.0") @minVersion(v: "18.0")
  }
}`,
			wantErr: true,
		},
		{
			name:   "aliased field with passing minVersion keeps alias",
			server: "17.0",
			input: `query Q {
  project {
    aliasName: updatedAt @minVersion(v: "17.0")
    name
  }
}`,
			wantContains: []string{"aliasName", "updatedAt", "name"},
		},
		{
			name:   "nested depth pruning",
			server: "16.0",
			input: `query Q {
  l1 {
    l2 {
      l3 {
        deepGated @minVersion(v: "17.0")
        sibling
      }
    }
  }
}`,
			wantContains:    []string{"sibling"},
			wantNotContains: []string{"deepGated"},
		},
		{
			// A kept inline fragment must still have its inner selections
			// filtered. Distinct from the whole-inline-fragment drop/keep
			// cases above, which never exercise the recursion.
			name:   "kept inline fragment still filters inner fields",
			server: "17.0",
			input: `query Q {
  project {
    ... on Project @minVersion(v: "16.0") {
      innerKept
      innerGated @minVersion(v: "18.0")
    }
  }
}`,
			wantContains:    []string{"innerKept"},
			wantNotContains: []string{"innerGated", "@minVersion"},
		},
		{
			// A fragment definition containing a spread, where the spread's
			// gate fails: the spread goes, the surrounding fragment body stays.
			name:   "nested fragments with mixed gating",
			server: "16.0",
			input: `query Q {
  project {
    ...Outer
  }
}
fragment Outer on Project {
  outerKept
  ...Inner @minVersion(v: "17.0")
}
fragment Inner on Project {
  innerField
}`,
			// Outer survives; Inner's only spread was gated out, so both the
			// spread and Inner's now-orphaned definition are removed.
			// Note the space: the formatter emits spreads as "... Name", so
			// asserting on "...Inner" would pass trivially.
			wantContains:    []string{"outerKept", "fragment Outer"},
			wantNotContains: []string{"... Inner", "fragment Inner", "@minVersion"},
		},
		{
			// Only ProjectB is affected; ProjectA must come through untouched.
			name:   "multiple fragments, only some affected",
			server: "16.0",
			input: `query Q {
  project {
    ...ProjectA
    ...ProjectB
  }
}
fragment ProjectA on Project {
  aOne
  aTwo
}
fragment ProjectB on Project {
  bKept
  bGated @minVersion(v: "17.0")
}`,
			wantContains:    []string{"aOne", "aTwo", "bKept"},
			wantNotContains: []string{"bGated", "@minVersion"},
		},
		{
			// Gated at field level AND at the enclosing inline-fragment level.
			// The outer gate fails, so the whole subtree goes regardless of the
			// inner field's own (passing) gate.
			name:   "field gated at both field and enclosing inline fragment level",
			server: "16.0",
			input: `query Q {
  project {
    ... on Project @minVersion(v: "17.0") {
      doubleGated @minVersion(v: "15.0")
    }
    sibling
  }
}`,
			wantContains:    []string{"sibling"},
			wantNotContains: []string{"doubleGated", "@minVersion"},
		},
		{
			name:   "duplicate @maxVersion errors out",
			server: "17.0",
			input: `query Q {
  project {
    field @maxVersion(v: "18.0") @maxVersion(v: "19.0")
  }
}`,
			wantErr: true,
		},
		{
			name:   "marker missing required v argument errors out",
			server: "17.0",
			input: `query Q {
  project {
    field @minVersion
  }
}`,
			wantErr: true,
		},
		{
			name:   "marker with non-string v argument errors out",
			server: "17.0",
			input: `query Q {
  project {
    field @minVersion(v: 17)
  }
}`,
			wantErr: true,
		},
		{
			name:    "unparseable query errors out",
			server:  "17.0",
			input:   `query Q { project {`,
			wantErr: true,
		},
		{
			name:   "maxVersion on a fragment spread is honored",
			server: "18.0",
			input: `query Q {
  project {
    ...Legacy @maxVersion(v: "17.0")
    kept
  }
}
fragment Legacy on Project {
  legacyField
}`,
			wantContains:    []string{"kept"},
			wantNotContains: []string{"... Legacy", "@maxVersion"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGate(tt.server)
			if err != nil {
				t.Fatalf("NewGate(%q): %v", tt.server, err)
			}
			out, pruned, err := stripVersionedNodes(tt.input, g)
			if (err != nil) != tt.wantErr {
				t.Fatalf("stripVersionedNodes error=%v wantErr=%v\noutput:\n%s",
					err, tt.wantErr, out)
			}
			if tt.wantErr {
				return
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

func keysOf(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

// TestStripVersionedNodes_OrphanedFragmentTakesItsVariables covers the cascade
// that a dropped spread sets off. Removing the last spread of a fragment
// orphans its definition, which GraphQL rejects:
//
//	Fragment HeavyFields was defined, but not used
//
// Removing that definition in turn strands any variable only its body
// referenced, which GraphQL also rejects:
//
//	Variable $filter is declared by Q but not used
//
// Both must be cleaned up in one pass.
func TestStripVersionedNodes_OrphanedFragmentTakesItsVariables(t *testing.T) {
	input := `query Q($filter: String) {
  project {
    ...HeavyFields @minVersion(v: "17.0")
    name
  }
}
fragment HeavyFields on Project {
  related(filter: $filter) {
    id
  }
}`
	g, err := NewGate("16.0")
	if err != nil {
		t.Fatal(err)
	}
	out, pruned, err := stripVersionedNodes(input, g)
	if err != nil {
		t.Fatalf("rewrite: %v", err)
	}
	if _, dropped := pruned["filter"]; !dropped {
		t.Errorf("$filter should be pruned: its only reference was the orphaned fragment\npruned=%v\noutput:\n%s",
			keysOf(pruned), out)
	}
	if strings.Contains(out, "$filter") {
		t.Errorf("stranded variable declaration left in output\noutput:\n%s", out)
	}
	if strings.Contains(out, "HeavyFields") {
		t.Errorf("orphaned fragment definition left in output\noutput:\n%s", out)
	}
}

// TestStripVersionedNodes_FragmentKeptWhenSpreadSurvives is the counterpart:
// a fragment reachable through any surviving spread must be retained, along
// with the variables its body references.
func TestStripVersionedNodes_FragmentKeptWhenSpreadSurvives(t *testing.T) {
	input := `query Q($filter: String) {
  project {
    ...HeavyFields @minVersion(v: "17.0")
    other {
      ...HeavyFields
    }
  }
}
fragment HeavyFields on Project {
  related(filter: $filter) {
    id
  }
}`
	g, err := NewGate("16.0")
	if err != nil {
		t.Fatal(err)
	}
	out, pruned, err := stripVersionedNodes(input, g)
	if err != nil {
		t.Fatalf("rewrite: %v", err)
	}
	if _, dropped := pruned["filter"]; dropped {
		t.Errorf("$filter must stay: a surviving spread still reaches its fragment\noutput:\n%s", out)
	}
	if !strings.Contains(out, "fragment HeavyFields") {
		t.Errorf("fragment must be retained while any spread survives\noutput:\n%s", out)
	}
	if !strings.Contains(out, "$filter") {
		t.Errorf("expected $filter to remain declared\noutput:\n%s", out)
	}
}

// TestStripVersionedNodes_TransitiveOrphanRemoval checks the walk is
// transitive: dropping A's only spread must also remove B, which only A
// spread, and C, which only B spread.
func TestStripVersionedNodes_TransitiveOrphanRemoval(t *testing.T) {
	input := `query Q {
  project {
    ...A @minVersion(v: "17.0")
    name
  }
}
fragment A on Project {
  ...B
}
fragment B on Project {
  ...C
}
fragment C on Project {
  id
}`
	g, err := NewGate("16.0")
	if err != nil {
		t.Fatal(err)
	}
	out, _, err := stripVersionedNodes(input, g)
	if err != nil {
		t.Fatalf("rewrite: %v", err)
	}
	for _, name := range []string{"fragment A", "fragment B", "fragment C"} {
		if strings.Contains(out, name) {
			t.Errorf("%s should have been removed transitively\noutput:\n%s", name, out)
		}
	}
	if !strings.Contains(out, "name") {
		t.Errorf("ungated sibling must survive\noutput:\n%s", out)
	}
}

// TestStripVersionedNodes_GatedFieldOrphansNestedSpread reproduces the exact
// shape that broke against a live 15.11 instance: a gated non-leaf field whose
// sub-selection held the only spread of a fragment.
//
//	runner @minVersion(v: "16.0") { ...RunnerReferenceFields }
func TestStripVersionedNodes_GatedFieldOrphansNestedSpread(t *testing.T) {
	input := `query Q {
  job {
    name
    runner @minVersion(v: "16.0") {
      ...RunnerReferenceFields
    }
  }
}
fragment RunnerReferenceFields on CiRunner {
  id
  shortSha
}`
	g, err := NewGate("15.11.13-ee")
	if err != nil {
		t.Fatal(err)
	}
	out, _, err := stripVersionedNodes(input, g)
	if err != nil {
		t.Fatalf("rewrite: %v", err)
	}
	if strings.Contains(out, "runner") {
		t.Errorf("gated field must be dropped\noutput:\n%s", out)
	}
	if strings.Contains(out, "RunnerReferenceFields") {
		t.Errorf("fragment orphaned by the dropped field must be removed\noutput:\n%s", out)
	}
	if !strings.Contains(out, "name") {
		t.Errorf("ungated sibling must survive\noutput:\n%s", out)
	}
}

// TestStripVersionedNodes_VariableInsideNestedObjectValue covers the
// argument-value walker's handling of nested object/list values that
// reference variables.
func TestStripVersionedNodes_VariableInsideNestedObjectValue(t *testing.T) {
	input := `query Q($id: ID!) {
  search(filter: {parentId: $id, kind: "x"}) {
    name
    gatedField @minVersion(v: "17.0")
  }
}`
	g, err := NewGate("16.0")
	if err != nil {
		t.Fatal(err)
	}
	out, pruned, err := stripVersionedNodes(input, g)
	if err != nil {
		t.Fatalf("rewrite: %v", err)
	}
	if _, dropped := pruned["id"]; dropped {
		t.Errorf("$id should be alive (used in surviving nested object arg)\noutput:\n%s", out)
	}
}

// TestStripForVersion covers the exported wrapper used by callers that need
// the rewritten query text rather than an HTTP round trip.
func TestStripForVersion(t *testing.T) {
	const input = `query Q {
  project {
    name
    updatedAt @minVersion(v: "17.0")
    exitCode @minVersion(v: "17.7")
  }
}`

	tests := []struct {
		name            string
		server          string
		wantContains    []string
		wantNotContains []string
	}{
		{
			name:            "old server drops both gated fields",
			server:          "16.10.0",
			wantContains:    []string{"name"},
			wantNotContains: []string{"updatedAt", "exitCode", "@minVersion"},
		},
		{
			// The boundary that matters in practice: 17.4 keeps updatedAt but
			// must still drop exitCode.
			name:            "mid server keeps one and drops the other",
			server:          "17.4.0-ee",
			wantContains:    []string{"name", "updatedAt"},
			wantNotContains: []string{"exitCode", "@minVersion"},
		},
		{
			name:            "new server keeps both",
			server:          "17.7.0",
			wantContains:    []string{"name", "updatedAt", "exitCode"},
			wantNotContains: []string{"@minVersion"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGate(tt.server)
			if err != nil {
				t.Fatalf("NewGate(%q): %v", tt.server, err)
			}
			out, err := StripForVersion(input, g)
			if err != nil {
				t.Fatalf("StripForVersion: %v", err)
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
		})
	}
}

// TestStripForVersion_RemovesFieldNotJustDirective pins the distinction that
// motivates this function: a directive-only strip would leave the gated field
// in place, which the server then rejects as an unknown field.
func TestStripForVersion_RemovesFieldNotJustDirective(t *testing.T) {
	const input = `query Q {
  job {
    name
    exitCode @minVersion(v: "17.7")
  }
}`
	g, err := NewGate("17.4.0")
	if err != nil {
		t.Fatal(err)
	}
	out, err := StripForVersion(input, g)
	if err != nil {
		t.Fatalf("StripForVersion: %v", err)
	}
	if strings.Contains(out, "exitCode") {
		t.Errorf("gated field must be removed, not merely un-annotated\noutput:\n%s", out)
	}
	if !strings.Contains(out, "name") {
		t.Errorf("ungated sibling must survive\noutput:\n%s", out)
	}
}

// TestStripForVersion_PropagatesError asserts a malformed marker surfaces as
// an error rather than a silently unmodified query.
func TestStripForVersion_PropagatesError(t *testing.T) {
	g, err := NewGate("17.0")
	if err != nil {
		t.Fatal(err)
	}
	out, err := StripForVersion(`query Q { x @minVersion(v: "not-a-version") }`, g)
	if err == nil {
		t.Fatalf("expected an error, got output:\n%s", out)
	}
	if out != "" {
		t.Errorf("expected empty output on error, got: %q", out)
	}
}
