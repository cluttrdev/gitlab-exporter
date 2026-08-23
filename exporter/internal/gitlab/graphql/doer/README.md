# `internal/gitlab/graphql/doer` — runtime query rewriter

This package adds runtime version gating to genqlient-generated GraphQL
queries. It is wired in front of the HTTP transport used by the GraphQL
client; on every outgoing request it parses the query, drops any selections
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
`make graphql` time rather than becoming silent runtime no-ops. Keep the two
files in sync.

The directives never reach the server; the rewriter strips them from
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

## Caching

The Doer caches rewritten queries by raw query string. The set of distinct
genqlient-emitted query strings is bounded by the number of operations, so the
cache is naturally small and stable. Calling `SetGate` (e.g. because a new
GitLab server version was detected) clears the cache.

## Argument-level gating (not supported)

The directives are declared on FIELD, FRAGMENT_SPREAD, and INLINE_FRAGMENT.
ARGUMENT_DEFINITION is intentionally omitted because no current annotation
needs argument-level gating. If a future change requires it, declare the
additional location in `spec/directives.graphql`, regenerate genqlient,
and extend `rewrite.go` to walk `ast.Argument` directives.

## Layout

```
exporter/internal/gitlab/graphql/
├── client.go                 # NewClient wraps the HTTP transport with
│                             # VersionGatedDoer and exposes SetVersion.
├── generated.go              # genqlient output (unchanged interface).
├── version_gating_test.go    # End-to-end: the real getProjects operation
│                             # against servers simulating several versions.
├── doer/
│   ├── doc.go                # Package overview.
│   ├── version.go            # ServerVersion + Gate + Allows.
│   ├── rewrite.go            # stripVersionedNodes (AST walker, var pruning).
│   ├── doer.go               # VersionGatedDoer (HTTP-level integration).
│   └── README.md             # This file.
└── spec/
    ├── directives.graphql        # @minVersion / @maxVersion declarations.
    ├── fragments/*.graphql       # Reusable fragment definitions; annotated
    │                             # with @minVersion / @maxVersion as needed.
    ├── queries/*.graphql         # Operation definitions; may also carry
    │                             # the directives on fields, spreads, and
    │                             # inline fragments.
    ├── schema/schema.graphql     # GitLab GraphQL schema (CI-only).
    └── complexity_test.{go,sh}   # Live complexity-measurement harness.
```

## Workflow for adding new gated fields

1. Edit the `.graphql` file under `spec/fragments/` or `spec/queries/`,
   adding `@minVersion(v: "X.Y")` / `@maxVersion(v: "X.Y")` to the field,
   spread, or inline fragment. Add a `# @genqlient(pointer: true)` comment
   above any newly gated field unless it is already nullable.
2. Run `make graphql` from `exporter/`. Genqlient validates the directives
   against `spec/directives.graphql` and regenerates `generated.go`.
3. Update any caller that previously had Go-side version dispatch for the
   same field (e.g. the pattern in `projects.go` was: `includeUpdatedAt := c.Version.Major >= 17`).
4. Add a unit test in `doer/rewrite_test.go` (or a Doer-level test in
   `doer/doer_test.go`) that exercises the new gating.
5. If the field is version-sensitive on the wire, add its name and
   introducing major version to `fieldsUnknownBefore` in
   `version_gating_test.go`. The simulated server then rejects any request
   that fails to strip it, so a gating regression fails the build.
6. Confirm `go test ./...` is green.
