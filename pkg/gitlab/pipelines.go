package gitlabclient

import (
	"context"
	"fmt"
	"time"

	gogitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

type ListProjectPipelineOptions struct {
	Page    int
	PerPage int

	Scope         *string
	Status        *string
	Source        *string
	Ref           *string
	SHA           *string
	YamlErrors    *bool
	Username      *string
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
	Name          *string
	OrderBy       *string
	Sort          *string
}

type ListProjectPipelinesResult struct {
	Pipeline *models.PipelineInfo
	Error    error
}

func (c *Client) ListProjectPipelines(ctx context.Context, projectID int64, opt *ListProjectPipelineOptions) <-chan ListProjectPipelinesResult {
	out := make(chan ListProjectPipelinesResult)

	go func() {
		defer close(out)

		gitlabOpts := &gogitlab.ListProjectPipelinesOptions{
			ListOptions: gogitlab.ListOptions{
				Page:    opt.Page,
				PerPage: opt.PerPage,
			},
			Scope:         opt.Scope,
			Status:        (*gogitlab.BuildStateValue)(opt.Status),
			Source:        opt.Source,
			Ref:           opt.Ref,
			SHA:           opt.SHA,
			YamlErrors:    opt.YamlErrors,
			Username:      opt.Username,
			UpdatedAfter:  opt.UpdatedAfter,
			UpdatedBefore: opt.UpdatedBefore,
			Name:          opt.Name,
			OrderBy:       opt.OrderBy,
			Sort:          opt.Sort,
		}

		for {
			c.RLock()
			ps, resp, err := c.client.Pipelines.ListProjectPipelines(int(projectID), gitlabOpts, gogitlab.WithContext(ctx))
			c.RUnlock()
			if err != nil {
				out <- ListProjectPipelinesResult{
					Error: err,
				}
				return
			}

			for _, pi := range ps {
				out <- ListProjectPipelinesResult{
					Pipeline: models.NewPipelineInfo(pi),
				}
			}

			if resp.NextPage == 0 {
				break
			}
			gitlabOpts.Page = resp.NextPage
		}

	}()

	return out
}

func (c *Client) GetPipeline(ctx context.Context, projectID int64, pipelineID int64) (*models.Pipeline, error) {
	p := new(models.Pipeline)
	if err := <-c.getPipeline(ctx, projectID, pipelineID, p); err != nil {
		return nil, err
	}
	return p, nil
}

type GetPipelineResult struct {
	Pipeline *models.Pipeline
	Error    error
}

func (c *Client) getPipeline(ctx context.Context, projectID int64, pipelineID int64, p *models.Pipeline) <-chan error {
	out := make(chan error)
	go func() {
		defer close(out)

		c.RLock()
		pipeline, _, err := c.client.Pipelines.GetPipeline(int(projectID), int(pipelineID), gogitlab.WithContext(ctx))
		c.RUnlock()
		if err != nil {
			out <- fmt.Errorf("[gitlab.Client.GetPipeline] %w", err)
			return
		}

		*p = *models.NewPipeline(pipeline)
	}()

	return out
}
