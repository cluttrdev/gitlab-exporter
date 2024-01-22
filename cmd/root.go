package cmd

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
)

const (
	exeName      string = "gitlab-exporter"
	envVarPrefix string = "GLE"
)

type RootConfig struct {
	filename string

	out   io.Writer
	flags flag.FlagSet
}

func NewRootCmd(out io.Writer) *cli.Command {
	fs := flag.NewFlagSet(exeName, flag.ContinueOnError)

	cfg := RootConfig{
		filename: "",

		out:   out,
		flags: *fs,
	}

	cfg.RegisterFlags(&cfg.flags)

	return &cli.Command{
		Name:       exeName,
		ShortUsage: fmt.Sprintf("%s <subcommand> [option]... [arg]...", exeName),
		Flags:      &cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *RootConfig) RegisterFlags(fs *flag.FlagSet) {
	defaults := config.Default()

	fs.String("gitlab-api-url", defaults.GitLab.Api.URL, fmt.Sprintf("The GitLab API URL (default: '%s').", defaults.GitLab.Api.URL))
	fs.String("gitlab-api-token", defaults.GitLab.Api.Token, fmt.Sprintf("The GitLab API Token (default: '%s').", defaults.GitLab.Api.Token))

	fs.StringVar(&c.filename, "config", "", "Configuration file to use.")
}

func (c *RootConfig) Exec(context.Context, []string) error {
	return flag.ErrHelp
}

func loadConfig(filename string, flags *flag.FlagSet, cfg *config.Config) error {
	if filename != "" {
		if err := config.LoadFile(filename, cfg); err != nil {
			return err
		}
	}

	flags.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "gitlab-api-url":
			cfg.GitLab.Api.URL = f.Value.String()
		case "gitlab-api-token":
			cfg.GitLab.Api.Token = f.Value.String()
		}
	})

	return nil
}

func writeConfig(cfg config.Config, out io.Writer) {
	fmt.Fprintln(out, "----")
	fmt.Fprintf(out, "GitLab URL: %s\n", cfg.GitLab.Api.URL)
	fmt.Fprintf(out, "GitLab Token: %x\n", sha256String(cfg.GitLab.Api.Token))
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
