package gitlab

import (
	"context"

	_gitlab "github.com/xanzy/go-gitlab"

	pb "github.com/cluttrdev/gitlab-exporter/grpc/exporterpb"
	"github.com/cluttrdev/gitlab-exporter/internal/models"
)

type ListPipelineJobsResult struct {
	Job   *pb.Job
	Error error
}

func (c *Client) ListPipelineJobs(ctx context.Context, projectID int64, pipelineID int64) <-chan ListPipelineJobsResult {
	ch := make(chan ListPipelineJobsResult)

	go func() {
		defer close(ch)

		opts := &_gitlab.ListJobsOptions{
			ListOptions: _gitlab.ListOptions{
				PerPage: 100,
				Page:    1,
			},
			IncludeRetried: &[]bool{false}[0],
		}

		for {
			c.RLock()
			jobs, res, err := c.client.Jobs.ListPipelineJobs(int(projectID), int(pipelineID), opts, _gitlab.WithContext(ctx))
			c.RUnlock()
			if err != nil {
				ch <- ListPipelineJobsResult{
					Error: err,
				}
				return
			}

			for _, j := range jobs {
				ch <- ListPipelineJobsResult{
					Job: models.ConvertJob(j),
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

type ListPipelineBridgesResult struct {
	Bridge *pb.Bridge
	Error  error
}

func (c *Client) ListPipelineBridges(ctx context.Context, projectID int64, pipelineID int64) <-chan ListPipelineBridgesResult {
	ch := make(chan ListPipelineBridgesResult)

	go func() {
		defer close(ch)

		opts := &_gitlab.ListJobsOptions{
			ListOptions: _gitlab.ListOptions{
				PerPage: 100,
				Page:    1,
			},
		}

		for {
			c.RLock()
			bridges, res, err := c.client.Jobs.ListPipelineBridges(int(projectID), int(pipelineID), opts, _gitlab.WithContext(ctx))
			c.RUnlock()
			if err != nil {
				ch <- ListPipelineBridgesResult{
					Error: err,
				}
				return
			}

			for _, b := range bridges {
				ch <- ListPipelineBridgesResult{
					Bridge: models.ConvertBridge(b),
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
