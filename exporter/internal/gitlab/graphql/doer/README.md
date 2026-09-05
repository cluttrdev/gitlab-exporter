# `internal/gitlab/graphql/doer` - runtime query rewriter

This package adds runtime version gating to genqlient-generated GraphQL
queries. It is wired in front of the HTTP transport used by the GraphQL
client. On every outgoing request it parses the query, drops any selections
whose `@minVersion` / `@maxVersion` directives exclude the connected GitLab
server's version, prunes now-unused query variables, and forwards the
rewritten request unchanged in every other respect.

The motivation is to let one query target multiple GitLab versions without
maintaining duplicate per-version variants in `.graphql` files or Go
dispatch code.

## Annotation form

Annotate any field, fragment spread, or inline fragment in the project's
`.graphql` files with the directives declared in
`internal/gitlab/graphql/spec/directives.graphql`:

```graphql
directive @minVersion(v: String!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT
directive @maxVersion(v: String!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT
```

Examples:

```graphql
fragment ProjectFieldsCore on Project {
  name
  updatedAt @minVersion(v: "17.0")
  legacyField @maxVersion(v: "17.0")
}

query getProjects {
  projects {
    nodes {
      ...ProjectFieldsCore
      ...NewProjectFields @minVersion(v: "17.0")
      ... on Project @minVersion(v: "17.0") {
        modernOnlyField
      }
    }
  }
}
```

Multiple directives compose:

```graphql
field @minVersion(v: "16.0") @maxVersion(v: "18.0")
```

The directives are declared in the schema input list of both
`genqlient.yaml` (used by `make graphql`) and `.graphqlrc.yaml` (used by
editor/LSP tooling), so misspellings (`@minVersio(...)`) are caught at
`make graphql` time rather than becoming silent runtime no-ops.

The directives never reach the server, the rewriter strips them from
every outgoing request.

## Pointer convention for gated fields

Any field gated by `@minVersion` or `@maxVersion` may be absent at runtime
(stripped by the rewriter when the gate disallows it). The corresponding
Go type must therefore be a pointer so that "absent" is distinguishable
from "zero value." Add an adjacent magic comment immediately above the
gated field:

```graphql
# @genqlient(pointer: true)
updatedAt @minVersion(v: "17.0")
```

Use `pointer: true` **only**. `omitempty` applies to genqlient *variables*,
not response fields; adding it to a field fails codegen with:

```
omitempty is not applicable to variables, not fields
```

Without the pointer, genqlient emits a non-pointer field whose zero value
would ambiguously stand for both "field returned but empty" and "field
stripped by the rewriter."

The project's `optional: pointer` setting in `genqlient.yaml` already makes
every *nullable* schema field a pointer, so the comment is redundant for
those (it is a no-op on `updatedAt` today). It is not redundant for
non-nullable fields. Default to including it on any gated field: it states
the intent explicitly rather than relying on schema nullability plus a
global config setting.

## How variables are pruned

When the rewriter drops a field, it walks the field's arguments and the
arguments of any nested directives, recording the variable names they
reference as pruning candidates. After the whole document has been pruned,
it performs a second pass to find any variable still referenced by a
surviving selection. Variables in `candidates \ stillReferenced` are
removed from each operation's `variableDefinitions` and from the request's
`variables` payload (whether sent as a URL query parameter on a GET request
or as a JSON body field on a POST request).

A variable referenced from a fragment definition stays alive for as long as
that fragment is still reachable. If pruning removes the fragment's last
spread, the definition is removed too (see below) and its variables become
pruning candidates like any others.

## Orphaned fragment definitions

GraphQL rejects a document that declares a fragment nothing uses:

```
Fragment RunnerReferenceFields was defined, but not used
```

Pruning can cause exactly that, most easily by gating a field whose
sub-selection holds the only spread of a fragment:

```graphql
runner @minVersion(v: "16.0") {
  ...RunnerReferenceFields
}
```

Below 16.0 the `runner` field and its subtree disappear, leaving
`fragment RunnerReferenceFields` unreferenced. So after pruning, the
rewriter recomputes which fragments are still reachable from the operations
- following spreads transitively, since removing one fragment can orphan
another that only it spread - and drops the rest.

Removing a fragment can in turn strand a variable that only its body
referenced, which the server also rejects:

```
Variable $filter is declared by Q but not used
```

The variables of a removed fragment therefore become pruning candidates in
the same pass, and are dropped unless something surviving still references
them.

This matters most when annotating a **non-leaf** node. Gating a plain scalar
leaf cannot orphan anything; gating a field with a sub-selection, a fragment
spread, or an inline fragment can.

## What happens to empty selection sets

GraphQL forbids empty selection sets on composite types. If pruning empties
a parent's selection set, the rewriter injects `__typename` to keep the
query syntactically valid. Converters can treat the resulting "present but
trivial" parent the same way they treat a nil pointer.

A fully filtered fragment definition (every selection gated out) is
handled the same way.

## Failure semantics

Any rewriter error — parse failure, malformed version string in a
directive `v:` argument, duplicate `@minVersion` on the same node — aborts
the request: the Doer returns the error and never forwards the original
(unrewritten) query.

A startup-time failure to install the gate is fatal in `cmd/run.go` (the
hardcoded fallback version is build-time configuration; a parse failure
there is a code bug).

## `@include` / `@skip` below GitLab 17.0

Several operations split their field sets with conditional fragment spreads,
so one generated operation can be issued cheaply for just the fields a given
caller needs:

```graphql
...PipelineFieldsCore      @include(if: $_core)
...PipelineFieldsRelations @include(if: $_relations)
```

That only reduces cost if the server's complexity analyzer honors the
directive. Until **graphql-ruby 2.3.0** it did not. Its analysis visitor sets
a "skipping" flag from a node's directives, but only inside `on_field`;
`on_fragment_spread` and `on_inline_fragment` gained the same handling in
2.3.0. Only GitLab 17.0 ships 2.3.2 so before this release `queryComplexity`
charges for conditional fragments the server never resolves.

So below `complexityHonorsConditionalFragments` (17.0) the rewriter evaluates
`@include` / `@skip` itself: nodes whose condition excludes them are removed
outright, and the directive is stripped from nodes that stay.

From 17.0 up, nothing happens, the directives go out untouched and the server
evaluates them, so its own accounting stays authoritative.

## Caching

The Doer caches rewritten queries by raw query string. The set of distinct
genqlient-emitted query strings is bounded by the number of operations, so the
cache is naturally small and stable. Calling `SetGate` (e.g. because a new
GitLab server version was detected) clears the cache.

On servers that need conditional materialization the rewrite also depends on
the request's variables, so the key gains the sorted set of **boolean**
variables. An `if:` condition is a boolean, so no other variable type can
change the result and none is included. The key space stays bounded by
"operations times flag combinations".

## Argument-level gating

`@minVersion` / `@maxVersion` drop a whole selection node. Sometimes a field
exists on older servers but one of its *arguments* was added later - for
example `MergeRequest.notes` accepts `filter:` only from 16.5. Dropping the
whole field would lose the notes entirely; dropping just the argument returns
all notes instead of activity-only, a graceful degradation.

GraphQL forbids directives on argument usages (there is no executable
`DirectiveLocation` for them, so `notes(filter: X @minVersion(...))` is a
syntax error). Two `on FIELD` directives express argument gating instead, with
the affected arguments listed by name:

```graphql
directive @minVersionArgs(v: String!, args: [String!]!) on FIELD
directive @maxVersionArgs(v: String!, args: [String!]!) on FIELD
```

```graphql
notes(filter: ONLY_ACTIVITY) @minVersionArgs(v: "16.5", args: ["filter"])
```

When the server version is out of range, the named arguments are removed from
the field while the field and its selection set stay. Any variable referenced
only by a removed argument is pruned by the same machinery that prunes
variables of dropped fields. The markers are stripped before the request is
sent. The check runs only after the field survives its own `@minVersion` /
`@maxVersion` - a dropped field makes its argument gating moot. Works
identically for literal and variable argument values.

Two rules that will otherwise cause issues:

- **Only gate additive/optional arguments.** The rewriter has no runtime schema
  and cannot tell a refinement (`filter:`) from a semantically required
  selector. Removing a required argument produces a query that fails or, worse,
  silently changes meaning.
- **One marker per field.** A single version covers every name in `args`. Two
  arguments on one field needing *different* versions is not expressible: a
  second marker of the same kind is rejected as an authoring error, and
  gqlparser's uniqueness rule
  (`validator/rules/unique_directives_per_location.go`) would reject it at
  codegen even if the directive were declared `repeatable`, because it compares
  the directive name against the literal string `"repeatable"` rather than
  consulting `IsRepeatable`.

The rewriter reports an authoring error for: an `args` name not present on the
field (caught on every server, since `args` entries are otherwise unvalidated
strings), an empty `args` list, a duplicate marker of the same kind, a
malformed or missing version, or a non-string list item.

## Workflow for adding new gated fields

1. Edit the `.graphql` file under `spec/fragments/` or `spec/queries/`,
   adding `@minVersion(v: "X.Y")` / `@maxVersion(v: "X.Y")` to the field,
   spread, or inline fragment. Add a `# @genqlient(pointer: true)` comment
   above any newly gated field unless it is already nullable. To gate only
   an argument of a field that itself exists on older servers, use
   `@minVersionArgs(v: "X.Y", args: ["argName"])` /
   `@maxVersionArgs(...)` instead (see "Argument-level gating").
2. Run `make graphql` from `exporter/`. Genqlient validates the directives
   against `spec/directives.graphql` and regenerates `generated.go`.
3. Add a unit test in `doer/rewrite_test.go` (or a Doer-level test in
   `doer/doer_test.go`) that exercises the new gating.
4. If the field is version-sensitive on the wire, add its name and
   introducing major version to `fieldsUnknownBefore` in
   `version_gating_test.go`. The simulated server then rejects any request
   that fails to strip it, so a gating regression fails the build.
5. Confirm `go test ./...` is green.
