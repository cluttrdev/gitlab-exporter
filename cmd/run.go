package cmd

import (
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"

	gitlab "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
)

type RunConfig struct {
    rootConfig *RootConfig
    out io.Writer
    projects projectList
}

type projectList []int64

func (f *projectList) String() string {
    return fmt.Sprintf("%v", []int64(*f))
}

func (f *projectList) Set(value string) error {
    values := strings.Split(value, ",")
    for _, s := range values {
        v, err := strconv.ParseInt(s, 10, 64)
        if err != nil {
            return err
        }
        *f = append(*f, v)
    }
    return nil
}


func NewRunCmd(rootConfig *RootConfig, out io.Writer) *ffcli.Command {
    config := RunConfig{
        rootConfig: rootConfig,
        out: out,
    }

    fs := flag.NewFlagSet(fmt.Sprintf("%s run", exeName), flag.ExitOnError)
    config.RegisterFlags(fs)
    config.rootConfig.RegisterFlags(fs)

    return &ffcli.Command{
        Name: "run",
        ShortUsage: fmt.Sprintf("%s run [flags]", exeName),
        ShortHelp: "Run in daemon mode",
        FlagSet: fs,
        Options: []ff.Option{ff.WithEnvVarPrefix(envVarPrefix)},
        Exec: config.Exec,
    }
}

func (c *RunConfig) RegisterFlags(fs *flag.FlagSet) {
    fs.Var(&c.projects, "projects", "Comma separated list of project ids.")
}

func (c *RunConfig) Exec(ctx context.Context, _ []string) error {
	// configure logging
	log.SetOutput(c.out)

    ctl := c.rootConfig.Controller

	// init controller
    if err := ctl.Init(ctx); err != nil {
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

    // log configuration
    printRunConfig(c, c.out)
    
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

			for _, project := range c.projects {
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

func printRunConfig(cfg *RunConfig, out io.Writer) {
    fmt.Fprintf(out, "GitLab URL: %s\n", cfg.rootConfig.Config.GitLab.URL)
    fmt.Fprintf(out, "GitLab Token: %x\n", sha256String(cfg.rootConfig.Config.GitLab.Token))
    fmt.Fprintln(out, "---")
    fmt.Fprintf(out, "ClickHouse Host: %s\n", cfg.rootConfig.Config.ClickHouse.Host)
    fmt.Fprintf(out, "ClickHouse Port: %d\n", cfg.rootConfig.Config.ClickHouse.Port)
    fmt.Fprintf(out, "ClickHouse Database: %s\n", cfg.rootConfig.Config.ClickHouse.Database)
    fmt.Fprintf(out, "ClickHouse User: %s\n", cfg.rootConfig.Config.ClickHouse.User)
    fmt.Fprintf(out, "ClickHouse Password: %x\n", sha256String(cfg.rootConfig.Config.ClickHouse.Password))
    fmt.Fprintln(out, "---")
    fmt.Fprintf(out, "Projects: %v\n", cfg.projects)
}

func sha256String(s string) []byte {
    h := sha256.New()
    h.Write([]byte(s))
    return h.Sum(nil)
}
