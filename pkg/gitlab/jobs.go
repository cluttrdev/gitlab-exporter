package gitlabclient

import (
    "context"
    "fmt"

    "github.com/xanzy/go-gitlab"

    "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

func (c *Client) GetJobs(ctx context.Context, projectID int, pipelineID int) ([]*models.Job, error) {
    opts := &gitlab.ListJobsOptions{
        ListOptions: gitlab.ListOptions{
            PerPage: 100,
            Page: 1,
        },
        IncludeRetried: &[]bool{false}[0],
    }

    jobs := []*models.Job{}
    for {
        js, res, err := c.client.Jobs.ListPipelineJobs(projectID, pipelineID, opts, gitlab.WithContext(ctx))
        if err != nil {
            return nil, fmt.Errorf("[gitlab.Client.GetJobs] %w", err)
        }

        for _, j := range js {
            jobs = append(jobs, models.NewJob(j))
        }

        if res.NextPage == 0 {
            break
        }

        opts.Page = res.NextPage
    }

    return jobs, nil
}

func (c *Client) GetBridges(ctx context.Context, projectID int, pipelineID int) ([]*models.Bridge, error) {
    opts := &gitlab.ListJobsOptions{
        ListOptions: gitlab.ListOptions{
            PerPage: 100,
            Page: 1,
        },
    }

    bridges := []*models.Bridge{}
    for {
        bs, res, err := c.client.Jobs.ListPipelineBridges(projectID, pipelineID, opts, gitlab.WithContext(ctx))
        if err != nil {
            return nil, fmt.Errorf("[gitlab.Client.GetBridges] %w", err)
        }

        for _, b := range bs {
            bridges = append(bridges, models.NewBridge(b))
        }

        if res.NextPage == 0 {
            break
        }

        opts.Page = res.NextPage
    }

    return bridges, nil
}
