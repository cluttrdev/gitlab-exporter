package oauth2

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
)

const (
	authURLPath     = "oauth/authorize"
	tokenURLPath    = "oauth/token"
	redirectURLPath = "auth/redirect"
	redirectURL     = "http://localhost:7171/" + redirectURLPath
)

type Config struct {
	*oauth2.Config
}

type Token struct {
	*oauth2.Token
}

var (
	HTTPClient = oauth2.HTTPClient
)

func Configure(hostURL string, clientID string, clientSecret string, scopes []string) *Config {
	if !strings.HasSuffix(hostURL, "/") {
		hostURL += "/"
	}

	return &Config{
		Config: &oauth2.Config{
			ClientID: clientID,
			Endpoint: oauth2.Endpoint{
				AuthURL:  hostURL + authURLPath,
				TokenURL: hostURL + tokenURLPath,
			},
			RedirectURL: redirectURL,
			Scopes:      scopes,
		},
	}
}

func StartPasswordFlow(ctx context.Context, cfg *Config, username string, password string) (*Token, error) {
	token, err := cfg.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return &Token{Token: token}, nil
}

func StartAuthorizationCodeFlow(ctx context.Context, cfg *Config) (*Token, error) {
	state := oauth2.GenerateVerifier()
	verifier := oauth2.GenerateVerifier()

	authURL := cfg.AuthCodeURL(state, oauth2.S256ChallengeOption(verifier))

	tokenCh := handleAuthRedirect(cfg.Config, "0.0.0.0:7171", state, verifier)
	defer close(tokenCh)

	if err := browser.OpenURL(authURL); err != nil {
		slog.Error("failed opening browser", "url", authURL, "err", err)
	}
	token := <-tokenCh
	if token == nil {
		return nil, fmt.Errorf("authentication failed: no token received")
	}

	return &Token{Token: token}, nil
}

func handleAuthRedirect(config *oauth2.Config, listenAddress, originalState, codeVerifier string) chan *oauth2.Token {
	tokenCh := make(chan *oauth2.Token)

	http.HandleFunc(fmt.Sprintf("/%s", redirectURLPath), func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		state := r.FormValue("state")

		if state != originalState {
			slog.Error("Invalid state")
			tokenCh <- nil
			return
		}

		token, err := config.Exchange(r.Context(), code, oauth2.VerifierOption(codeVerifier))
		if err != nil {
			slog.Error("Error requesting access token", "err", err)
			tokenCh <- nil
			return
		}

		_, _ = w.Write([]byte("You have authenticated successfully. You can now close this browser window."))
		tokenCh <- token
	})

	go func() {
		err := http.ListenAndServe(listenAddress, nil)
		if err != nil {
			slog.Error("Error setting up server", "err", err)
			tokenCh <- nil
		}
	}()

	return tokenCh
}

func RefreshToken(ctx context.Context, hostURL string, refreshToken string) (*Token, error) {
	data := url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{refreshToken},
	}

	tokenURL, err := url.JoinPath(hostURL, tokenURLPath)
	if err != nil {
		return nil, err
	}

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		text := http.StatusText(resp.StatusCode)
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%d - %s: %s", resp.StatusCode, text, string(respBody))
	}

	var token *oauth2.Token

	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return nil, err
	}

	return &Token{Token: token}, nil
}
