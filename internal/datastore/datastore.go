package datastore

import (
	"context"
	"time"

	"github.com/cluttrdev/gitlab-exporter/pkg/models"
)

type DataStore interface {
	Initialize(context.Context) error
	CheckReadiness(context.Context) error

	InsertPipelines(context.Context, []*models.Pipeline) error
	InsertJobs(context.Context, []*models.Job) error
	InsertSections(context.Context, []*models.Section) error
	InsertBridges(context.Context, []*models.Bridge) error

	InsertPipelineHierarchy(context.Context, *models.PipelineHierarchy) error
	InsertTraces(context.Context, []models.Trace) error

	InsertTestReports(context.Context, []*models.PipelineTestReport) error
	InsertTestSuites(context.Context, []*models.PipelineTestSuite) error
	InsertTestCases(context.Context, []*models.PipelineTestCase) error

	QueryProjectPipelinesLatestUpdate(context.Context, int64) (map[int64]time.Time, error)
}
