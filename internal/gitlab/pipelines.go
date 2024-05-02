package gitlab

import (
	"context"
	"fmt"
	"strconv"

	gitlab "github.com/xanzy/go-gitlab"

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
					Pipeline: convertPipelineInfo(pi),
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
	return convertPipeline(pipeline), nil
}

func convertPipelineInfo(pipeline *gitlab.PipelineInfo) *typespb.PipelineInfo {
	if pipeline == nil {
		return nil
	}
	return &typespb.PipelineInfo{
		Id:        int64(pipeline.ID),
		Iid:       int64(pipeline.IID),
		ProjectId: int64(pipeline.ProjectID),
		Status:    pipeline.Status,
		Source:    pipeline.Source,
		Ref:       pipeline.Ref,
		Sha:       pipeline.SHA,
		WebUrl:    pipeline.WebURL,
		CreatedAt: convertTime(pipeline.CreatedAt),
		UpdatedAt: convertTime(pipeline.UpdatedAt),
	}
}

func convertPipeline(pipeline *gitlab.Pipeline) *typespb.Pipeline {
	return &typespb.Pipeline{
		Id:             int64(pipeline.ID),
		Iid:            int64(pipeline.IID),
		ProjectId:      int64(pipeline.ProjectID),
		Status:         pipeline.Status,
		Source:         pipeline.Source,
		Ref:            pipeline.Ref,
		Sha:            pipeline.SHA,
		BeforeSha:      pipeline.BeforeSHA,
		Tag:            pipeline.Tag,
		YamlErrors:     pipeline.YamlErrors,
		CreatedAt:      convertTime(pipeline.CreatedAt),
		UpdatedAt:      convertTime(pipeline.UpdatedAt),
		StartedAt:      convertTime(pipeline.StartedAt),
		FinishedAt:     convertTime(pipeline.FinishedAt),
		CommittedAt:    convertTime(pipeline.CommittedAt),
		Duration:       convertDuration(float64(pipeline.Duration)),
		QueuedDuration: convertDuration(float64(pipeline.QueuedDuration)),
		Coverage:       convertCoverage(pipeline.Coverage),
		WebUrl:         pipeline.WebURL,
		User:           convertBasicUser(pipeline.User),
	}
}

func convertCoverage(coverage string) float64 {
	cov, err := strconv.ParseFloat(coverage, 64)
	if err != nil {
		cov = 0.0
	}
	return cov
}
