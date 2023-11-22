package gitlab

import (
	"context"
	"fmt"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/pkg/models"
)

func (c *Client) GetPipelineTestReport(ctx context.Context, projectID int64, pipelineID int64) (*models.PipelineTestReport, error) {
	c.RLock()
	report, _, err := c.client.Pipelines.GetPipelineTestReport(int(projectID), int(pipelineID), _gitlab.WithContext(ctx))
	c.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("[gitlab.Client.GetPipelineTestReport] %w", err)
	}

	return models.NewPipelineTestReport(pipelineID, report), nil
}

func (c *Client) GetPipelineHierarchyTestReports(ctx context.Context, ph *models.PipelineHierarchy) ([]*models.PipelineTestReport, error) {
	tr, err := c.GetPipelineTestReport(ctx, ph.Pipeline.ProjectID, ph.Pipeline.ID)
	if err != nil {
		return nil, fmt.Errorf("[gitlab.Client.GetPipelineHierarchyTestReports] %w", err)
	}

	reports := []*models.PipelineTestReport{tr}
	for _, dph := range ph.DownstreamPipelines {
		trs, err := c.GetPipelineHierarchyTestReports(ctx, dph)
		if err != nil {
			return nil, err
		}
		reports = append(reports, trs...)
	}

	return reports, nil
}
