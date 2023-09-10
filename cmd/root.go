package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/controller"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/internal/ffyaml"
)

const (
	exeName      string = "gitlab-clickhouse-exporter"
	envVarPrefix string = "GLCHE"
)

var (
	rootCmdOptions = []ff.Option{
		ff.WithEnvVarPrefix(envVarPrefix),
		ff.WithConfigFileFlag("config"),
		ff.WithConfigFileParser(ffyaml.Parser),
		ff.WithIgnoreUndefined(true),
	}
)

type RootConfig struct {
	Config config.Config
}

func NewRootCmd() (*ffcli.Command, *RootConfig) {
	var config RootConfig

	fs := flag.NewFlagSet(exeName, flag.ContinueOnError)
	config.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       exeName,
		ShortUsage: fmt.Sprintf("%s <subcommand> [flags] [<args>...]", exeName),
		UsageFunc:  usageFunc,
		FlagSet:    fs,
		Options:    rootCmdOptions,
		Exec:       config.Exec,
	}, &config
}

func (c *RootConfig) RegisterFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.Config.GitLab.Api.URL, "gitlab-api-url", config.DefaultGitLabApiUrl, fmt.Sprintf("The GitLab API URL (default: '%s').", config.DefaultGitLabApiUrl))
	fs.StringVar(&c.Config.GitLab.Api.Token, "gitlab-api-token", config.DefaultGitLabApiToken, fmt.Sprintf("The GitLab API Token (default: '%s').", config.DefaultGitLabApiToken))
	fs.Float64Var(&c.Config.GitLab.Client.Rate.Limit, "gitlab-client-rate-limit", config.DefaultGitLabClientRateLimit, fmt.Sprintf("Rate limit for the GitLab client (default: '%.1f').", config.DefaultGitLabClientRateLimit))

	fs.StringVar(&c.Config.ClickHouse.Host, "clickhouse-host", config.DefaultClickHouseHost, fmt.Sprintf("The ClickHouse server name (default: '%s').", config.DefaultClickHouseHost))
	fs.StringVar(&c.Config.ClickHouse.Port, "clickhouse-port", config.DefaultClickHousePort, fmt.Sprintf("The ClickHouse port to connect to (default: '%s')", config.DefaultClickHousePort))
	fs.StringVar(&c.Config.ClickHouse.Database, "clickhouse-database", config.DefaultClickHouseDatabase, fmt.Sprintf("Select the current default ClickHouse database (default: '%s').", config.DefaultClickHouseDatabase))
	fs.StringVar(&c.Config.ClickHouse.User, "clickhouse-user", config.DefaultClickHouseUser, fmt.Sprintf("The ClickHouse username to connect with (default: '%s').", config.DefaultClickHouseUser))
	fs.StringVar(&c.Config.ClickHouse.Password, "clickhouse-password", config.DefaultClickHousePassword, fmt.Sprintf("The ClickHouse password (default: '%s').", config.DefaultClickHousePassword))

	fs.String("config", "", "A configuration file.")
}

func (c *RootConfig) Exec(context.Context, []string) error {
	return flag.ErrHelp
}

func (c *RootConfig) newController() (*controller.Controller, error) {
	ctl, err := controller.NewController(c.Config)
	if err != nil {
		return nil, fmt.Errorf("error constructing controller: %w", err)
	}
	return &ctl, nil
}
