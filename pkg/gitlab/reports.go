package gitlabclient

import (
	"context"
	"fmt"

	gogitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

func (c *Client) GetPipelineReport(ctx context.Context, projectID int64, pipelineID int64) (*models.PipelineTestReport, error) {
	report, _, err := c.client.Pipelines.GetPipelineTestReport(int(projectID), int(pipelineID), gogitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("[gitlab.Client.GetPipelineTestReport] %w", err)
	}

	return models.NewPipelineTestReport(pipelineID, report), nil
}
