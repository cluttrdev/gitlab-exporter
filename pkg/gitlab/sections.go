package gitlabclient

import (
    "context"
    "fmt"

    "github.com/xanzy/go-gitlab"

    "github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

func (c *Client) GetSections(ctx context.Context, projectID int, jobID int64) ([]*models.Section, error) {
    job, _, err := c.client.Jobs.GetJob(projectID, int(jobID), gitlab.WithContext(ctx))
    if err != nil {
        return nil, fmt.Errorf("[gitlab.Client.GetSections] %w", err)
    }

    trace, _, err := c.client.Jobs.GetTraceFile(projectID, int(jobID), gitlab.WithContext(ctx))
    if err != nil {
        return nil, fmt.Errorf("[gitlab.Client.GetSections] %w", err)
    }

    sections, err := models.ParseSections(trace)
    if err != nil {
        return nil, fmt.Errorf("[gitlab.Client.GetSections] %w", err)
    }

    for secnum, section := range sections {
        section.ID = int64(job.ID * 1000 + secnum)
        section.Job.ID = int64(job.ID)
        section.Job.Name = job.Name
        section.Job.Status = job.Status
        section.Pipeline.ID = int64(job.Pipeline.ID)
        section.Pipeline.ProjectID = int64(job.Pipeline.ProjectID)
        section.Pipeline.Ref = job.Pipeline.Ref
        section.Pipeline.Sha = job.Pipeline.Sha
        section.Pipeline.Status = job.Pipeline.Status
    }

    return sections, nil
}

