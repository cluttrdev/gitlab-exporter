package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strconv"

	"github.com/cluttrdev/cli"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/config"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/exporter"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/types"

	grpc_client "go.cluttr.dev/gitlab-exporter/grpc/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ExportPipelineConfig struct {
	ExportConfig

	exportSections    bool
	exportTestReports bool
	exportTraces      bool
	exportMetrics     bool
}

func NewExportPipelineCmd(out io.Writer) *cli.Command {
	cfg := &ExportPipelineConfig{
		ExportConfig: ExportConfig{
			RootConfig: RootConfig{
				out:   out,
				flags: flag.NewFlagSet(fmt.Sprintf("%s export pipeline", exeName), flag.ContinueOnError),
			},
		},
	}

	cfg.RegisterFlags(cfg.flags)

	return &cli.Command{
		Name:       "pipeline",
		ShortUsage: fmt.Sprintf("%s export pipeline [option]... project_id pipeline_id", exeName),
		ShortHelp:  "Export pipeline data",
		Flags:      cfg.flags,
		Exec:       cfg.Exec,
	}
}

func (c *ExportPipelineConfig) RegisterFlags(fs *flag.FlagSet) {
	c.ExportConfig.RegisterFlags(fs)

	fs.BoolVar(&c.exportSections, "export-sections", true, "Export job sections. (default: true)")
	fs.BoolVar(&c.exportTraces, "export-traces", true, "Export pipeline trace. (default: true)")
	fs.BoolVar(&c.exportTestReports, "export-testreports", true, "Export pipeline test reports. (default: true)")
	fs.BoolVar(&c.exportMetrics, "export-metrics", true, "Export job log embedded metrics. (default: true)")
}

func (c *ExportPipelineConfig) Exec(ctx context.Context, args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("invalid number of positional arguments: %v", args)
	}

	projectId, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing `project_id` argument: %w", err)
	}

	pipelineId, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return fmt.Errorf("error parsing `pipeline_id` argument: %w", err)
	}

	cfg := config.Default()
	if err := loadConfig(c.ExportConfig.RootConfig.filename, c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// create gitlab client
	glab, err := createGitLabClient(cfg)
	if err != nil {
		return fmt.Errorf("create gitlab client: %w", err)
	}

	// initialize grpc clients (external only)
	var recorderConfigs []config.Recorder
	var clients []*grpc_client.Client
	for _, endpoint := range cfg.Endpoints {
		recorderConfigs = append(recorderConfigs, config.Recorder{
			Address: endpoint.Address,
			Mode:    config.RecorderModeExternal,
			Enabled: true,
		})
	}
	recorderConfigs = append(recorderConfigs, cfg.Recorders...)
	for _, rec := range recorderConfigs {
		if !rec.Enabled {
			continue
		}
		if rec.Mode != config.RecorderModeExternal {
			continue
		}

		if rec.Address == "" {
			return fmt.Errorf("external recorder %s: address is required", rec.Type)
		}

		client, err := grpc_client.NewCLient(rec.Address,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return fmt.Errorf("connect to external recorder %s at %s: %w", rec.Type, rec.Address, err)
		}

		clients = append(clients, client)
	}

	// create exporter
	exp := exporter.New()
	for _, client := range clients {
		if err := exp.AddClient(client); err != nil {
			return fmt.Errorf("add grpc client: %w", err)
		}
	}

	projectGid := graphql.GlobalIdProjectPrefix + strconv.FormatInt(projectId, 10)
	pipelineGid := graphql.GlobalIdPipelinePrefix + strconv.FormatInt(pipelineId, 10)

	pipelineFields, err := glab.GraphQL.GetProjectPipeline(ctx, projectGid, pipelineGid)
	if err != nil {
		return fmt.Errorf("get pipeline fields: %w", err)
	}

	pipeline, err := graphql.ConvertPipeline(pipelineFields)
	if err != nil {
		return fmt.Errorf("convert pipeline fields: %w", err)
	}

	return exp.ExportPipelines(ctx, []types.Pipeline{pipeline})
}
