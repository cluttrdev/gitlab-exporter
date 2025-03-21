package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/cluttrdev/cli"
	"go.cluttr.dev/gitlab-exporter/internal/config"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab/graphql"
)

type FetchProjectsConfig struct {
	FetchConfig

	ids           []int64
	updatedAfter  *time.Time
	updatedBefore *time.Time
}

func NewFetchProjectsCommand(out io.Writer) *cli.Command {
	cfg := FetchProjectsConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet(fmt.Sprintf("%s fetch project updates", exeName), flag.ExitOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "projects",
		ShortUsage: fmt.Sprintf("%s fetch projects [option]...", exeName),
		ShortHelp:  "Fetch projects",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *FetchProjectsConfig) RegisterFlags(fs *flag.FlagSet) {
	c.FetchConfig.RegisterFlags(fs)

	fs.Func("ids", "", func(s string) error {
		for _, s_ := range strings.Split(s, ",") {
			id, err := strconv.ParseInt(s_, 10, 64)
			if err != nil {
				return err
			}
			c.ids = append(c.ids, id)
		}
		return nil
	})
	fs.Func("updated-after", "", func(s string) error {
		t, err := parseTimeISO8601(s)
		if err != nil {
			return err
		}
		c.updatedAfter = &t
		return nil
	})
	fs.Func("updated-before", "", func(s string) error {
		t, err := parseTimeISO8601(s)
		if err != nil {
			return err
		}
		c.updatedBefore = &t
		return nil
	})
}

func (c *FetchProjectsConfig) Exec(ctx context.Context, args []string) error {
	cfg := config.Default()
	if err := loadConfig(c.FetchConfig.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	glab, err := createGitLabClient(cfg)
	if err != nil {
		return fmt.Errorf("create gitlab client: %w", err)
	}

	var gids []string
	if len(c.ids) > 0 {
		gids = make([]string, 0, len(c.ids))
		for _, id := range c.ids {
			gids = append(gids, graphql.FormatId(id, graphql.GlobalIdProjectPrefix))
		}
	}

	updates, err := glab.GraphQL.GetProjects(ctx, gids, c.updatedAfter, c.updatedBefore)
	if err != nil {
		return err
	}

	out, err := json.Marshal(updates)
	if err != nil {
		return err
	}

	fmt.Fprintln(c.FetchConfig.RootConfig.out, string(out))
	return nil
}

func parseTimeISO8601(s string) (time.Time, error) {
	layouts := []string{
		"2006-01-02T15:04:05+07:00",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}

	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t.In(time.UTC), nil
		}
	}

	return time.Time{}, fmt.Errorf("failed to parse time from %q", s)
}
