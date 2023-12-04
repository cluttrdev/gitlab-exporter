package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/cluttrdev/gitlab-exporter/pkg/models"
)

type ClickHouseDataStore struct {
	client *Client
}

func NewClickHouseDataStore(c *Client) *ClickHouseDataStore {
	return &ClickHouseDataStore{
		client: c,
	}
}

func (ds *ClickHouseDataStore) Initialize(ctx context.Context) error {
	if err := ds.client.CreateTables(ctx); err != nil {
		return fmt.Errorf("error creating tables: %w", err)
	}

	return nil
}

func (ds *ClickHouseDataStore) CheckReadiness(ctx context.Context) error {
	return ds.client.CheckReadiness(ctx)
}

func (ds *ClickHouseDataStore) InsertPipelines(ctx context.Context, pipelines []*models.Pipeline) error {
	return InsertPipelines(ctx, pipelines, ds.client)
}

func (ds *ClickHouseDataStore) InsertJobs(ctx context.Context, jobs []*models.Job) error {
	return InsertJobs(ctx, jobs, ds.client)
}

func (ds *ClickHouseDataStore) InsertSections(ctx context.Context, sections []*models.Section) error {
	return InsertSections(ctx, sections, ds.client)
}

func (ds *ClickHouseDataStore) InsertBridges(ctx context.Context, bridges []*models.Bridge) error {
	return InsertBridges(ctx, bridges, ds.client)
}

func (ds *ClickHouseDataStore) InsertPipelineHierarchy(ctx context.Context, hierarchy *models.PipelineHierarchy) error {
	return InsertPipelineHierarchy(ctx, hierarchy, ds.client)
}

func (ds *ClickHouseDataStore) InsertTraces(ctx context.Context, traces []models.Trace) error {
	return InsertTraces(ctx, traces, ds.client)
}

func (ds *ClickHouseDataStore) InsertTestReports(ctx context.Context, reports []*models.PipelineTestReport) error {
	return InsertTestReports(ctx, reports, ds.client)
}

func (ds *ClickHouseDataStore) InsertTestSuites(ctx context.Context, suites []*models.PipelineTestSuite) error {
	return InsertTestSuites(ctx, suites, ds.client)
}

func (ds *ClickHouseDataStore) InsertTestCases(ctx context.Context, cases []*models.PipelineTestCase) error {
	return InsertTestCases(ctx, cases, ds.client)
}

func (ds *ClickHouseDataStore) InsertJobMetrics(ctx context.Context, metrics []*models.JobMetric) error {
	return InsertJobMetrics(ctx, metrics, ds.client)
}

func (ds *ClickHouseDataStore) QueryProjectPipelinesLatestUpdate(ctx context.Context, projectID int64) (map[int64]time.Time, error) {
	return ds.client.QueryProjectPipelinesLatestUpdate(ctx, projectID)
}
