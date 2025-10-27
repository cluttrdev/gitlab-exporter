package rest

import (
	"strings"

	"golang.org/x/time/rate"

	"gitlab.com/gitlab-org/api/client-go"
)

const apiPath string = "api/v4/"

type Client struct {
	client *gitlab.Client
}

func NewClient(url string, token string, rateLimit float64) (*Client, error) {
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	if !strings.HasSuffix(url, apiPath) {
		url += apiPath
	}

	opts := []gitlab.ClientOptionFunc{
		gitlab.WithBaseURL(url),
	}

	if rateLimit > 0 {
		limit := rate.Limit(rateLimit * 0.66)
		burst := rateLimit * 0.33
		limiter := rate.NewLimiter(limit, int(burst))

		opts = append(opts, gitlab.WithCustomLimiter(limiter))
	}

	client, err := gitlab.NewOAuthClient(token, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: client,
	}, nil
}

func (c *Client) Client() *gitlab.Client {
	return c.client
}
