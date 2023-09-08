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
	fs.StringVar(&c.Config.GitLab.Api.URL, "gitlab-api-url", config.DefaultGitLabApiUrl, "The GitLab API URL.")
	fs.StringVar(&c.Config.GitLab.Api.Token, "gitlab-api-token", config.DefaultGitLabApiToken, "The GitLab API Token.")
	fs.Int64Var(&c.Config.GitLab.Client.RequestsPerSecond, "gitlab-client-requests-per-second", config.DefaultGitLabClientRequestsPerSecond, "Rate limit for the GitLab client.")

	fs.StringVar(&c.Config.ClickHouse.Host, "clickhouse-host", config.DefaultClickHouseHost, "The ClickHouse server name (default: 'localhost').")
	fs.StringVar(&c.Config.ClickHouse.Port, "clickhouse-port", config.DefaultClickHousePort, "The ClickHouse port to connect to (default: 9000)")
	fs.StringVar(&c.Config.ClickHouse.Database, "clickhouse-database", config.DefaultClickHouseDatabase, "Select the current default ClickHouse database (default: 'default').")
	fs.StringVar(&c.Config.ClickHouse.User, "clickhouse-user", config.DefaultClickHouseUser, "The ClickHouse username to connect with (default: 'default').")
	fs.StringVar(&c.Config.ClickHouse.Password, "clickhouse-password", config.DefaultClickHousePassword, "The ClickHouse password (default: '').")

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
