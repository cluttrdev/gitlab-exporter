package gitlabclient

import (
	"context"
	"fmt"
	"time"

	"github.com/xanzy/go-gitlab"

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

func (c *Client) ListProjectPipelines(ctx context.Context, projectID int64, opt *ListProjectPipelineOptions) ([]*models.PipelineInfo, error) {
	gitlabOpts := &gitlab.ListProjectPipelinesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    opt.Page,
			PerPage: opt.PerPage,
		},
		Scope:         opt.Scope,
		Status:        (*gitlab.BuildStateValue)(opt.Status),
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

	pipelines := []*models.PipelineInfo{}
	for {
		ps, resp, err := c.client.Pipelines.ListProjectPipelines(int(projectID), gitlabOpts, gitlab.WithContext(ctx))
		if err != nil {
			return nil, err
		}

		for _, pi := range ps {
			pipelines = append(pipelines, models.NewPipelineInfo(pi))
		}

		if resp.NextPage == 0 {
			break
		}
		gitlabOpts.Page = resp.NextPage
	}

	return pipelines, nil
}

func (c *Client) GetPipeline(ctx context.Context, projectID int64, pipelineID int64) (*models.Pipeline, error) {

	pipeline, _, err := c.client.Pipelines.GetPipeline(int(projectID), int(pipelineID), gitlab.WithContext(ctx))
	if err != nil {
		return nil, fmt.Errorf("[gitlab.Client.GetPipeline] %w", err)
	}

	return models.NewPipeline(pipeline), nil
}
