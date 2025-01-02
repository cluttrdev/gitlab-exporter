package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
)

type FetchArtifactsConfig struct {
	FetchConfig

	projectPath string
	jobId       int64
	fileType    string

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
