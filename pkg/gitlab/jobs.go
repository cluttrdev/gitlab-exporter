package gitlabclient

import (
	"context"
	"fmt"

	"github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

func (c *Client) GetPipelineJobs(ctx context.Context, projectID int64, pipelineID int64) ([]*models.Job, error) {
	jobs := []*models.Job{}
	for r := range c.ListPipelineJobs(ctx, projectID, pipelineID) {
		if r.Error != nil {
			return nil, fmt.Errorf("[gitlab.Client.GetJobs] %w", r.Error)
		}
		jobs = append(jobs, r.Job)
	}
	return jobs, nil
}

type ListPipelineJobsResult struct {
	Job   *models.Job
	Error error
}

func (c *Client) ListPipelineJobs(ctx context.Context, projectID int64, pipelineID int64) <-chan ListPipelineJobsResult {
	ch := make(chan ListPipelineJobsResult)

	go func() {
		defer close(ch)

		opts := &gitlab.ListJobsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    1,
			},
			IncludeRetried: &[]bool{false}[0],
		}

		for {
			c.RLock()
			jobs, res, err := c.client.Jobs.ListPipelineJobs(int(projectID), int(pipelineID), opts, gitlab.WithContext(ctx))
			c.RUnlock()
			if err != nil {
				ch <- ListPipelineJobsResult{
					Error: err,
				}
				return
			}

			for _, j := range jobs {
				ch <- ListPipelineJobsResult{
					Job: models.NewJob(j),
				}
			}

			if res.NextPage == 0 {
				break
			}

			opts.Page = res.NextPage
		}
	}()

	return ch
}

func (c *Client) GetPipelineBridges(ctx context.Context, projectID int64, pipelineID int64) ([]*models.Bridge, error) {
	bridges := []*models.Bridge{}
	for r := range c.ListPipelineBridges(ctx, projectID, pipelineID) {
		if r.Error != nil {
			return nil, fmt.Errorf("[gitlab.Client.GetBridges] %w", r.Error)
		}
		bridges = append(bridges, r.Bridge)
	}
	return bridges, nil
}

type ListPipelineBridgesResult struct {
	Bridge *models.Bridge
	Error  error
}

func (c *Client) ListPipelineBridges(ctx context.Context, projectID int64, pipelineID int64) <-chan ListPipelineBridgesResult {
	ch := make(chan ListPipelineBridgesResult)

	go func() {
		defer close(ch)

		opts := &gitlab.ListJobsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    1,
			},
		}

		for {
			c.RLock()
			bridges, res, err := c.client.Jobs.ListPipelineBridges(int(projectID), int(pipelineID), opts, gitlab.WithContext(ctx))
			c.RUnlock()
			if err != nil {
				ch <- ListPipelineBridgesResult{
					Error: err,
				}
				return
			}

			for _, b := range bridges {
				ch <- ListPipelineBridgesResult{
					Bridge: models.NewBridge(b),
				}
			}

			if res.NextPage == 0 {
				break
			}

			opts.Page = res.NextPage
		}
	}()

	return ch
}
