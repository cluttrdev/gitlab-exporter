package cmd

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/cluttrdev/cli"

	"go.cluttr.dev/gitlab-exporter/internal/cobertura"
	"go.cluttr.dev/gitlab-exporter/internal/config"
	"go.cluttr.dev/gitlab-exporter/internal/junitxml"
)

type FetchArtifactsConfig struct {
	FetchConfig

	projectPath string
	jobId       int64
	fileType    string

	parse  bool
	output string
}

func NewFetchArtifactsCmd(out io.Writer) *cli.Command {
	cfg := FetchArtifactsConfig{
		FetchConfig: FetchConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet(fmt.Sprintf("%s download artifacts", exeName), flag.ContinueOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "artifacts",
		ShortUsage: fmt.Sprintf("%s fetch artifacts [option]...", exeName),
		ShortHelp:  "Download job artifacts file.",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *FetchArtifactsConfig) RegisterFlags(fs *flag.FlagSet) {
	c.FetchConfig.RegisterFlags(fs)

	fs.StringVar(&c.projectPath, "project-path", "", "The full project path.")
	fs.Int64Var(&c.jobId, "job-id", 0, "The job id.")
	fs.StringVar(&c.fileType, "file-type", "", "The artifacts file type.")

	fs.BoolVar(&c.parse, "parse", false, "Whether to parse the artifacts file before writing to the output file.")
	fs.StringVar(&c.output, "output", "", "The output file.")
}

func (c *FetchArtifactsConfig) Exec(ctx context.Context, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("too many arguments: %v", args)
	}

	if c.projectPath == "" {
		return fmt.Errorf("missing required option: --project-path")
	}
	if c.jobId == 0 {
		return fmt.Errorf("missing required option: --job-id")
	}
	if c.output == "" {
		return fmt.Errorf("missing required option: --output")
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

	r, err := glab.HTTP.GetProjectJobArtifactsFile(ctx, c.projectPath, c.jobId, c.fileType)
	if err != nil {
		return fmt.Errorf("fetch artifact file: %w", err)
	}

	if c.parse {
		var (
			b   []byte
			err error
		)
		switch c.fileType {
		case "junit":
			b, err = c.parseJunitArtifact(r)
		case "cobertura":
			b, err = c.parseCoberturaArtifact(r)
		default:
			return fmt.Errorf("unsupported file type for parsing: %s", c.fileType)
		}
		if err != nil {
			return fmt.Errorf("parse artifact file: %w", err)
		}

		r = bytes.NewReader(b)
	}

	file, err := os.Create(c.output)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	return nil
}

func (c *FetchArtifactsConfig) parseJunitArtifact(r io.Reader) ([]byte, error) {
	reader, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	reports, err := junitxml.ParseMany(reader)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(reports)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *FetchArtifactsConfig) parseCoberturaArtifact(r io.Reader) ([]byte, error) {
	report, err := cobertura.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	b, err := json.Marshal(report)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}

	return b, nil
}
