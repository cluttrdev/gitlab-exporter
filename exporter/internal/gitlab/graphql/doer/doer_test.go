package doer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// captureDoer records every request body it sees, then returns a static 200.
type captureDoer struct {
	mu       sync.Mutex
	calls    int
	bodies   []string
	requests []*http.Request
}

func (c *captureDoer) Do(req *http.Request) (*http.Response, error) {
	var body []byte
	if req.Body != nil {
		body, _ = io.ReadAll(req.Body)
		_ = req.Body.Close()
	}
	c.mu.Lock()
	c.calls++
	c.bodies = append(c.bodies, string(body))
	c.requests = append(c.requests, req)
	c.mu.Unlock()
	return &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"data":{}}`))),
		Header:     make(http.Header),
	}, nil
}

func newRequest(t *testing.T, body string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, "http://example/graphql", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestDoer_GateUnset(t *testing.T) {
	cap := &captureDoer{}
	d := NewVersionGatedDoer(cap)

	body := `{"query":"query Q { x }","variables":{},"operationName":"Q"}`
	_, err := d.Do(newRequest(t, body))
	if !errors.Is(err, ErrGateUnset) {
		t.Errorf("expected ErrGateUnset, got %v", err)
	}
	if cap.calls != 0 {
		t.Errorf("inner Doer should not be called when gate is unset, got %d calls", cap.calls)
	}
}

func TestDoer_PrunesFieldsAndVariables(t *testing.T) {
	cap := &captureDoer{}
	d := NewVersionGatedDoer(cap)
	g, _ := NewGate("16.0")
	d.SetGate(g)

	body := `{
		"query": "query Q($cursor: String, $alive: Boolean!) { project { gated(after: $cursor) @minVersion(v: \"17.0\") alwaysHere(flag: $alive) } }",
		"variables": {"cursor": "abc", "alive": true},
		"operationName": "Q"
	}`
	if _, err := d.Do(newRequest(t, body)); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if cap.calls != 1 {
		t.Fatalf("expected 1 inner call, got %d", cap.calls)
	}

	var got payload
	if err := json.Unmarshal([]byte(cap.bodies[0]), &got); err != nil {
		t.Fatalf("unmarshal forwarded body: %v", err)
	}

	if strings.Contains(got.Query, "gated") {
		t.Errorf("expected gated field to be stripped, got:\n%s", got.Query)
	}
	if !strings.Contains(got.Query, "alwaysHere") {
		t.Errorf("expected alwaysHere to remain, got:\n%s", got.Query)
	}
	if _, ok := got.Variables["cursor"]; ok {
		t.Errorf("expected $cursor to be pruned from variables, got %v", got.Variables)
	}
	if _, ok := got.Variables["alive"]; !ok {
		t.Errorf("expected $alive to remain, got %v", got.Variables)
	}

	// Both the ContentLength field and the Content-Length header must match
	// the rewritten body length; a stale header would make the server read a
	// truncated or over-long body.
	wantLen := len(cap.bodies[0])
	if got := cap.requests[0].ContentLength; got != int64(wantLen) {
		t.Errorf("ContentLength = %d, want %d", got, wantLen)
	}
	if got := cap.requests[0].Header.Get("Content-Length"); got != strconv.Itoa(wantLen) {
		t.Errorf("Content-Length header = %q, want %q", got, strconv.Itoa(wantLen))
	}
}

func TestDoer_CachesRewrite(t *testing.T) {
	cap := &captureDoer{}
	d := NewVersionGatedDoer(cap)
	g, _ := NewGate("17.0")
	d.SetGate(g)

	body := `{"query":"query Q { updatedAt @minVersion(v: \"17.0\") }","variables":{}}`
	for i := 0; i < 3; i++ {
		if _, err := d.Do(newRequest(t, body)); err != nil {
			t.Fatalf("Do %d: %v", i, err)
		}
	}

	// All three calls must reach the inner doer.
	if cap.calls != 3 {
		t.Errorf("expected 3 inner calls, got %d", cap.calls)
	}

	// The cache should hold exactly one entry. Count via sync.Map iteration.
	count := 0
	d.cache.Range(func(_, _ any) bool {
		count++
		return true
	})
	if count != 1 {
		t.Errorf("expected 1 cache entry, got %d", count)
	}
}

func TestDoer_SetGateClearsCache(t *testing.T) {
	cap := &captureDoer{}
	d := NewVersionGatedDoer(cap)
	g17, _ := NewGate("17.0")
	d.SetGate(g17)

	body := `{"query":"query Q { x @minVersion(v: \"17.0\") }","variables":{}}`
	if _, err := d.Do(newRequest(t, body)); err != nil {
		t.Fatal(err)
	}

	g16, _ := NewGate("16.0")
	d.SetGate(g16)

	count := 0
	d.cache.Range(func(_, _ any) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("SetGate should clear cache, got %d entries", count)
	}
}

// TestDoer_SetGateNilClearsGate pins the documented ability to un-install a
// gate, which returns the Doer to its fail-closed initial state rather than
// silently forwarding unrewritten queries.
func TestDoer_SetGateNilClearsGate(t *testing.T) {
	cap := &captureDoer{}
	d := NewVersionGatedDoer(cap)

	g17, _ := NewGate("17.0")
	d.SetGate(g17)
	if d.Gate() == nil {
		t.Fatal("Gate() should report the installed gate")
	}

	body := `{"query":"query Q { x }","variables":{}}`
	if _, err := d.Do(newRequest(t, body)); err != nil {
		t.Fatalf("request with gate installed: %v", err)
	}

	d.SetGate(nil)
	if d.Gate() != nil {
		t.Error("Gate() should report nil after SetGate(nil)")
	}

	before := cap.calls
	if _, err := d.Do(newRequest(t, body)); !errors.Is(err, ErrGateUnset) {
		t.Errorf("after SetGate(nil), error = %v, want ErrGateUnset", err)
	}
	if cap.calls != before {
		t.Errorf("request must not be forwarded without a gate: calls %d -> %d", before, cap.calls)
	}
}

// TestDoer_PassThroughPaths covers the requests the Doer must forward
// untouched rather than reject: no body, empty body, and a JSON body carrying
// no query. A regression here would break non-GraphQL traffic on the same
// client.
func TestDoer_PassThroughPaths(t *testing.T) {
	tests := []struct {
		name string
		body string
		nil  bool
	}{
		{name: "nil body", nil: true},
		{name: "empty body", body: ""},
		{name: "json without query", body: `{"variables":{"a":1}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cap := &captureDoer{}
			d := NewVersionGatedDoer(cap)
			g, _ := NewGate("17.0")
			d.SetGate(g)

			req := newRequest(t, tt.body)
			if tt.nil {
				req.Body = nil
			}
			if _, err := d.Do(req); err != nil {
				t.Fatalf("expected pass-through, got error: %v", err)
			}
			if cap.calls != 1 {
				t.Errorf("expected request forwarded once, got %d calls", cap.calls)
			}
		})
	}
}

// TestDoer_MalformedJSONBody asserts a body that is not valid JSON is
// rejected rather than forwarded, since the Doer cannot know whether it
// contains a gated query.
func TestDoer_MalformedJSONBody(t *testing.T) {
	cap := &captureDoer{}
	d := NewVersionGatedDoer(cap)
	g, _ := NewGate("17.0")
	d.SetGate(g)

	if _, err := d.Do(newRequest(t, `{"query": bogus`)); err == nil {
		t.Error("expected error for malformed JSON body, got nil")
	}
	if cap.calls != 0 {
		t.Errorf("malformed body must not be forwarded, got %d calls", cap.calls)
	}
}

func TestDoer_RewriteErrorAbortsRequest(t *testing.T) {
	cap := &captureDoer{}
	d := NewVersionGatedDoer(cap)
	g, _ := NewGate("17.0")
	d.SetGate(g)

	// Malformed version in marker - rewriter returns error.
	body := `{"query":"query Q { x @minVersion(v: \"not-a-version\") }","variables":{}}`
	_, err := d.Do(newRequest(t, body))
	if err == nil {
		t.Fatal("expected rewriter error")
	}
	if cap.calls != 0 {
		t.Errorf("inner Doer should not be called on rewrite error, got %d calls", cap.calls)
	}
}

func TestDoer_GetBodySupportsRetry(t *testing.T) {
	cap := &captureDoer{}
	d := NewVersionGatedDoer(cap)
	g, _ := NewGate("17.0")
	d.SetGate(g)

	body := `{"query":"query Q { name }","variables":{}}`
	req := newRequest(t, body)
	if _, err := d.Do(req); err != nil {
		t.Fatal(err)
	}

	if cap.requests[0].GetBody == nil {
		t.Fatal("expected GetBody to be set on forwarded request")
	}
	rc, err := cap.requests[0].GetBody()
	if err != nil {
		t.Fatalf("GetBody: %v", err)
	}
	defer rc.Close()
	replay, _ := io.ReadAll(rc)
	if string(replay) != cap.bodies[0] {
		t.Errorf("GetBody replay differs from original\nreplay=%q\noriginal=%q", replay, cap.bodies[0])
	}
}

func TestDoer_RewritesGETRequest(t *testing.T) {
	cap := &captureDoer{}
	d := NewVersionGatedDoer(cap)
	g, _ := NewGate("16.0")
	d.SetGate(g)

	rawQuery := `query Q($cursor: String) { project { gated(after: $cursor) @minVersion(v: "17.0") name } }`
	vars := `{"cursor":"abc"}`

	u, _ := http.NewRequest(http.MethodGet, "http://example/graphql", nil)
	q := u.URL.Query()
	q.Set("query", rawQuery)
	q.Set("variables", vars)
	q.Set("operationName", "Q")
	u.URL.RawQuery = q.Encode()

	if _, err := d.Do(u); err != nil {
		t.Fatalf("Do: %v", err)
	}
	if cap.calls != 1 {
		t.Fatalf("expected 1 inner call, got %d", cap.calls)
	}
	forwarded := cap.requests[0].URL.Query()

	if strings.Contains(forwarded.Get("query"), "gated") {
		t.Errorf("expected gated stripped, got query:\n%s", forwarded.Get("query"))
	}
	if !strings.Contains(forwarded.Get("query"), "name") {
		t.Errorf("expected name preserved, got query:\n%s", forwarded.Get("query"))
	}
	var fwdVars map[string]interface{}
	if err := json.Unmarshal([]byte(forwarded.Get("variables")), &fwdVars); err != nil {
		t.Fatalf("unmarshal variables: %v", err)
	}
	if _, ok := fwdVars["cursor"]; ok {
		t.Errorf("expected $cursor pruned from GET variables, got %v", fwdVars)
	}
	if forwarded.Get("operationName") != "Q" {
		t.Errorf("operationName lost, got %q", forwarded.Get("operationName"))
	}
}

// TestDoer_Phase5_UpdatedAtGating pins the Phase 5 acceptance criterion: a
// query carrying @minVersion(v: "17.0") on the projects' updatedAt field is
// stripped for a pre-17 server and kept for a 17+ server. This mirrors the
// real genqlient-emitted query for getProjects after the Phase 5 fragment
// annotation lands.
func TestDoer_Phase5_UpdatedAtGating(t *testing.T) {
	const query = `query getProjects($ids: [ID!]) {
  projects(ids: $ids) {
    nodes {
      id
      updatedAt @minVersion(v: "17.0")
      lastActivityAt
    }
  }
}`

	cases := []struct {
		name          string
		server        string
		wantUpdatedAt bool
	}{
		{"pre-17 strips updatedAt", "16.10.0", false},
		{"17.0 keeps updatedAt", "17.0.0", true},
		{"17.5 keeps updatedAt", "17.5.2-ee", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cap := &captureDoer{}
			d := NewVersionGatedDoer(cap)
			g, err := NewGate(c.server)
			if err != nil {
				t.Fatal(err)
			}
			d.SetGate(g)

			u, _ := http.NewRequest(http.MethodGet, "http://example/graphql", nil)
			q := u.URL.Query()
			q.Set("query", query)
			q.Set("variables", `{"ids":["gid://gitlab/Project/1"]}`)
			q.Set("operationName", "getProjects")
			u.URL.RawQuery = q.Encode()

			if _, err := d.Do(u); err != nil {
				t.Fatal(err)
			}

			fwdQuery := cap.requests[0].URL.Query().Get("query")
			hasUpdatedAt := strings.Contains(fwdQuery, "updatedAt")
			hasDirective := strings.Contains(fwdQuery, "@minVersion")

			if hasUpdatedAt != c.wantUpdatedAt {
				t.Errorf("updatedAt present=%v want=%v\nforwarded query:\n%s",
					hasUpdatedAt, c.wantUpdatedAt, fwdQuery)
			}
			if hasDirective {
				t.Errorf("@minVersion should never reach the wire\nforwarded query:\n%s", fwdQuery)
			}
		})
	}
}

// TestDoer_AgainstHTTPServer is an end-to-end style check: the Doer is
// composed with a real *http.Client (not a mock) and we point it at an
// httptest server. This exercises the full body-buffer + Content-Length
// handshake the way genqlient will actually drive it.
func TestDoer_AgainstHTTPServer(t *testing.T) {
	var received string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		received = string(b)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{}}`))
	}))
	defer srv.Close()

	httpDoer := srv.Client()
	d := NewVersionGatedDoer(httpDoer)
	g, _ := NewGate("16.0")
	d.SetGate(g)

	body := `{"query":"query Q { kept dropped @minVersion(v: \"17.0\") }","variables":{}}`
	req, err := http.NewRequest(http.MethodPost, srv.URL, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if received == "" {
		t.Fatal("server received no body")
	}
	var p payload
	if err := json.Unmarshal([]byte(received), &p); err != nil {
		t.Fatalf("server body unparseable: %v", err)
	}
	if strings.Contains(p.Query, "dropped") {
		t.Errorf("expected 'dropped' field stripped, server saw:\n%s", p.Query)
	}
	if !strings.Contains(p.Query, "kept") {
		t.Errorf("expected 'kept' field present, server saw:\n%s", p.Query)
	}
}
