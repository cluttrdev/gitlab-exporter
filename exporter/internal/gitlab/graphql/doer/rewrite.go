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

// Names of the standard conditional directives. On servers whose complexity
// analyzer ignores them on fragment spreads and inline fragments, the rewriter
// evaluates them itself and removes the excluded nodes outright.
const (
	directiveInclude      = "include"
	directiveSkip         = "skip"
	directiveConditionArg = "if"
)

// StripForVersion rewrites query as the VersionGatedDoer would for the given
// gate: every selection node excluded by its @minVersion / @maxVersion
// directives is removed, and the marker directives are stripped from the nodes
// that remain.
//
// @include / @skip are left untouched. Use StripForVersionWithVariables to
// also materialize those against a concrete set of variables.
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

// StripForVersionWithVariables rewrites query exactly as the VersionGatedDoer
// would for the given gate and request variables. In addition to the
// @minVersion / @maxVersion pruning StripForVersion performs, it materializes
// @include / @skip when the gate's server needs it
// (see Gate.MaterializesConditionals): every node whose condition resolves to
// "excluded" is removed, and the directive is stripped from nodes that stay.
//
// Callers measuring query cost should prefer this over StripForVersion: below
// GitLab 17.0 it is the difference between the query the server is charged for
// and the query it actually resolves.
func StripForVersionWithVariables(query string, gate *Gate, vars map[string]any) (string, error) {
	rewritten, _, err := stripNodes(query, rewriteOpts{
		gate:        gate,
		vars:        vars,
		materialize: gate.MaterializesConditionals(),
	})
	if err != nil {
		return "", err
	}
	return rewritten, nil
}

// rewriteOpts configures a single rewrite pass.
type rewriteOpts struct {
	gate *Gate

	// vars carries the request's GraphQL variables, used to resolve the
	// `if:` argument of @include / @skip. Only consulted when materialize
	// is set.
	vars map[string]any

	// materialize enables client-side evaluation of @include / @skip. When
	// false the rewriter leaves those directives alone and the server decides.
	materialize bool
}

// pruneCtx is the resolved, per-document form of rewriteOpts handed to the
// selection walk.
type pruneCtx struct {
	gate *Gate

	// conditionals maps variable name to the boolean value the request
	// supplies for it. A nil map disables conditional materialization
	// entirely; a name missing from a non-nil map is unresolvable and makes
	// the rewriter leave the node - and its directive - in place.
	conditionals map[string]bool
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
	return stripNodes(raw, rewriteOpts{gate: gate})
}

// stripNodes is the full rewrite pass behind stripVersionedNodes and
// StripForVersionWithVariables. See stripVersionedNodes for the contract; the
// only addition here is the optional @include / @skip materialization
// configured by opts.
func stripNodes(raw string, opts rewriteOpts) (string, map[string]struct{}, error) {
	doc, perr := parser.ParseQuery(&ast.Source{Name: "request", Input: raw})
	if perr != nil {
		return "", nil, fmt.Errorf("parse query: %w", perr)
	}

	pc := pruneCtx{
		gate:         opts.gate,
		conditionals: resolveConditionals(doc, opts),
	}

	candidateVars := make(map[string]struct{})

	for _, frag := range doc.Fragments {
		pruned, err := pruneSelectionSet(pc, frag.SelectionSet, candidateVars)
		if err != nil {
			return "", nil, fmt.Errorf("fragment %s: %w", frag.Name, err)
		}
		if len(pruned) == 0 {
			pruned = injectTypename()
		}
		frag.SelectionSet = pruned
	}

	for _, op := range doc.Operations {
		pruned, err := pruneSelectionSet(pc, op.SelectionSet, candidateVars)
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

// resolveConditionals builds the variable -> boolean map used to materialize
// @include / @skip, or returns nil to leave those directives to the server.
//
// Values come from the request variables, falling back to the default in the
// operation's variable definition. A variable that resolves to neither is
// omitted, which makes keepFor leave its node untouched.
//
// Materialization is restricted to single-operation documents. Fragment
// definitions are shared across every operation in a document, so with more
// than one there is no single correct value to prune a fragment body against.
// Every genqlient-issued request carries exactly one operation.
func resolveConditionals(doc *ast.QueryDocument, opts rewriteOpts) map[string]bool {
	if !opts.materialize || len(doc.Operations) != 1 {
		return nil
	}

	out := make(map[string]bool)
	for _, vd := range doc.Operations[0].VariableDefinitions {
		if v, ok := opts.vars[vd.Variable]; ok {
			if b, ok := v.(bool); ok {
				out[vd.Variable] = b
			}
			// A non-boolean value for a variable used in an `if:` position is
			// a caller error; leaving it unresolved makes the server say so.
			continue
		}
		if vd.DefaultValue != nil && vd.DefaultValue.Kind == ast.BooleanValue {
			out[vd.Variable] = vd.DefaultValue.Raw == "true"
		}
	}
	return out
}

// pruneSelectionSet walks a selection set and returns a new SelectionSet with
// nodes dropped according to their @minVersion / @maxVersion directives and,
// when pc.conditionals is set, their @include / @skip directives. Dropped
// subtrees contribute variable references to candidateVars so the second
// pass can decide which variable definitions are truly unused.
func pruneSelectionSet(pc pruneCtx, in ast.SelectionSet, candidateVars map[string]struct{}) (ast.SelectionSet, error) {
	if len(in) == 0 {
		return nil, nil
	}
	out := make(ast.SelectionSet, 0, len(in))
	for _, sel := range in {
		switch n := sel.(type) {
		case *ast.Field:
			keep, stripped, err := keepFor(pc, n.Directives, candidateVars)
			if err != nil {
				return nil, fmt.Errorf("field %q: %w", n.Name, err)
			}
			if !keep {
				collectVariableReferencesField(n, candidateVars)
				continue
			}
			child, err := pruneSelectionSet(pc, n.SelectionSet, candidateVars)
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
			keep, stripped, err := keepFor(pc, n.Directives, candidateVars)
			if err != nil {
				return nil, fmt.Errorf("spread %q: %w", n.Name, err)
			}
			if !keep {
				// A spread has no arguments or sub-selections of its own, so
				// the only variables it can strand are the ones in its
				// directives - the `if:` of an @include, typically. The
				// fragment definition stays in the document for now;
				// reachableFragments removes it, and harvests its variables,
				// if this was its last spread.
				collectVariableReferencesDirectives(n.Directives, candidateVars)
				continue
			}
			n.Directives = stripped
			out = append(out, n)

		case *ast.InlineFragment:
			keep, stripped, err := keepFor(pc, n.Directives, candidateVars)
			if err != nil {
				return nil, fmt.Errorf("inline fragment on %q: %w", n.TypeCondition, err)
			}
			if !keep {
				collectVariableReferencesDirectives(n.Directives, candidateVars)
				collectVariableReferences(n.SelectionSet, candidateVars)
				continue
			}
			child, err := pruneSelectionSet(pc, n.SelectionSet, candidateVars)
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

// keepFor evaluates a node's directives and reports whether the node should
// be kept. Two independent decisions are folded in:
//
//   - @minVersion / @maxVersion against the gate's server version;
//   - @include / @skip against pc.conditionals, but only when that map is
//     non-nil (see resolveConditionals).
//
// Either one excluding the node drops it. The returned DirectiveList carries
// the original directives with every directive the rewriter consumed removed;
// anything it did not act on is preserved unchanged - including @include and
// @skip whose condition it could not resolve, which the server then evaluates
// as usual.
//
// Consuming a resolved @include / @skip can leave its `if:` variable
// unreferenced, so those names are added to candidateVars for the caller's
// variable-pruning pass to intersect.
//
// Multiple @minVersion or @maxVersion markers on the same node are an
// authoring error.
func keepFor(pc pruneCtx, in ast.DirectiveList, candidateVars map[string]struct{}) (bool, ast.DirectiveList, error) {
	if len(in) == 0 {
		return true, in, nil
	}

	var (
		minV, maxV        string
		hasMin, hasMax    bool
		conditionExcludes bool
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
		case directiveInclude, directiveSkip:
			cond, ok := readConditionArg(d, pc.conditionals)
			if !ok {
				// Unresolvable: no conditionals map, no `if:` argument, or a
				// variable the request did not supply a boolean for. Leave the
				// directive in place and let the server decide - pruning a node
				// the server would have kept silently loses data.
				stripped = append(stripped, d)
				continue
			}
			// @skip(if: true) and @include(if: false) both exclude.
			if (d.Name == directiveSkip) == cond {
				conditionExcludes = true
			}
			collectVariableReferencesDirectives(ast.DirectiveList{d}, candidateVars)
		default:
			stripped = append(stripped, d)
		}
	}

	if conditionExcludes {
		return false, nil, nil
	}

	if !hasMin && !hasMax {
		return true, stripped, nil
	}

	keep, err := pc.gate.Allows(minV, maxV)
	if err != nil {
		return false, nil, err
	}
	return keep, stripped, nil
}

// readConditionArg resolves the `if:` argument of an @include / @skip
// directive to a concrete boolean. It reports false as its second result when
// the condition cannot be resolved statically, which is the signal to leave
// the directive alone: a literal of the wrong kind, a missing argument, a
// variable absent from conditionals, or a nil conditionals map (materialization
// disabled).
func readConditionArg(d *ast.Directive, conditionals map[string]bool) (bool, bool) {
	if conditionals == nil {
		return false, false
	}
	for _, arg := range d.Arguments {
		if arg.Name != directiveConditionArg || arg.Value == nil {
			continue
		}
		switch arg.Value.Kind {
		case ast.BooleanValue:
			return arg.Value.Raw == "true", true
		case ast.Variable:
			v, ok := conditionals[arg.Value.Raw]
			return v, ok
		}
		return false, false
	}
	return false, false
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
