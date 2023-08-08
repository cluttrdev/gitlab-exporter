package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/controller"
)

const (
    exeName string = "glche"
    envVarPrefix string = "GLCHE"
)

type RootConfig struct {
    Config config.Config
    Controller *controller.Controller
}

func NewRootCmd() (*ffcli.Command, *RootConfig) {
    var config RootConfig

    fs := flag.NewFlagSet(exeName, flag.ExitOnError)
    config.RegisterFlags(fs)

    return &ffcli.Command{
        Name: exeName,
        ShortUsage: fmt.Sprintf("%s [flags] <subcommand> [flags]", exeName),
        FlagSet: fs,
        Options: []ff.Option{ff.WithEnvVarPrefix(envVarPrefix)},
        Exec: config.Exec,
    }, &config
}

func (c *RootConfig) RegisterFlags(fs *flag.FlagSet) {
    fs.StringVar(&c.Config.GitLab.URL, "gitlab-api-url", config.DefaultGitLabApiUrl, "The GitLab API URL.")
    fs.StringVar(&c.Config.GitLab.Token, "gitlab-api-token", "", "The GitLab API Token.")

    fs.StringVar(&c.Config.ClickHouse.Host, "clickhouse-host", config.DefaultClickHouseHost, "The ClickHouse server name (default: 'localhost').")
    fs.IntVar(&c.Config.ClickHouse.Port, "clickhouse-port", config.DefaultClickHousePort, "The ClickHouse port to connect to (default: 9000)")
    fs.StringVar(&c.Config.ClickHouse.Database, "clickhouse-database", config.DefaultClickHouseDatabase, "Select the current default ClickHouse database (default: 'default').")
    fs.StringVar(&c.Config.ClickHouse.User, "clickhouse-user", config.DefaultClickHouseUser, "The ClickHouse username to connect with (default: 'default').")
    fs.StringVar(&c.Config.ClickHouse.Password, "clickhouse-password", config.DefaultClickHousePassword, "The ClickHouse password (default: '').")
}

func (c *RootConfig) Exec(context.Context, []string) error {
    return flag.ErrHelp
}
