package gitlab

import (
	"context"
	"fmt"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/types"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type ListOptions = gitlab.ListOptions
type ListProjectPipelinesOptions = gitlab.ListProjectPipelinesOptions

type ListProjectPipelinesResult struct {
	Pipeline *typespb.PipelineInfo
	Error    error
}

func (c *Client) ListProjectPipelines(ctx context.Context, projectID int64, opt ListProjectPipelinesOptions) <-chan ListProjectPipelinesResult {
	out := make(chan ListProjectPipelinesResult)

	go func() {
		defer close(out)

		for {
			c.RLock()
			ps, resp, err := c.client.Pipelines.ListProjectPipelines(int(projectID), &opt, gitlab.WithContext(ctx))
			c.RUnlock()
			if err != nil {
				out <- ListProjectPipelinesResult{
					Error: err,
				}
				return
			}

			for _, pi := range ps {
				out <- ListProjectPipelinesResult{
					Pipeline: types.ConvertPipelineInfo(pi),
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

func (c *Client) GetPipeline(ctx context.Context, projectID int64, pipelineID int64) (*typespb.Pipeline, error) {
	c.RLock()
	pipeline, _, err := c.client.Pipelines.GetPipeline(int(projectID), int(pipelineID), gitlab.WithContext(ctx))
	c.RUnlock()
	if err != nil {
		return nil, fmt.Errorf("error getting pipeline: %w", err)
	}
	return types.ConvertPipeline(pipeline), nil
}
