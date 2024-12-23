package gitlab

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/oauth2"
	"github.com/cluttrdev/gitlab-exporter/internal/httpclient"
)

type HTTPClient struct {
	client *http.Client

	url string
}

func NewHTTPClient(url string, config *OAuthConfig) (*HTTPClient, error) {
	client := httpclient.New().StandardClient()

	if config != nil {
		var (
			token *oauth2.Token
			err   error
		)

		cfg := oauth2.Configure(url, config.ClientId, config.ClientSecret, config.Scopes)

		ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)
		switch flowType := config.FlowType; flowType {
		case "authorization_code":
			token, err = oauth2.StartAuthorizationCodeFlow(ctx, cfg)
		case "password":
			token, err = oauth2.StartPasswordFlow(ctx, cfg, config.Username, config.Password)
		default:
			err = fmt.Errorf("unsupported flow: %q", flowType)
		}
		if err != nil {
			return nil, err
		}

		client = cfg.Client(ctx, token.Token)
	}

	return &HTTPClient{
		client: client,
		url:    url,
	}, nil
}

func (c *HTTPClient) ChechAuth() error {
	_, err := c.GetPath("oauth/token/info")
	if err != nil {
		return err
	}

	return nil
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *HTTPClient) GetPath(path string) (*http.Response, error) {
	p, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse url path: %w", err)
	}
	if p.Host != "" {
		return nil, fmt.Errorf("non-empty host: %v", p.Host)
	}

	u, err := url.Parse(c.url)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}

	u.Path = p.Path
	u.RawQuery = p.RawQuery

	return c.client.Get(u.String())
}

func (c *HTTPClient) GetProjectJobArtifactsFile(ctx context.Context, projectPath string, jobId int64, fileType string) (*bytes.Reader, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}

	path, err := url.JoinPath("", projectPath, "-", "jobs", strconv.FormatInt(jobId, 10), "artifacts", "download")
	if err != nil {
		return nil, fmt.Errorf("create url path: %w", err)
	}
	u.Path = path

	query := url.Values{
		"file_type": []string{fileType},
	}
	u.RawQuery = query.Encode()

	resp, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	defer io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%d - %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)

	return bytes.NewReader(buf.Bytes()), err
}
