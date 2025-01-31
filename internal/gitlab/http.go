package gitlab

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/oauth2"
	"github.com/cluttrdev/gitlab-exporter/internal/httpclient"
	"golang.org/x/net/html"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPClient struct {
	doer Doer
	url  string
}

func NewHTTPClient(baseUrl string, cfg AuthConfig) (*HTTPClient, error) {
	switch cfg.AuthType {
	case SessionAuth:
		return newSessionAuthedHTTPClient(baseUrl, cfg.Basic)
	case OAuth:
		return newOAuthedHTTPClient(baseUrl, cfg.OAuth)
	default:
		return &HTTPClient{
			doer: httpclient.New(),
			url:  baseUrl,
		}, nil
	}
}

func (c *HTTPClient) CheckAuthed() error {
	switch d := c.doer.(type) {
	case *sessionAuthedClient:
		return d.CheckAuthed()
	}

	return errors.ErrUnsupported
}

func (c *HTTPClient) Do(req *http.Request) (*http.Response, error) {
	return c.doer.Do(req)
}

func (c *HTTPClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(req)
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

	return c.Get(u.String())
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

	resp, err := c.Get(u.String())
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

func newOAuthedHTTPClient(baseUrl string, cfg OAuthConfig) (*HTTPClient, error) {
	var (
		token *oauth2.Token
		err   error
	)

	config := oauth2.Configure(baseUrl, cfg.ClientId, cfg.ClientSecret, cfg.Scopes)
	client := httpclient.New().StandardClient()
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, client)

	switch flowType := cfg.FlowType; flowType {
	case "authorization_code":
		token, err = oauth2.StartAuthorizationCodeFlow(ctx, config)
	case "password":
		token, err = oauth2.StartPasswordFlow(ctx, config, cfg.Username, cfg.Password)
	default:
		err = fmt.Errorf("unsupported flow: %q", flowType)
	}
	if err != nil {
		return nil, err
	}

	return &HTTPClient{
		doer: config.Client(ctx, token.Token),
		url:  baseUrl,
	}, nil
}

func newSessionAuthedHTTPClient(baseUrl string, cfg BasicAuthConfig) (*HTTPClient, error) {
	session := newSessionAuthedClient(baseUrl, cfg.Username, cfg.Password)

	// https://gitlab.com/gitlab-org/gitlab/-/issues/395038
	if err := session.signIn(); err != nil {
		return nil, fmt.Errorf("create session authed http client: %w", err)
	}

	return &HTTPClient{
		doer: session,
		url:  baseUrl,
	}, nil
}

type sessionAuthedClient struct {
	*http.Client
	mx sync.RWMutex

	url      string
	username string
	password string

	signedInAt time.Time
}

func newSessionAuthedClient(baseURL string, username string, password string) *sessionAuthedClient {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	client := httpclient.New().StandardClient()
	jar, _ := cookiejar.New(nil)
	client.Jar = jar

	return &sessionAuthedClient{
		Client:   client,
		url:      baseURL,
		username: username,
		password: password,
	}
}

func (s *sessionAuthedClient) Do(req *http.Request) (*http.Response, error) {
	if err := s.ensureSession(); err != nil {
		return nil, fmt.Errorf("ensure session: %w", err)
	}
	return s.Client.Do(req)
}

func (s *sessionAuthedClient) CheckAuthed() error {
	url := s.url + "users/" + s.username + "/exists"
	resp, err := s.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return errors.New("unauthorized")
	} else if code := resp.StatusCode; code != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d - %s", code, http.StatusText(code))
	}

	return nil
}

func (s *sessionAuthedClient) ensureSession() error {
	s.mx.Lock()
	defer s.mx.Unlock()

	if time.Now().UTC().Sub(s.signedInAt) > 1*time.Hour {
		return nil
	}

	if !s.signedInAt.IsZero() {
		if err := s.signOut(); err != nil {
			return fmt.Errorf("sign out: %w", err)
		}
		s.signedInAt = time.Time{}
	}

	if err := s.signIn(); err != nil {
		return fmt.Errorf("sign in: %w", err)
	}
	s.signedInAt = time.Now().UTC()

	return nil
}

func (s *sessionAuthedClient) signIn() error {
	signInURL := s.url + "users/sign_in?auto_sign_in=false"

	client := &http.Client{
		Jar: s.Jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// 1. Get CSRF token from login page
	req, err := http.NewRequest(
		http.MethodGet,
		signInURL,
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := client.Get(signInURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	csrfParam, csrfToken := parseCSRFToken(resp.Body)

	// 2. Sign in
	data := url.Values{
		"user[login]":    {s.username},
		"user[password]": {s.password},
		csrfParam:        {csrfToken},
	}

	req, err = http.NewRequest(
		http.MethodPost,
		signInURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// we expect to be redirected on successful sign in
	if code := resp.StatusCode; code != http.StatusFound {
		return fmt.Errorf("unexpected sign in response: %d - %s", code, http.StatusText(code))
	}

	return nil
}

func (s *sessionAuthedClient) signOut() error {
	signOutURL := s.url + "users/sign_out"

	req, err := http.NewRequest(http.MethodPost, signOutURL, nil)
	if err != nil {
		return err
	}

	_, err = s.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionAuthedClient) sessionCookie() *http.Cookie {
	s.mx.RLock()
	defer s.mx.RUnlock()

	u, _ := url.Parse(s.url)
	for _, c := range s.Jar.Cookies(u) {
		if c.Name == "_gitlab_session" {
			return c
		}
	}
	return nil
}

func parseCSRFToken(r io.Reader) (string, string) {
	var csrfParam, csrfToken string

	lexer := html.NewTokenizer(r)
	for {
		tokenType := lexer.Next()
		if tokenType == html.ErrorToken {
			break
		}

		tok := lexer.Token()
		if tokenType != html.SelfClosingTagToken || tok.Data != "meta" {
			continue
		}

		attrs := attrMap(tok.Attr)
		switch {
		case hasKeyVal(attrs, "name", "csrf-param"):
			csrfParam = attrs["content"]
		case hasKeyVal(attrs, "name", "csrf-token"):
			csrfToken = attrs["content"]
		}

		if csrfParam != "" && csrfToken != "" {
			break
		}
	}

	return csrfParam, csrfToken
}

func attrMap(attrs []html.Attribute) map[string]string {
	m := make(map[string]string, len(attrs))
	for _, a := range attrs {
		m[a.Key] = a.Val
	}
	return m
}

func hasKeyVal(m map[string]string, key string, val string) bool {
	if v, ok := m[key]; ok {
		return v == val
	}
	return false
}
