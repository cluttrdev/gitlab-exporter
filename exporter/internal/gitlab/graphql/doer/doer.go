package doer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
)

// ErrGateUnset is returned by VersionGatedDoer when a request arrives before
// any Gate has been installed via SetGate.
var ErrGateUnset = errors.New("versiongate: no version gate set")

// HTTPDoer is the minimal interface VersionGatedDoer wraps. It is satisfied
// by *http.Client. We use the interface form so callers can inject mocks
// during tests.
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// VersionGatedDoer wraps an inner HTTPDoer and rewrites the outgoing GraphQL
// request to remove selections gated by @minVersion / @maxVersion directives
// that disallow the configured server version. It handles both POST requests
// (with a JSON body) and GET requests (with query/variables/operationName
// URL parameters).
//
// Construct one with NewVersionGatedDoer. The gate may be installed later via
// SetGate; requests before SetGate is called will fail with ErrGateUnset.
//
// The rewritten query strings are cached keyed by raw query string. The set
// of distinct query strings is bounded by the number of genqlient operations,
// so the cache is naturally small and stable. Resetting the gate (via a
// follow-up SetGate call) clears the cache.
type VersionGatedDoer struct {
	inner HTTPDoer

	gate  atomic.Pointer[Gate]
	cache sync.Map // map[string]*rewriteCacheEntry, keyed by raw query string
}

type rewriteCacheEntry struct {
	rewritten  string
	prunedVars map[string]struct{}
}

// NewVersionGatedDoer returns a Doer that delegates to inner. The Doer is
// returned with no gate installed; callers must call SetGate before sending
// requests.
func NewVersionGatedDoer(inner HTTPDoer) *VersionGatedDoer {
	return &VersionGatedDoer{inner: inner}
}

// SetGate installs or replaces the gate used to evaluate version directives.
// Any previously cached rewrites are discarded because the version they were
// computed against is no longer relevant.
//
// It is safe to call SetGate concurrently with in-flight requests; the
// atomic swap means in-flight requests use the gate that was current when
// they read it.
//
// Passing a nil gate clears the gate; subsequent requests will fail with
// ErrGateUnset until a non-nil gate is installed again.
func (d *VersionGatedDoer) SetGate(g *Gate) {
	d.gate.Store(g)
	d.cache.Range(func(k, _ any) bool {
		d.cache.Delete(k)
		return true
	})
}

// Gate returns the currently-installed gate, or nil if none is set.
func (d *VersionGatedDoer) Gate() *Gate {
	return d.gate.Load()
}

// payload is the standard GraphQL POST body shape.
type payload struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables,omitempty"`
	OperationName string                 `json:"operationName,omitempty"`
}

// Do implements genqlient's graphql.Doer interface. It dispatches on request
// method:
//
//   - GET requests: rewrite the `query` URL parameter and delete pruned
//     variables from the `variables` URL parameter.
//   - POST requests with a JSON body: rewrite the body's `query` field
//     and prune `variables` likewise.
//
// On any rewriter error the request is not forwarded; the error is returned
// unchanged.
func (d *VersionGatedDoer) Do(req *http.Request) (*http.Response, error) {
	gate := d.gate.Load()
	if gate == nil {
		return nil, ErrGateUnset
	}

	if req.Method == http.MethodGet {
		return d.doGet(req, gate)
	}
	return d.doPost(req, gate)
}

// doGet rewrites a GET-style GraphQL request whose payload is carried in URL
// query parameters (genqlient.NewClientUsingGet emits these).
func (d *VersionGatedDoer) doGet(req *http.Request, gate *Gate) (*http.Response, error) {
	if req.URL == nil {
		return d.inner.Do(req)
	}
	q := req.URL.Query()
	query := q.Get("query")
	if query == "" {
		return d.inner.Do(req)
	}

	entry, err := d.rewrite(query, gate)
	if err != nil {
		return nil, err
	}

	q.Set("query", entry.rewritten)

	if len(entry.prunedVars) > 0 {
		varsRaw := q.Get("variables")
		if varsRaw != "" {
			var vars map[string]interface{}
			if err := json.Unmarshal([]byte(varsRaw), &vars); err != nil {
				return nil, fmt.Errorf("versiongate: unmarshal GET variables: %w", err)
			}
			for name := range entry.prunedVars {
				delete(vars, name)
			}
			newVarsRaw, err := json.Marshal(vars)
			if err != nil {
				return nil, fmt.Errorf("versiongate: marshal pruned GET variables: %w", err)
			}
			q.Set("variables", string(newVarsRaw))
		}
	}

	// Build a shallow copy of the URL so we don't mutate the caller's
	// reference (defensive; genqlient's createGetRequest builds a fresh URL
	// per request, but Doers are conventionally non-mutating).
	newURL := *req.URL
	newURL.RawQuery = q.Encode()
	req.URL = &newURL
	return d.inner.Do(req)
}

// doPost rewrites a POST-style GraphQL request whose payload is a JSON body.
func (d *VersionGatedDoer) doPost(req *http.Request, gate *Gate) (*http.Response, error) {
	if req.Body == nil {
		return d.inner.Do(req)
	}

	raw, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("versiongate: read request body: %w", err)
	}
	_ = req.Body.Close()

	if len(raw) == 0 {
		return d.forwardBody(req, raw)
	}

	var p payload
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, fmt.Errorf("versiongate: unmarshal request body: %w", err)
	}

	if p.Query == "" {
		return d.forwardBody(req, raw)
	}

	entry, err := d.rewrite(p.Query, gate)
	if err != nil {
		return nil, err
	}

	if len(entry.prunedVars) > 0 && len(p.Variables) > 0 {
		for name := range entry.prunedVars {
			delete(p.Variables, name)
		}
	}
	p.Query = entry.rewritten

	newBody, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("versiongate: marshal rewritten body: %w", err)
	}

	return d.forwardBody(req, newBody)
}

// rewrite looks up or computes the cached rewrite for a given raw query
// string. Cache entries are immutable; concurrent callers either find an
// existing entry or both compute the same result and the last writer wins
// (which is fine because the result is deterministic).
func (d *VersionGatedDoer) rewrite(raw string, gate *Gate) (*rewriteCacheEntry, error) {
	if cached, ok := d.cache.Load(raw); ok {
		return cached.(*rewriteCacheEntry), nil
	}
	rewritten, pruned, err := stripVersionedNodes(raw, gate)
	if err != nil {
		return nil, fmt.Errorf("versiongate: rewrite query: %w", err)
	}
	entry := &rewriteCacheEntry{rewritten: rewritten, prunedVars: pruned}
	d.cache.Store(raw, entry)
	return entry, nil
}

// forwardBody installs body on req and delegates to the inner Doer. It sets
// up GetBody so retries and redirects can replay the body.
func (d *VersionGatedDoer) forwardBody(req *http.Request, body []byte) (*http.Response, error) {
	req.Body = io.NopCloser(bytes.NewReader(body))
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	req.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(body)), nil
	}
	return d.inner.Do(req)
}
