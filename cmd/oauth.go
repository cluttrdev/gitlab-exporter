package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"

	"github.com/cluttrdev/cli"
	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/oauth2"
)

type OAuthConfig struct {
	RootConfig
}

func NewOAuthCmd(out io.Writer) *cli.Command {
	cfg := OAuthConfig{
		RootConfig: RootConfig{
			out:   out,
			flags: flag.NewFlagSet(fmt.Sprintf("%s oauth", exeName), flag.ExitOnError),
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "oauth",
		ShortUsage: fmt.Sprintf("%s oauth <subcommand> [option]...", exeName),
		ShortHelp:  "Use GitLab as an OAuth 2.0 identity provider",
		Flags:      cfg.flags,
		Subcommands: []*cli.Command{
			NewOAuthRequestCmd(out),
			NewOAuthRefreshCmd(out),
		},
		Exec: cfg.Exec,
	}
}

func (c *OAuthConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootConfig.RegisterFlags(fs)
}

func (c *OAuthConfig) Exec(ctx context.Context, _ []string) error {
	return flag.ErrHelp
}

// Request

type OAuthRequestConfig struct {
	OAuthConfig

	scope string
}

func NewOAuthRequestCmd(out io.Writer) *cli.Command {
	cfg := OAuthRequestConfig{
		OAuthConfig: OAuthConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet(fmt.Sprintf("%s oauth request", exeName), flag.ExitOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "request",
		ShortUsage: fmt.Sprintf("%s oauth request [option]...", exeName),
		ShortHelp:  "Request token",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *OAuthRequestConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootConfig.RegisterFlags(fs)

	fs.StringVar(&c.scope, "scope", "", "The token scope to request.")
}

func (c *OAuthRequestConfig) Exec(ctx context.Context, args []string) error {
	var err error

	cfg := config.Default()
	if err = loadConfig(c.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	secrets := cfg.GitLab.OAuth.GitLabOAuthSecrets
	if cfg.GitLab.OAuth.SecretsFile != "" {
		secrets, err = config.LoadOAuthSecretsFile(cfg.GitLab.OAuth.SecretsFile)
	}

	scopes := []string{c.scope}

	config := oauth2.Configure(cfg.GitLab.Url, secrets.ClientId, secrets.ClientSecret, scopes)

	var token *oauth2.Token
	switch cfg.GitLab.OAuth.FlowType {
	case "authorization_code":
		token, err = oauth2.StartAuthorizationCodeFlow(ctx, config)
	case "password":
		token, err = oauth2.StartPasswordFlow(ctx, config, secrets.Username, secrets.Password)
	default:
		err = fmt.Errorf("invalid flow type: %s", cfg.GitLab.OAuth.FlowType)
	}
	if err != nil {
		return err
	}

	return json.NewEncoder(c.RootConfig.out).Encode(token)
}

// Refresh

type OAuthRefreshConfig struct {
	OAuthConfig

	refreshToken string
}

func NewOAuthRefreshCmd(out io.Writer) *cli.Command {
	cfg := OAuthRefreshConfig{
		OAuthConfig: OAuthConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet(fmt.Sprintf("%s oauth refresh", exeName), flag.ExitOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "refresh",
		ShortUsage: fmt.Sprintf("%s oauth refresh [option]...", exeName),
		ShortHelp:  "Refresh a token",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *OAuthRefreshConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootConfig.RegisterFlags(fs)

	fs.StringVar(&c.refreshToken, "refresh-token", "", "The refresh token.")
}

func (c *OAuthRefreshConfig) Exec(ctx context.Context, args []string) error {
	var err error

	cfg := config.Default()
	if err = loadConfig(c.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	token, err := oauth2.RefreshToken(ctx, cfg.GitLab.Url, c.refreshToken)
	if err != nil {
		return err
	}

	return json.NewEncoder(c.RootConfig.out).Encode(token)
}
