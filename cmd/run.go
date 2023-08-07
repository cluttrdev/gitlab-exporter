package cmd

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/controller"
	gitlab "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
)

func Run(ctx context.Context, cfg config.Config, out io.Writer) error {
	// configure logging
	log.SetOutput(out)

	// parse commandline arguments and flags
	flag.Parse()
	if flag.NArg() != 1 {
		return fmt.Errorf("usage: gitlab-clickhouse-exporter PROJECT_ID")
	}

	projectID, err := strconv.ParseInt(flag.Arg(0), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	// setup controller
	ctl, err := controller.NewController(cfg)
	if err != nil {
		return err
	}

	if err = ctl.Init(ctx); err != nil {
		return err
	}

	// setup daemon
	ctx, cancel := context.WithCancel(ctx)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	go func() {
		select {
		case <-signalChan:
			log.Println("Got SIGINT/SIGTERM, exiting")
			cancel()
		case <-ctx.Done():
			log.Println("Done")
		}
	}()

	// run daemon
	var firstRun bool = true
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if firstRun {
				ticker.Stop()
				ticker = time.NewTicker(60 * time.Second)
				firstRun = false
			}

			pu, err := ctl.QueryLatestProjectPipelineUpdate(ctx)
			if err != nil {
				log.Println(err)
			}

			scope := "finished"
			opt := &gitlab.ListProjectPipelineOptions{
				PerPage: 100,
				Page:    1,

				Scope: &scope,
			}

			projects := []int64{projectID}
			for _, project := range projects {
				pis, err := ctl.GitLab.ListProjectPipelines(ctx, project, opt)
				if err != nil {
					log.Println(err)
					continue
				}

				for _, pi := range pis {
					if pi.UpdatedAt.Before(pu[project]) {
						continue
					}

					log.Printf("Exporting projects/%d/pipelines/%d\n", project, pi.ID)
					if err := ctl.ExportPipeline(ctx, project, pi.ID); err != nil {
						log.Printf("Error: %s\n", err)
					}
				}
			}
		}
	}
}
