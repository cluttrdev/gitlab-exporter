package controller

import (
	"context"
	"fmt"
	"time"

	clickhouse "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/clickhouse"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/config"
	gitlab "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/gitlab"
	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
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

func (c *Controller) QueryLatestProjectPipelineUpdates(ctx context.Context, projectID int64) (map[int64]time.Time, error) {
	const (
		msPerSecond float64 = 1000
	)

	var results []struct {
		PipelineID   int64   `ch:"id"`
		LatestUpdate float64 `ch:"latest_update"`
	}

	query := `SELECT id, max(updated_at) AS latest_update FROM gitlab_ci.pipelines GROUP BY id`
	if err := c.ClickHouse.Conn.Select(ctx, &results, query); err != nil {
		return nil, fmt.Errorf("[controller.QueryLatestProjectPipelineUpdates] %w", err)
	}

	m := map[int64]time.Time{}
	for _, r := range results {
		m[r.PipelineID] = time.UnixMilli(int64(r.LatestUpdate * msPerSecond)).UTC()
	}

	return m, nil
}

func (c *Controller) ExportPipeline(ctx context.Context, projectID int64, pipelineID int64) error {
	if err := <-c.exportPipeline(ctx, projectID, pipelineID); err != nil {
		return err
	}
	return nil
}

func (c *Controller) exportPipeline(ctx context.Context, projectID int64, pipelineID int64) <-chan error {
	errChan := make(chan error)

	go func() {
		defer close(errChan)

		cfg, err := c.projectConfig(projectID)
		if err != nil {
			errChan <- err
			return
		}

		opt := &gitlab.GetPipelineHierarchyOptions{
			FetchSections: cfg.Sections.Enabled,
		}

		phr := <-c.GitLab.GetPipelineHierarchy(ctx, projectID, pipelineID, opt)
		if err := phr.Error; err != nil {
			errChan <- fmt.Errorf("[controller.ExportPipeline/GetHierarchy] %w", err)
			return
		}
		ph := phr.PipelineHierarchy

		if err := clickhouse.InsertPipelineHierarchy(ctx, ph, c.ClickHouse); err != nil {
			errChan <- fmt.Errorf("[controller.ExportPipeline/InsertHierarchy] %w", err)
			return
		}

		pts := ph.GetAllTraces()
		if err := clickhouse.InsertTraces(ctx, pts, c.ClickHouse); err != nil {
			errChan <- fmt.Errorf("[controller.ExportPipeline/InsertTraces] %w", err)
			return
		}

		trs, err := c.GitLab.GetPipelineHierarchyTestReports(ctx, ph)
		if err != nil {
			errChan <- fmt.Errorf("[controller.ExportPipeline/GetTestRerports] %w", err)
			return
		}
		tss := []*models.PipelineTestSuite{}
		tcs := []*models.PipelineTestCase{}
		for _, tr := range trs {
			tss = append(tss, tr.TestSuites...)
			for _, ts := range tr.TestSuites {
				tcs = append(tcs, ts.TestCases...)
			}
		}
		if err = clickhouse.InsertTestReports(ctx, trs, c.ClickHouse); err != nil {
			errChan <- fmt.Errorf("[controller.ExportPipeline/InsertTestReports] %w", err)
			return
		}
		if err = clickhouse.InsertTestSuites(ctx, tss, c.ClickHouse); err != nil {
			errChan <- fmt.Errorf("[controller.ExportPipeline/InsertTestSuites] %w", err)
			return
		}
		if err = clickhouse.InsertTestCases(ctx, tcs, c.ClickHouse); err != nil {
			errChan <- fmt.Errorf("[controller.ExportPipeline/InsertTestCases] %w", err)
			return
		}
	}()

	return errChan
}
