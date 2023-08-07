package controller

import (
	"context"
	"fmt"
	"time"

	clickhouse "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/clickhouse"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	gitlab "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
)

type Controller struct {
	Config     config.Config
	GitLab     *gitlab.Client
	ClickHouse *clickhouse.Client
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
		URL:   cfg.URL,
		Token: cfg.Token,
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

func (c *Controller) Init(ctx context.Context) error {
	if err := c.ClickHouse.CreateDatabase(ctx); err != nil {
		return err
	}

	return c.ClickHouse.CreateTables(ctx)
}

func (c *Controller) QueryLatestProjectPipelineUpdate(ctx context.Context) (map[int64]time.Time, error) {
	const (
		msPerSecond float64 = 1000
	)

	var results []struct {
		ProjectID    int64   `ch:"project_id"`
		LatestUpdate float64 `ch:"latest_update"`
	}

	query := `SELECT project_id, max(updated_at) AS latest_update FROM gitlab_ci.pipelines GROUP BY project_id`
	if err := c.ClickHouse.Conn.Select(ctx, &results, query); err != nil {
		return nil, fmt.Errorf("[controller.QueryLatestProjectPipelineUpdate] %w", err)
	}

	m := map[int64]time.Time{}
	for _, r := range results {
		m[r.ProjectID] = time.UnixMilli(int64(r.LatestUpdate * msPerSecond)).UTC()
	}

	return m, nil

}

func (c *Controller) ExportPipeline(ctx context.Context, projectID int64, pipelineID int64) error {
	ph, err := c.GitLab.GetPipelineHierarchy(ctx, projectID, pipelineID)
	if err != nil {
		return fmt.Errorf("[controller.ExportPipeline/GetHierarchy] %w", err)
	}

	if err = clickhouse.InsertPipelineHierarchy(ctx, ph, c.ClickHouse); err != nil {
		return fmt.Errorf("[controller.ExportPipeline/InsertHierarchy] %w", err)
	}

	pts := ph.GetAllTraces()
	if err = clickhouse.InsertTraces(ctx, pts, c.ClickHouse); err != nil {
		return fmt.Errorf("[controller.ExportPipeline/GetTraces] %w", err)
	}

	if err = clickhouse.InsertTraces(ctx, pts, c.ClickHouse); err != nil {
		return fmt.Errorf("[controller.ExportPipeline/InsertTraces] %w", err)
	}

	return nil
}
