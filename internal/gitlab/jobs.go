package gitlab

import (
	"context"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/types"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

func (c *Client) GetPipelineJobs(ctx context.Context, projectID int64, pipelineID int64) ([]*typespb.Job, error) {
	var jobs []*typespb.Job

	opt := &gitlab.ListJobsOptions{
		ListOptions: gitlab.ListOptions{
			Pagination: "keyset",
			PerPage:    100,
			OrderBy:    "updated_at",
			Sort:       "desc",
		},
		IncludeRetried: gitlab.Ptr(false),
	}

	options := []gitlab.RequestOptionFunc{
		gitlab.WithContext(ctx),
	}

	for {
		js, resp, err := c.client.Jobs.ListPipelineJobs(int(projectID), int(pipelineID), opt, options...)
		if err != nil {
			return nil, err
		}

		for _, j := range js {
			jobs = append(jobs, types.ConvertJob(j))
		}

		if resp.NextLink == "" {
			break
		}

		options = []gitlab.RequestOptionFunc{
			gitlab.WithContext(ctx),
			gitlab.WithKeysetPaginationParameters(resp.NextLink),
		}
	}

	return jobs, nil
}

func (c *Client) GetPipelineBridges(ctx context.Context, projectID int64, pipelineID int64) ([]*typespb.Bridge, error) {
	var bridges []*typespb.Bridge

	opts := &gitlab.ListJobsOptions{
		ListOptions: gitlab.ListOptions{
			Pagination: "keyset",
			PerPage:    100,
			OrderBy:    "updated_at",
			Sort:       "desc",
		},
		IncludeRetried: gitlab.Ptr(false),
	}

	options := []gitlab.RequestOptionFunc{
		gitlab.WithContext(ctx),
	}

	for {
		bs, resp, err := c.client.Jobs.ListPipelineBridges(int(projectID), int(pipelineID), opts, options...)
		if err != nil {
			return nil, err
		}

		for _, b := range bs {
			bridges = append(bridges, types.ConvertBridge(b))
		}

		if resp.NextLink == "" {
			break
		}

		options = []gitlab.RequestOptionFunc{
			gitlab.WithContext(ctx),
			gitlab.WithKeysetPaginationParameters(resp.NextLink),
		}
	}

	return bridges, nil
}
