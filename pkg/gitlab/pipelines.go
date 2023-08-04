package gitlabclient

import (
    "context"
    "fmt"

    "github.com/xanzy/go-gitlab"

    "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

func (c *Client) ListProjectPipelines(ctx context.Context, pid int) ([]*models.PipelineInfo, error) {
    opts := &gitlab.ListProjectPipelinesOptions{
        Sort: gitlab.String("desc"),
    }

    ps, _, err := c.client.Pipelines.ListProjectPipelines(pid, opts, gitlab.WithContext(ctx))
    if err != nil {
        return nil, err
    }

    pipelines := []*models.PipelineInfo{}
    for _, pi := range ps {
        pipelines = append(pipelines, models.NewPipelineInfo(pi))
    }

    return pipelines, nil
}

func (c *Client) GetPipeline(ctx context.Context, projectID int, pipelineID int) (*models.Pipeline, error) {


    pipeline, _, err := c.client.Pipelines.GetPipeline(projectID, pipelineID, gitlab.WithContext(ctx))
    if err != nil {
        return nil, fmt.Errorf("[gitlab.Client.GetPipeline] %w", err)
    }

    return models.NewPipeline(pipeline), nil
}
