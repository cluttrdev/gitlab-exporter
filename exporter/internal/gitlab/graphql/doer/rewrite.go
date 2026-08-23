package doer

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
	"github.com/vektah/gqlparser/v2/parser"
)

// Names of the version-gating directives recognized by the rewriter. They are
// declared in spec/directives.graphql so genqlient validates them at codegen
// time; this package strips them before any request reaches the server.
const (
	directiveMinVersion = "minVersion"
	directiveMaxVersion = "maxVersion"
	directiveVersionArg = "v"
)

// StripForVersion rewrites query as the VersionGatedDoer would for the given
// gate: every selection node excluded by its @minVersion / @maxVersion
// directives is removed, and the marker directives are stripped from the nodes
// that remain.
//
// This function performs no variable pruning, so a variable left unreferenced
// by the rewrite stays in the caller's variables map. That is harmless for
// servers, which ignore variables an operation does not declare.
func StripForVersion(query string, gate *Gate) (string, error) {
	rewritten, _, err := stripVersionedNodes(query, gate)
	if err != nil {
		return "", err
	}
	return rewritten, nil
}

// stripVersionedNodes parses raw, prunes any selection node whose
// @minVersion / @maxVersion directives exclude the gate's server version,
// and returns the re-formatted query alongside the set of variables that
// became unreferenced as a result of pruning.
//
// Matched marker directives are removed from kept nodes; all other directives
// (notably @include for field-set splits) pass through unchanged.
//
// Variables referenced exclusively from pruned subtrees are removed from
// each operation's VariableDefinitions and returned in prunedVars so the
// caller can also strip them from the variables map.
//
// On any error (parse, format, malformed version string in a marker
// directive, duplicate marker on the same node), the rewriter returns an
// error and no rewritten string. The caller MUST NOT forward the original
// query in that case.
func stripVersionedNodes(raw string, gate *Gate) (string, map[string]struct{}, error) {
	doc, perr := parser.ParseQuery(&ast.Source{Name: "request", Input: raw})
	if perr != nil {
		return "", nil, fmt.Errorf("parse query: %w", perr)
	}

	candidateVars := make(map[string]struct{})

	for _, frag := range doc.Fragments {
		pruned, err := pruneSelectionSet(gate, frag.SelectionSet, candidateVars)
		if err != nil {
			return "", nil, fmt.Errorf("fragment %s: %w", frag.Name, err)
		}
		if len(pruned) == 0 {
			pruned = injectTypename()
		}
		frag.SelectionSet = pruned
	}

	for _, op := range doc.Operations {
		pruned, err := pruneSelectionSet(gate, op.SelectionSet, candidateVars)
		if err != nil {
			return "", nil, fmt.Errorf("operation %s: %w", op.Name, err)
		}
		if len(pruned) == 0 {
			pruned = injectTypename()
		}
		op.SelectionSet = pruned
	}

	// Drop fragment definitions that nothing spreads any more. GraphQL
	// rejects a document declaring an unused fragment:
	//
	//	Fragment RunnerReferenceFields was defined, but not used
	//
	// which happens whenever pruning removes a fragment's last spread - most
	// easily by gating a field whose sub-selection held that spread:
	//
	//	runner @minVersion(v: "16.0") { ...RunnerReferenceFields }
	//
	// This runs before the variable pass below so that variables referenced
	// only from a removed fragment are pruned along with it. Without that,
	// removing the fragment would strand its variable definitions and the
	// server would reject the query for a second reason:
	//
	//	Variable $filter is declared by Q but not used
	if len(doc.Fragments) > 0 {
		reachable := reachableFragments(doc)
		if len(reachable) != len(doc.Fragments) {
			kept := doc.Fragments[:0]
			for _, frag := range doc.Fragments {
				if _, ok := reachable[frag.Name]; !ok {
					// Its body is going away, so anything it referenced
					// becomes a pruning candidate.
					collectVariableReferences(frag.SelectionSet, candidateVars)
					continue
				}
				kept = append(kept, frag)
			}
			doc.Fragments = kept
		}
	}

	// Second pass: a candidate variable is only really pruned if no
	// surviving selection still references it. Walk the document one more
	// time and intersect.
	stillReferenced := make(map[string]struct{})
	for _, frag := range doc.Fragments {
		collectVariableReferences(frag.SelectionSet, stillReferenced)
	}
	for _, op := range doc.Operations {
		collectVariableReferences(op.SelectionSet, stillReferenced)
	}

	prunedVars := make(map[string]struct{})
	for name := range candidateVars {
		if _, alive := stillReferenced[name]; alive {
			continue
		}
		prunedVars[name] = struct{}{}
	}

	// Remove pruned variable definitions from each operation.
	for _, op := range doc.Operations {
		if len(prunedVars) == 0 || len(op.VariableDefinitions) == 0 {
			continue
		}
		filtered := op.VariableDefinitions[:0]
		for _, vd := range op.VariableDefinitions {
			if _, drop := prunedVars[vd.Variable]; drop {
				continue
			}
			filtered = append(filtered, vd)
		}
		op.VariableDefinitions = filtered
	}

	var buf strings.Builder
	formatter.NewFormatter(&buf).FormatQueryDocument(doc)
	out := buf.String()
	if out == "" {
		return "", nil, fmt.Errorf("formatter produced empty output")
	}

	return out, prunedVars, nil
}

// pruneSelectionSet walks a selection set and returns a new SelectionSet with
// nodes dropped according to @minVersion / @maxVersion directives. Dropped
// subtrees contribute variable references to candidateVars so the second
// pass can decide which variable definitions are truly unused.
func pruneSelectionSet(gate *Gate, in ast.SelectionSet, candidateVars map[string]struct{}) (ast.SelectionSet, error) {
	if len(in) == 0 {
		return nil, nil
	}
	out := make(ast.SelectionSet, 0, len(in))
	for _, sel := range in {
		switch n := sel.(type) {
		case *ast.Field:
			keep, stripped, err := keepFor(gate, n.Directives)
			if err != nil {
				return nil, fmt.Errorf("field %q: %w", n.Name, err)
			}
			if !keep {
				collectVariableReferencesField(n, candidateVars)
				continue
			}
			child, err := pruneSelectionSet(gate, n.SelectionSet, candidateVars)
			if err != nil {
				return nil, fmt.Errorf("field %q: %w", n.Name, err)
			}
			// If the source had sub-selections but pruning emptied them,
			// inject __typename so the resulting query stays valid.
			if len(n.SelectionSet) > 0 && len(child) == 0 {
				child = injectTypename()
			}
			n.Directives = stripped
			n.SelectionSet = child
			out = append(out, n)

		case *ast.FragmentSpread:
			keep, stripped, err := keepFor(gate, n.Directives)
			if err != nil {
				return nil, fmt.Errorf("spread %q: %w", n.Name, err)
			}
			if !keep {
				// Spreads carry no arguments / sub-selections; the fragment
				// definition the spread targets stays in the document and
				// its variable references stay live.
				continue
			}
			n.Directives = stripped
			out = append(out, n)

		case *ast.InlineFragment:
			keep, stripped, err := keepFor(gate, n.Directives)
			if err != nil {
				return nil, fmt.Errorf("inline fragment on %q: %w", n.TypeCondition, err)
			}
			if !keep {
				collectVariableReferences(n.SelectionSet, candidateVars)
				continue
			}
			child, err := pruneSelectionSet(gate, n.SelectionSet, candidateVars)
			if err != nil {
				return nil, fmt.Errorf("inline fragment on %q: %w", n.TypeCondition, err)
			}
			if len(n.SelectionSet) > 0 && len(child) == 0 {
				child = injectTypename()
			}
			n.Directives = stripped
			n.SelectionSet = child
			out = append(out, n)
		}
	}
	return out, nil
}

// keepFor evaluates @minVersion / @maxVersion directives on a node and
// reports whether the node should be kept for the gate's server version.
// The returned DirectiveList carries the original directives with the
// matched marker directives removed; other directives (e.g. @include) are
// preserved unchanged.
//
// Multiple @minVersion or @maxVersion markers on the same node are an
// authoring error.
func keepFor(gate *Gate, in ast.DirectiveList) (bool, ast.DirectiveList, error) {
	if len(in) == 0 {
		return true, in, nil
	}

	var (
		minV, maxV     string
		hasMin, hasMax bool
	)
	stripped := make(ast.DirectiveList, 0, len(in))

	for _, d := range in {
		switch d.Name {
		case directiveMinVersion:
			if hasMin {
				return false, nil, fmt.Errorf("duplicate @%s directive", directiveMinVersion)
			}
			v, err := readVersionArg(d)
			if err != nil {
				return false, nil, err
			}
			hasMin = true
			minV = v
		case directiveMaxVersion:
			if hasMax {
				return false, nil, fmt.Errorf("duplicate @%s directive", directiveMaxVersion)
			}
			v, err := readVersionArg(d)
			if err != nil {
				return false, nil, err
			}
			hasMax = true
			maxV = v
		default:
			stripped = append(stripped, d)
		}
	}

	if !hasMin && !hasMax {
		return true, stripped, nil
	}

	keep, err := gate.Allows(minV, maxV)
	if err != nil {
		return false, nil, err
	}
	return keep, stripped, nil
}

// readVersionArg extracts the v: "..." argument from a @minVersion or
// @maxVersion directive. The value must be a literal StringValue; variables,
// enums, or any other value kind are rejected.
func readVersionArg(d *ast.Directive) (string, error) {
	for _, arg := range d.Arguments {
		if arg.Name != directiveVersionArg {
			continue
		}
		if arg.Value == nil {
			return "", fmt.Errorf("directive @%s: argument %q has no value", d.Name, directiveVersionArg)
		}
		if arg.Value.Kind != ast.StringValue && arg.Value.Kind != ast.BlockValue {
			return "", fmt.Errorf("directive @%s: argument %q must be a string, got kind %d", d.Name, directiveVersionArg, arg.Value.Kind)
		}
		return arg.Value.Raw, nil
	}
	return "", fmt.Errorf("directive @%s: missing required argument %q", d.Name, directiveVersionArg)
}

// reachableFragments returns the names of fragment definitions reachable from
// at least one operation, following spreads transitively: a fragment spread
// only by another unreachable fragment is itself unreachable.
//
// Reachability is computed from the already-pruned document, so a spread that
// version-gating removed does not keep its definition alive.
func reachableFragments(doc *ast.QueryDocument) map[string]struct{} {
	byName := make(map[string]*ast.FragmentDefinition, len(doc.Fragments))
	for _, frag := range doc.Fragments {
		byName[frag.Name] = frag
	}

	reachable := make(map[string]struct{}, len(doc.Fragments))

	var pending []string
	for _, op := range doc.Operations {
		pending = append(pending, collectSpreadNames(op.SelectionSet)...)
	}

	// Breadth-first over spread names. The visited check both terminates the
	// walk and guards against fragment cycles, which GraphQL forbids but
	// which must not hang the rewriter if one slips through.
	for len(pending) > 0 {
		name := pending[len(pending)-1]
		pending = pending[:len(pending)-1]

		if _, seen := reachable[name]; seen {
			continue
		}
		reachable[name] = struct{}{}

		if frag, ok := byName[name]; ok {
			pending = append(pending, collectSpreadNames(frag.SelectionSet)...)
		}
	}

	return reachable
}

// collectSpreadNames returns the names of every fragment spread in set,
// including those nested inside fields and inline fragments.
func collectSpreadNames(set ast.SelectionSet) []string {
	var names []string
	for _, sel := range set {
		switch n := sel.(type) {
		case *ast.FragmentSpread:
			names = append(names, n.Name)
		case *ast.Field:
			names = append(names, collectSpreadNames(n.SelectionSet)...)
		case *ast.InlineFragment:
			names = append(names, collectSpreadNames(n.SelectionSet)...)
		}
	}
	return names
}

// injectTypename returns a one-element SelectionSet containing a synthetic
// __typename field. Used when version-pruning empties a selection set whose
// source had at least one selection: GraphQL forbids empty selection sets
// on composite types, so __typename keeps the query syntactically valid.
func injectTypename() ast.SelectionSet {
	return ast.SelectionSet{
		&ast.Field{
			Alias: "__typename",
			Name:  "__typename",
		},
	}
}

// collectVariableReferences walks set and adds every variable-typed argument
// value (including nested object / list values, including directive
// arguments) to dst.
func collectVariableReferences(set ast.SelectionSet, dst map[string]struct{}) {
	for _, sel := range set {
		switch n := sel.(type) {
		case *ast.Field:
			collectVariableReferencesField(n, dst)
		case *ast.InlineFragment:
			collectVariableReferencesDirectives(n.Directives, dst)
			collectVariableReferences(n.SelectionSet, dst)
		case *ast.FragmentSpread:
			collectVariableReferencesDirectives(n.Directives, dst)
		}
	}
}

func collectVariableReferencesField(f *ast.Field, dst map[string]struct{}) {
	for _, arg := range f.Arguments {
		collectVariableReferencesValue(arg.Value, dst)
	}
	collectVariableReferencesDirectives(f.Directives, dst)
	collectVariableReferences(f.SelectionSet, dst)
}

func collectVariableReferencesDirectives(dirs ast.DirectiveList, dst map[string]struct{}) {
	for _, d := range dirs {
		for _, arg := range d.Arguments {
			collectVariableReferencesValue(arg.Value, dst)
		}
	}
}

func collectVariableReferencesValue(v *ast.Value, dst map[string]struct{}) {
	if v == nil {
		return
	}
	switch v.Kind {
	case ast.Variable:
		if v.Raw != "" {
			dst[v.Raw] = struct{}{}
		}
	case ast.ListValue, ast.ObjectValue:
		for _, child := range v.Children {
			collectVariableReferencesValue(child.Value, dst)
		}
	}
}
