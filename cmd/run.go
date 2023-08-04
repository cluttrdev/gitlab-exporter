package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/controller"
)

func Run(ctx context.Context, cfg config.Config, out io.Writer) error {
    log.SetOutput(out)

    flag.Parse()
    if flag.NArg() != 2 {
        return fmt.Errorf("usage: gitlab-clickhouse-exporter PROJECT_ID PIPELINE_ID")
    }
    
    projectID, err := strconv.ParseInt(flag.Arg(0), 10, 64)
    if err != nil {
        log.Fatal(err)
    }
    pipelineID, err := strconv.ParseInt(flag.Arg(1), 10, 64)
    if err != nil {
        return err
    }

    ctl, err := controller.NewController(cfg)
    if err != nil {
        return err
    }

    if err = ctl.Start(ctx); err != nil {
        return err
    }

    return ctl.ExportPipeline(ctx, projectID, pipelineID)
}
