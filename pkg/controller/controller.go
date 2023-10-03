package controller

import (
	"context"
	"fmt"

	clickhouse "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/clickhouse"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	gitlab "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/worker"
)

type Controller struct {
	Config     config.Config
	GitLab     *gitlab.Client
	ClickHouse *clickhouse.Client

	workers []worker.Worker
}

func NewController(cfg config.Config) (c Controller, err error) {
	c.Config = cfg

	if err = c.configureGitLabClient(cfg.GitLab); err != nil {
		return
	}

	if err = c.configureClickHouseClient(cfg.ClickHouse); err != nil {
		return
	}

	return c, nil
}

func (c *Controller) configureGitLabClient(cfg config.GitLab) (err error) {
	c.GitLab, err = gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL:   cfg.Api.URL,
		Token: cfg.Api.Token,

		RateLimit: cfg.Client.Rate.Limit,
	})
	return
}

func (c *Controller) configureClickHouseClient(cfg config.ClickHouse) (err error) {
	c.ClickHouse, err = clickhouse.NewClickHouseClient(clickhouse.ClientConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Database: cfg.Database,
		User:     cfg.User,
		Password: cfg.Password,
	})
	return
}

func (c *Controller) projectConfig(pid int64) (*config.Project, error) {
	for _, cfg := range c.Config.Projects {
		if cfg.Id == pid {
			return &cfg, nil
		}
	}
	return nil, fmt.Errorf("Config not found for project id %d", pid)
}

func (c *Controller) Init(ctx context.Context) error {
	if err := c.ClickHouse.CreateDatabase(ctx); err != nil {
		return err
	}

	return c.ClickHouse.CreateTables(ctx)
}
