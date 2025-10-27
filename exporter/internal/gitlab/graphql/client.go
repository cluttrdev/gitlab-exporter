package graphql

import (
	"net/http"
	"strings"

	"github.com/Khan/genqlient/graphql"

	"go.cluttr.dev/gitlab-exporter/internal/httpclient"
)

const apiPath string = "api/graphql/"

type Client struct {
	client graphql.Client
}

func NewClient(url string, token string) *Client {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	if !strings.HasSuffix(url, apiPath) {
		url += apiPath
	}

	doer := newAuthedClient(token)

	return &Client{
		client: graphql.NewClientUsingGet(url, doer),
	}
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
