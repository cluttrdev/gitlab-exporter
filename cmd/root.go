package cmd

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/controller"
)

const (
	exeName      string = "gitlab-clickhouse-exporter"
	envVarPrefix string = "GLCHE"
)

var (
	rootCmdOptions = []ff.Option{
		ff.WithEnvVarPrefix(envVarPrefix),
	}
)

type RootConfig struct {
	filename string

	flags *flag.FlagSet
}

func NewRootCmd() (*ffcli.Command, *RootConfig) {
	fs := flag.NewFlagSet(exeName, flag.ContinueOnError)

	cfg := RootConfig{
		filename: "",

		flags: fs,
	}

	cfg.RegisterFlags(cfg.flags)

	return &ffcli.Command{
		Name:       exeName,
		ShortUsage: fmt.Sprintf("%s <subcommand> [flags] [<args>...]", exeName),
		UsageFunc:  usageFunc,
		FlagSet:    cfg.flags,
		Options:    rootCmdOptions,
		Exec:       cfg.Exec,
	}, &cfg
}

func (c *RootConfig) RegisterFlags(fs *flag.FlagSet) {
	defaults := config.Default()

	fs.String("gitlab-api-url", defaults.GitLab.Api.URL, fmt.Sprintf("The GitLab API URL (default: '%s').", defaults.GitLab.Api.URL))
	fs.String("gitlab-api-token", defaults.GitLab.Api.Token, fmt.Sprintf("The GitLab API Token (default: '%s').", defaults.GitLab.Api.Token))

	fs.String("clickhouse-host", defaults.ClickHouse.Host, fmt.Sprintf("The ClickHouse server name (default: '%s').", defaults.ClickHouse.Host))
	fs.String("clickhouse-port", defaults.ClickHouse.Port, fmt.Sprintf("The ClickHouse port to connect to (default: '%s')", defaults.ClickHouse.Port))
	fs.String("clickhouse-database", defaults.ClickHouse.Database, fmt.Sprintf("Select the current default ClickHouse database (default: '%s').", defaults.ClickHouse.Database))
	fs.String("clickhouse-user", defaults.ClickHouse.User, fmt.Sprintf("The ClickHouse username to connect with (default: '%s').", defaults.ClickHouse.User))
	fs.String("clickhouse-password", defaults.ClickHouse.Password, fmt.Sprintf("The ClickHouse password (default: '%s').", defaults.ClickHouse.Password))

	fs.StringVar(&c.filename, "config", "", "Configuration file to use.")
}

func (c *RootConfig) Exec(context.Context, []string) error {
	return flag.ErrHelp
}

func newConfig(filename string, flags *flag.FlagSet) (*config.Config, error) {
	cfg := config.Default()

	if filename != "" {
		if err := config.LoadFile(filename, cfg); err != nil {
			return nil, err
		}
	}

	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "gitlab-api-url":
			cfg.GitLab.Api.URL = f.Value.String()
		case "gitlab-api-token":
			cfg.GitLab.Api.Token = f.Value.String()
		case "clickhouse-host":
			cfg.ClickHouse.Host = f.Value.String()
		case "clickhouse-port":
			cfg.ClickHouse.Port = f.Value.String()
		case "clickhouse-database":
			cfg.ClickHouse.Database = f.Value.String()
		case "clickhouse-user":
			cfg.ClickHouse.User = f.Value.String()
		case "clickhouse-password":
			cfg.ClickHouse.Password = f.Value.String()
		}
	})

	return cfg, nil
}

func newController(cfg *config.Config) (*controller.Controller, error) {
	ctl, err := controller.NewController(*cfg)
	if err != nil {
		return nil, fmt.Errorf("error constructing controller: %w", err)
	}

	return &ctl, nil
}

func writeConfig(cfg *config.Config, out io.Writer) {
	fmt.Fprintln(out, "----")
	fmt.Fprintf(out, "GitLab URL: %s\n", cfg.GitLab.Api.URL)
	fmt.Fprintf(out, "GitLab Token: %x\n", sha256String(cfg.GitLab.Api.Token))
	fmt.Fprintln(out, "----")
	fmt.Fprintf(out, "ClickHouse Host: %s\n", cfg.ClickHouse.Host)
	fmt.Fprintf(out, "ClickHouse Port: %s\n", cfg.ClickHouse.Port)
	fmt.Fprintf(out, "ClickHouse Database: %s\n", cfg.ClickHouse.Database)
	fmt.Fprintf(out, "ClickHouse User: %s\n", cfg.ClickHouse.User)
	fmt.Fprintf(out, "ClickHouse Password: %x\n", sha256String(cfg.ClickHouse.Password))
	fmt.Fprintln(out, "----")

	projects := []int64{}
	for _, p := range cfg.Projects {
		projects = append(projects, p.Id)
	}
	fmt.Fprintf(out, "Projects: %v\n", projects)
	fmt.Fprintln(out, "----")
}

func sha256String(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}
