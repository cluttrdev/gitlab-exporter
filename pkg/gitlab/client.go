package gitlabclient

import (
	"context"
    "fmt"
    "log"

	"github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

type Client struct {
    client *gitlab.Client
}

type ClientConfig struct {
    URL string
    Token string
}

func NewGitLabClient(cfg ClientConfig) (*Client, error) {
    opts := []gitlab.ClientOptionFunc{
        gitlab.WithBaseURL(cfg.URL),
    }

    client, err := gitlab.NewOAuthClient(cfg.Token, opts...)
    if err != nil {
        return nil, err
    }

    return &Client{
        client: client,
    }, nil
}

func (c *Client) GetPipelineHierarchy(ctx context.Context, projectID int64, pipelineID int64) (*models.PipelineHierarchy, error) {
    pipeline, err := c.GetPipeline(ctx, int(projectID), int(pipelineID))
    if err != nil {
        return nil, fmt.Errorf("[gitlab.Client.GetPipelineFull] %w", err)
    }

    jobs, err := c.GetJobs(ctx, int(projectID), int(pipelineID))
    if err != nil {
        return nil, fmt.Errorf("[gitlab.Client.GetPipelineFull] %w", err)
    }

    sections := []*models.Section{}
    for _, job := range jobs {
        log.Printf("job %s\n", job.WebURL)
        s, err := c.GetSections(ctx, int(projectID), job.ID)
        if err != nil {
            return nil, fmt.Errorf("[gitlab.Client.GetPipelineFull] %w", err)
        }
        sections = append(sections, s...)
    }

    bridges, err := c.GetBridges(ctx, int(projectID), int(pipelineID))
    if err != nil {
        return nil, fmt.Errorf("[gitlab.Client.GetPipelineFull] %w", err)
    }

    dps := []*models.PipelineHierarchy{}
    for _, bridge := range bridges {
        if bridge.DownstreamPipeline == nil {
            continue
        }
        ph, err := c.GetPipelineHierarchy(
            ctx, bridge.DownstreamPipeline.ProjectID, bridge.DownstreamPipeline.ID,
        )
        if err != nil {
            return nil, fmt.Errorf("[gitlab.Client.GetPipelineFull] %w", err)
        }
        dps = append(dps, ph)
    }

    return &models.PipelineHierarchy{
        Pipeline: pipeline,
        Jobs: jobs,
        Sections: sections,
        Bridges: bridges,
        DownstreamPipelines: dps,
    }, nil
}
