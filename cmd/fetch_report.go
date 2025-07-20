package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/cluttrdev/cli"
	"go.cluttr.dev/junitxml"

	"go.cluttr.dev/gitlab-exporter/internal/cobertura"
	"go.cluttr.dev/gitlab-exporter/internal/config"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/internal/types"
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
	if len(c.artifactPaths) == 0 {
		return fmt.Errorf("missing report artifact paths")
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

	var b []byte
	switch c.fileType {
	case "junit":
		b, err = c.fetchJunitReports(ctx, glab)
	case "cobertura":
		b, err = c.fetchCoberturaReport(ctx, glab)
	default:
		return fmt.Errorf("unsupported file type: %s", c.fileType)
	}
	if err != nil {
		return err
	}

	fmt.Fprint(c.FetchConfig.RootConfig.out, string(b))

	return nil
}

func (c *FetchReportConfig) fetchJunitReports(ctx context.Context, glab *gitlab.Client) ([]byte, error) {
	type result struct {
		TestReport types.TestReport  `json:"testreport"`
		TestSuites []types.TestSuite `json:"testsuites"`
		TestCases  []types.TestCase  `json:"testcases"`
	}

	var results []result

	for _, path := range c.artifactPaths {
		reader, err := glab.Rest.GetProjectJobArtifact(ctx, c.projectPath, c.jobId, path)
		if errors.Is(err, gitlab.ErrNotFound) {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("download file: %w", err)
		}

		report, err := junitxml.Parse(reader)
		if err != nil {
			return nil, fmt.Errorf("parse file: %w", err)
		}

		tr, ts, tc := types.ConvertTestReport(report, types.JobReference{
			Id: c.jobId,
			Pipeline: types.PipelineReference{
				Project: types.ProjectReference{
					FullPath: c.projectPath,
				},
			},
		})
		results = append(results, result{
			TestReport: tr,
			TestSuites: ts,
			TestCases:  tc,
		})
	}

	b, err := json.Marshal(results)
	if err != nil {
		return nil, fmt.Errorf("error marshalling pipeline testreport: %w", err)
	}

	return b, nil
}

func (c *FetchReportConfig) fetchCoberturaReport(ctx context.Context, glab *gitlab.Client) ([]byte, error) {
	var (
		reports  []types.CoverageReport
		packages []types.CoveragePackage
		classes  []types.CoverageClass
		methods  []types.CoverageMethod
	)

	reportCounter := 0
	for _, path := range c.artifactPaths {
		reader, err := glab.Rest.GetProjectJobArtifact(ctx, c.projectPath, c.jobId, path)
		if errors.Is(err, gitlab.ErrNotFound) {
			continue
		} else if err != nil {
			return nil, fmt.Errorf("download file: %w", err)
		}

		report, err := cobertura.Parse(reader)
		if err != nil {
			return nil, fmt.Errorf("parse file: %w", err)
		}

		cr, cp, cc, cm := cobertura.ConvertCoverageReport(reportCounter, report, types.JobReference{
			Id: c.jobId,
			Pipeline: types.PipelineReference{
				Project: types.ProjectReference{
					FullPath: c.projectPath,
				},
			},
		})
		reports = append(reports, cr)
		packages = append(packages, cp...)
		classes = append(classes, cc...)
		methods = append(methods, cm...)

		reportCounter++
	}

	b, err := json.Marshal(map[string]any{
		"reports":  reports,
		"packages": packages,
		"classes":  classes,
		"methods":  methods,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshalling cobertura reports: %w", err)
	}

	return b, nil
}
