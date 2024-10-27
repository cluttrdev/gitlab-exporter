package rest

import (
	"context"
	"fmt"

	_gitlab "github.com/xanzy/go-gitlab"
)

func ListProjectPipelines(ctx context.Context, glab *_gitlab.Client, pid int64, opt _gitlab.ListProjectPipelinesOptions, yield func(p []*_gitlab.PipelineInfo) bool) error {
	opt.ListOptions.Pagination = "keyset"
	if opt.ListOptions.OrderBy == "" {
		opt.ListOptions.OrderBy = "updated_at"
	}
	if opt.ListOptions.Sort == "" {
		opt.ListOptions.Sort = "desc"
	}

	options := []_gitlab.RequestOptionFunc{
		_gitlab.WithContext(ctx),
	}

	for {
		ps, resp, err := glab.Pipelines.ListProjectPipelines(int(pid), &opt, options...)
		if err != nil {
			return err
		}

		if !yield(ps) {
			break
		}

		if resp.NextLink == "" {
			break
		}

		options = []_gitlab.RequestOptionFunc{
			_gitlab.WithContext(ctx),
			_gitlab.WithKeysetPaginationParameters(resp.NextLink),
		}
	}

	return nil
}

func FetchProjectPipelines(ctx context.Context, glab *_gitlab.Client, pid int64, opt _gitlab.ListProjectPipelinesOptions) ([]*_gitlab.PipelineInfo, error) {
	var pipelines []*_gitlab.PipelineInfo

	err := ListProjectPipelines(ctx, glab, pid, opt, func(ps []*_gitlab.PipelineInfo) bool {
		pipelines = append(pipelines, ps...)
		return true
	})
	if err != nil {
		return nil, err
	}

	return pipelines, nil
}

func (c *Client) GetProjectPipeline(ctx context.Context, projectID int64, pipelineID int64) (*_gitlab.Pipeline, error) {
	pipeline, _, err := c.client.Pipelines.GetPipeline(int(projectID), int(pipelineID), _gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("error getting pipeline: %w", err)
	}
	return pipeline, nil
}
