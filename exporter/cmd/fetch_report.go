package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/cluttrdev/cli"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/config"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/tasks"
)

type FetchReportConfig struct {
	FetchConfig

	projectPath   string
	jobId         int64
	fileType      string
	artifactPaths stringList
}

type stringList []string

func (f *stringList) String() string {
	return fmt.Sprintf("%v", []string(*f))
}

func (f *stringList) Set(value string) error {
	values := strings.Split(value, ",")
	*f = append(*f, values...)
	return nil
}

func NewFetchReportCmd(out io.Writer) *cli.Command {
	cfg := FetchReportConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet("fetch report", flag.ExitOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "report",
		ShortUsage: fmt.Sprintf("%s fetch report [option]...", exeName),
		ShortHelp:  "Fetch job report",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *FetchReportConfig) RegisterFlags(fs *flag.FlagSet) {
	c.FetchConfig.RegisterFlags(fs)

	fs.StringVar(&c.projectPath, "project-path", "", "The full project path.")
	fs.Int64Var(&c.jobId, "job-id", 0, "The job id.")
	fs.StringVar(&c.fileType, "file-type", "", "The report file type.")
	fs.Var(&c.artifactPaths, "artifact-paths", "Comma separated list of artifact paths.")
}

func (c *FetchReportConfig) Exec(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("too many arguments: %v", args)
	}

	if c.projectPath == "" {
		return fmt.Errorf("missing required option: --project-path")
	}
	if c.jobId == 0 {
		return fmt.Errorf("missing required option: --job-id")
	}

	cfg := config.Default()
	if err := loadConfig(c.FetchConfig.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	// create gitlab client
	glab, err := createGitLabClient(cfg)
	if err != nil {
		return fmt.Errorf("create gitlab client: %w", err)
	}

	jobId := "gid://gitlab/Ci::Build/" + strconv.FormatInt(c.jobId, 10)

	artifacts, err := glab.GraphQL.GetProjectJobArtifacts(ctx, c.projectPath, jobId)
	if err != nil {
		return fmt.Errorf("get job artifacts list: %w", err)
	}

	i := slices.IndexFunc(artifacts, func(a graphql.JobArtifactFields) bool {
		return a.FileType != nil && strings.EqualFold(string(*a.FileType), c.fileType)
	})
	if i == -1 {
		return fmt.Errorf("file type not found: %s", c.fileType)
	}

	var data any
	switch c.fileType {
	case "junit":
		data, err = c.fetchJunitReport(ctx, glab, artifacts[i], c.artifactPaths)
	case "cobertura":
		data, err = c.fetchCoberturaReport(ctx, glab, artifacts[i], c.artifactPaths)
	default:
		return fmt.Errorf("unsupported file type: %s", c.fileType)
	}

	out, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	_, err = fmt.Fprint(c.FetchConfig.RootConfig.out, string(out))
	if err != nil {
		return err
	}

	return nil
}

func (c *FetchReportConfig) fetchJunitReport(ctx context.Context, glab *gitlab.Client, artifact graphql.JobArtifactFields, artifactPaths []string) (any, error) {
	tr, ts, tc, err := tasks.FetchProjectJobArtifactJunitReports(ctx, glab, artifact, artifactPaths)
	if err != nil {
		return nil, fmt.Errorf("fetch report: junit: %w", err)
	}
	return map[string]any{
		"reports": tr,
		"suites":  ts,
		"cases":   tc,
	}, nil
}

func (c *FetchReportConfig) fetchCoberturaReport(ctx context.Context, glab *gitlab.Client, artifact graphql.JobArtifactFields, artifactPaths []string) (any, error) {
	cr, cp, cc, cm, err := tasks.FetchProjectJobArtifactCoberturaReports(ctx, glab, artifact, artifactPaths)
	if err != nil {
		return nil, fmt.Errorf("fetch report: cobertura: %w", err)
	}
	return map[string]any{
		"reports":  cr,
		"packages": cp,
		"classes":  cc,
		"methods":  cm,
	}, nil
}
