package gitlab

import (
	"context"
	"fmt"

	_gitlab "github.com/xanzy/go-gitlab"

	pb "github.com/cluttrdev/gitlab-exporter/grpc/exporterpb"
	"github.com/cluttrdev/gitlab-exporter/internal/models"
)

type ListOptions = _gitlab.ListOptions
type ListProjectPipelinesOptions = _gitlab.ListProjectPipelinesOptions

type ListProjectPipelinesResult struct {
	Pipeline *pb.PipelineInfo
	Error    error
}

func (c *Client) ListProjectPipelines(ctx context.Context, projectID int64, opt ListProjectPipelinesOptions) <-chan ListProjectPipelinesResult {
	out := make(chan ListProjectPipelinesResult)

	go func() {
		defer close(out)

		for {
			c.RLock()
			ps, resp, err := c.client.Pipelines.ListProjectPipelines(int(projectID), &opt, _gitlab.WithContext(ctx))
			c.RUnlock()
			if err != nil {
				out <- ListProjectPipelinesResult{
					Error: err,
				}
				return
			}

			for _, pi := range ps {
				out <- ListProjectPipelinesResult{
					Pipeline: models.ConvertPipelineInfo(pi),
				}
			}

			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
	}()

	return out
}

func (c *Client) GetPipeline(ctx context.Context, projectID int64, pipelineID int64) (*pb.Pipeline, error) {
	c.RLock()
	pipeline, _, err := c.client.Pipelines.GetPipeline(int(projectID), int(pipelineID), _gitlab.WithContext(ctx))
	c.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("error getting pipeline: %w", err)
	}
	return models.ConvertPipeline(pipeline), nil
}
