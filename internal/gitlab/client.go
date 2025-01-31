package gitlab

import (
	"context"
	"fmt"
	"net/http"

	gitlab "gitlab.com/gitlab-org/api/client-go"
	"golang.org/x/sync/semaphore"

	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/graphql"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/rest"
)

type Client struct {
	URL string

	Rest    *rest.Client
	GraphQL *graphql.Client
	HTTP    *HTTPClient

	sem *semaphore.Weighted
}

type ClientConfig struct {
	URL   string
	Token string

	Auth AuthConfig

	RateLimit float64

	MaxWorkers int
}

type AuthType int

const (
	_ AuthType = iota
	SessionAuth
)

type AuthConfig struct {
	AuthType AuthType
	Basic    BasicAuthConfig
}

type BasicAuthConfig struct {
	Username string
	Password string
}

func NewGitLabClient(cfg ClientConfig) (*Client, error) {
	restClient, err := rest.NewClient(cfg.URL, cfg.Token, cfg.RateLimit)
	if err != nil {
		return nil, err
	}

	graphqlClient := graphql.NewClient(cfg.URL, cfg.Token)

	httpClient, err := NewHTTPClient(cfg.URL, cfg.Auth)
	if err != nil {
		return nil, err
	}

	var n int64 = 42
	if cfg.MaxWorkers > 0 {
		n = int64(cfg.MaxWorkers)
	}

	return &Client{
		URL: cfg.URL,

		Rest:    restClient,
		GraphQL: graphqlClient,
		HTTP:    httpClient,

		sem: semaphore.NewWeighted(n),
	}, nil
}

func (c *Client) Acquire(ctx context.Context, n int64) error {
	return c.sem.Acquire(ctx, n)
}

func (c *Client) Release(n int64) {
	c.sem.Release(n)
}

func (c *Client) CheckReadiness(ctx context.Context) error {
	const readinessEndpoint string = "version"

	req, err := c.Rest.Client().NewRequest(
		http.MethodGet,
		readinessEndpoint,
		nil,
		[]gitlab.RequestOptionFunc{gitlab.WithContext(ctx)},
	)
	if err != nil {
		return err
	}

	res, err := c.Rest.Client().Do(req, nil)
	if err != nil {
		return err
	}

	if res == nil {
		return fmt.Errorf("http error: empty response")
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http error: %d", res.StatusCode)
	}

	return nil
}
