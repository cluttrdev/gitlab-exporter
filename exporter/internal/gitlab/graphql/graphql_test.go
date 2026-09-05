package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/graphql/doer"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/version"
)

// gatedField describes a field that only exists on the server from a given
// GitLab version onward, and which therefore carries a @minVersion directive
// in the .graphql sources.
type gatedField struct {
	name  string // Field name as it appears in the query text.
	owner string // GraphQL type that owns it, used in the error message.
	since string // First version that knows the field.
}

// The fields currently gated in spec/fragments. Verified against a live
// GitLab 17.4.0-ee instance by schema introspection: updatedAt and closedAt
// exist there, exitCode does not.
var (
	projectUpdatedAt = gatedField{name: "updatedAt", owner: "Project", since: "17.0"}
	mrClosedAt       = gatedField{name: "closedAt", owner: "MergeRequest", since: "17.4"}
	jobExitCode      = gatedField{name: "exitCode", owner: "CiJob", since: "17.7"}
)

// simulatedGitLab serves a trivial response but rejects the request if the
// query mentions field.name while the simulated server version predates
// field.since, mirroring how the real API answers an unknown field:
//
//	{"errors":[{"message":"Field 'exitCode' doesn't exist on type 'CiJob'"}]}
//
// This is what gives the test teeth: it fails loudly if the Doer stops
// stripping a gated field, instead of silently sending a query that happens to
// work against a modern schema.
//
// Only one field is checked per server, deliberately. A blanket search for
// every gated name would produce false rejections, because several gated names
// also occur ungated on other types - getProjectsMergeRequests selects
// MergeRequest.updatedAt, which exists in every supported version even though
// Project.updatedAt is gated.
func simulatedGitLab(t *testing.T, serverVersion string, field gatedField) *httptest.Server {
	t.Helper()

	gate, err := doer.NewGate(serverVersion)
	if err != nil {
		t.Fatalf("NewGate(%q): %v", serverVersion, err)
	}
	known, err := gate.Allows(field.since, "")
	if err != nil {
		t.Fatalf("Allows(%q): %v", field.since, err)
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := receivedQuery(t, r)
		w.Header().Set("Content-Type", "application/json")

		// A version directive must never reach the server: it is declared only
		// in spec/directives.graphql, a local codegen-time input, and the real
		// GitLab schema does not define it.
		for _, d := range []string{"@minVersion", "@maxVersion"} {
			if strings.Contains(query, d) {
				_, _ = fmt.Fprintf(w, `{"errors":[{"message":"Directive %s is not defined"}]}`, d)
				return
			}
		}

		if !known && strings.Contains(query, field.name) {
			_, _ = fmt.Fprintf(w,
				`{"errors":[{"message":"Field '%s' doesn't exist on type '%s'"}]}`,
				field.name, field.owner)
			return
		}

		// genqlient unmarshals an empty data object cleanly, leaving every
		// field at its zero value. Reaching here means the server accepted the
		// query, which is the assertion.
		_, _ = w.Write([]byte(`{"data":{}}`))
	}))
}

// receivedQuery extracts the GraphQL query from either a GET URL parameter
// (the client is built with NewClientUsingGet) or a JSON POST body.
func receivedQuery(t *testing.T, r *http.Request) string {
	t.Helper()

	if q := r.URL.Query().Get("query"); q != "" {
		return q
	}
	if r.Body == nil {
		return ""
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("read request body: %v", err)
		return ""
	}
	var p struct{ Query string }
	if err := json.Unmarshal(body, &p); err != nil {
		t.Errorf("unmarshal request body %q: %v", body, err)
		return ""
	}
	return p.Query
}

// TestVersionGating_OperationsAcrossVersions drives the real generated
// operations through the real genqlient client and the version-gating Doer
// against servers simulating a spread of GitLab versions. It covers the
// end-to-end promise: one operation per concern, many server versions, no
// validation errors.
func TestVersionGating_OperationsAcrossVersions(t *testing.T) {
	ptr := func(b bool) *bool { return &b }

	operations := []struct {
		name  string
		field gatedField
		run   func(ctx context.Context, c *Client) error
	}{
		{
			name:  "getProjects",
			field: projectUpdatedAt,
			run: func(ctx context.Context, c *Client) error {
				_, err := getProjects(ctx, c.client, nil, nil, nil, nil)
				return err
			},
		},
		{
			name:  "getProjectsMergeRequests",
			field: mrClosedAt,
			run: func(ctx context.Context, c *Client) error {
				_, err := getProjectsMergeRequests(ctx, c.client, nil, nil, nil, nil,
					ptr(true), ptr(true), ptr(false), ptr(false))
				return err
			},
		},
		{
			name:  "getProjectsPipelinesJobs",
			field: jobExitCode,
			run: func(ctx context.Context, c *Client) error {
				_, err := getProjectsPipelinesJobs(ctx, c.client, nil, nil, nil, nil, nil, true, true)
				return err
			},
		},
	}

	// Spans every gate boundary: below all of them, exactly on each, and above
	// all of them. 17.4.0-ee and 17.5.2-ee also pin the suffix-tolerant
	// parsing, and the "exactly on" cases pin the inclusive lower bound.
	versions := []string{
		"16.10.0",
		"17.0.0",
		"17.4.0-ee",
		"17.5.2-ee",
		"17.7.0",
		"18.1.0",
	}

	for _, op := range operations {
		for _, v := range versions {
			t.Run(op.name+"/"+v, func(t *testing.T) {
				srv := simulatedGitLab(t, v, op.field)
				defer srv.Close()

				c := NewClient(srv.URL, "token")
				if err := c.SetVersion(version.GitLabVersion{Raw: v}); err != nil {
					t.Fatalf("SetVersion(%q): %v", v, err)
				}

				if err := op.run(context.Background(), c); err != nil {
					t.Errorf("%s against simulated GitLab %s: %v", op.name, v, err)
				}
			})
		}
	}
}

// TestVersionGating_FailsClosedWithoutVersion asserts the client refuses to
// talk to the server until a version is installed. Sending an ungated query
// would risk referencing fields the server does not have.
func TestVersionGating_FailsClosedWithoutVersion(t *testing.T) {
	var reached bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached = true
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{}}`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "token")

	if _, err := getProjects(context.Background(), c.client, nil, nil, nil, nil); err == nil {
		t.Error("expected getProjects to fail before SetVersion, got nil error")
	}
	if reached {
		t.Error("no request may reach the server before a version gate is installed")
	}
}
