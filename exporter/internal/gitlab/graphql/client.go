package graphql

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/graphql/doer"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/version"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/httpclient"
)

const apiPath string = "api/graphql/"

type Client struct {
	client  graphql.Client
	doer    *doer.VersionGatedDoer
	Version version.GitLabVersion
}

func NewClient(url string, token string) *Client {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	if !strings.HasSuffix(url, apiPath) {
		url += apiPath
	}

	inner := newAuthedClient(token)
	doer := doer.NewVersionGatedDoer(inner)

	return &Client{
		client: graphql.NewClientUsingGet(url, doer),
		doer:   doer,
	}
}

// SetVersion installs the given GitLab version on the GraphQL client. The
// version is used to evaluate @minVersion / @maxVersion directives on each
// outgoing query.
//
// Returns an error if v cannot be expressed as a valid semver string.
func (c *Client) SetVersion(v version.GitLabVersion) error {
	raw := v.Raw
	if raw == "" {
		raw = fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	}
	gate, err := doer.NewGate(raw)
	if err != nil {
		return fmt.Errorf("graphql: set version %q: %w", raw, err)
	}
	c.doer.SetGate(gate)
	c.Version = v
	return nil
}

type authedClient struct {
	*http.Client
	token string
}

func newAuthedClient(token string) *authedClient {
	return &authedClient{
		Client: httpclient.New().StandardClient(),
		token:  token,
	}
}

func (c *authedClient) Do(req *http.Request) (*http.Response, error) {
	if values := req.Header.Values("Authorization"); len(values) == 0 {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return c.Client.Do(req)
}
