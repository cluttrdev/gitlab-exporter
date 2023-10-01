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
	"strings"
	"syscall"

	"github.com/peterbourgon/ff/v3/ffcli"
	"golang.org/x/exp/slices"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/worker"
)

type RunConfig struct {
	rootConfig *RootConfig
	out        io.Writer
	projects   projectList

	flags *flag.FlagSet
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
	fs := flag.NewFlagSet(fmt.Sprintf("%s run", exeName), flag.ContinueOnError)

	config := RunConfig{
		rootConfig: rootConfig,
		out:        out,

		flags: fs,
	}

	config.RegisterFlags(fs)

	return &ffcli.Command{
		Name:       "run",
		ShortUsage: fmt.Sprintf("%s run [flags]", exeName),
		ShortHelp:  "Run in daemon mode",
		UsageFunc:  usageFunc,
		FlagSet:    fs,
		Options:    rootCmdOptions,
		Exec:       config.Exec,
	}
}

func (c *RunConfig) RegisterFlags(fs *flag.FlagSet) {
	c.rootConfig.RegisterFlags(fs)

	fs.Var(&c.projects, "projects", "Comma separated list of project ids.")
}

func (c *RunConfig) Exec(ctx context.Context, _ []string) error {
	// configure logging
	log.SetOutput(c.out)

	// create config from root flags
	log.Printf("Loading configuration from %s\n", c.rootConfig.filename)
	cfg, err := newConfig(c.rootConfig.filename, c.flags)
	if err != nil {
		return err
	}

	// add projects pass to run command
	for _, pid := range c.projects {
		exists := slices.ContainsFunc(cfg.Projects, func(p config.Project) bool {
			return p.Id == pid
		})

		if !exists {
			cfg.Projects = append(cfg.Projects, config.Project{
				ProjectSettings: *config.DefaultProjectSettings(),
				Id:              pid,
			})
		}
	}

	// setup controller
	ctl, err := newController(cfg)
	if err != nil {
		return err
	}

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
	writeConfig(cfg, c.out)

	// run daemon

	workers := []worker.Worker{}

	for _, pc := range cfg.Projects {
		if pc.CatchUp.Enabled {
			workers = append(workers, worker.NewCatchUpProjectWorker(ctl, &pc))
		}
		workers = append(workers, worker.NewExportProjectWorker(ctl, &pc))
	}

	log.Println("Starting workers...")
	for _, w := range workers {
		go w.Start()
	}

	<-ctx.Done()

	log.Println("Stopping workers...")
	for _, w := range workers {
		w.Stop()
	}

	return nil
}
