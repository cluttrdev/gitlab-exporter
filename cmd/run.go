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

	"golang.org/x/exp/slices"

	"github.com/cluttrdev/cli"

	"github.com/cluttrdev/gitlab-exporter/internal/config"
	"github.com/cluttrdev/gitlab-exporter/internal/controller"
	"github.com/cluttrdev/gitlab-exporter/internal/server"
)

type RunConfig struct {
	RootConfig

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

func NewRunCmd(out io.Writer) *cli.Command {
	fs := flag.NewFlagSet(fmt.Sprintf("%s run", exeName), flag.ContinueOnError)

	config := RunConfig{
		RootConfig: RootConfig{
			out: out,
		},
	}

	config.RegisterFlags(fs)

	return &cli.Command{
		Name:       "run",
		ShortUsage: fmt.Sprintf("%s run [option]...", exeName),
		ShortHelp:  "Run in daemon mode",
		Flags:      fs,
		Exec:       config.Exec,
	}
}

func (c *RunConfig) RegisterFlags(fs *flag.FlagSet) {
	c.RootConfig.RegisterFlags(fs)

	fs.Var(&c.projects, "projects", "Comma separated list of project ids.")
}

func (c *RunConfig) Exec(ctx context.Context, _ []string) error {
	// configure logging
	log.SetOutput(c.out)

	// load configuration
	log.Printf("Loading configuration from %s\n", c.RootConfig.filename)
	cfg := config.Default()
	if err := loadConfig(c.RootConfig.filename, &c.flags, &cfg); err != nil {
		return fmt.Errorf("error loading configuration: %w", err)
	}

	// add projects passed to run command
	for _, pid := range c.projects {
		exists := slices.ContainsFunc(cfg.Projects, func(p config.Project) bool {
			return p.Id == pid
		})

		if !exists {
			cfg.Projects = append(cfg.Projects, config.Project{
				ProjectSettings: config.DefaultProjectSettings(),
				Id:              pid,
			})
		}
	}

	// log configuration
	writeConfig(cfg, c.out)

	// setup controller
	ctl, err := controller.NewController(cfg)
	if err != nil {
		return fmt.Errorf("error constructing controller: %w", err)
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

	go startServer(ctx, cfg.Server, ctl)

	// run daemon
	return ctl.Run(ctx)
}

func startServer(ctx context.Context, cfg config.Server, ctl *controller.Controller) {
	srv := server.New(server.ServerConfig{
		Host:  cfg.Host,
		Port:  cfg.Port,
		Debug: false,

		ReadinessCheck: func() error { return ctl.CheckReadiness(ctx) },
	})

	if err := srv.Serve(ctx); err != nil {
		log.Printf("error during server shutdown: %v\n", err)
	}
}
