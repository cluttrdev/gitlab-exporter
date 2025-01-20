package graphql

import (
	"context"
	"fmt"

	"github.com/cluttrdev/gitlab-exporter/internal/types"
)

type JobArtifactFields struct {
	Job      JobReferenceFields
	Pipeline PipelineReferenceFields
	Project  ProjectReferenceFields

	JobArtifactFieldsCore
}

func ConvertJobArtifact(jaf JobArtifactFields) (types.JobArtifact, error) {
	var (
		jobId, pipelineId, pipelineIid, projectId int64
		err                                       error
	)
	if jobId, err = parseJobId(valOrZero(jaf.Job.Id)); err != nil {
		return types.JobArtifact{}, fmt.Errorf("parse job id: %w", err)
	}
	if pipelineId, err = ParseId(jaf.Pipeline.Id, GlobalIdPipelinePrefix); err != nil {
		return types.JobArtifact{}, fmt.Errorf("parse pipeline id: %w", err)
	}
	if pipelineIid, err = ParseId(jaf.Pipeline.Iid, ""); err != nil {
		return types.JobArtifact{}, fmt.Errorf("parse pipeline iid: %w", err)
	}
	if projectId, err = ParseId(jaf.Project.Id, GlobalIdProjectPrefix); err != nil {
		return types.JobArtifact{}, fmt.Errorf("parse project id: %w", err)
	}

	return types.JobArtifact{
		Job: types.JobReference{
			Id:   jobId,
			Name: valOrZero(jaf.Name),
			Pipeline: types.PipelineReference{
				Id:  pipelineId,
				Iid: pipelineIid,
				Project: types.ProjectReference{
					Id:       projectId,
					FullPath: jaf.Project.FullPath,
				},
			},
		},
	}, nil
}

func (c *Client) GetProjectPipelineJobsArtifacts(ctx context.Context, projectPath string, pipelineIid string) ([]JobArtifactFields, error) {
	return c.getProjectPipelineJobsArtifacts(ctx, projectPath, pipelineIid, nil)
}

func (c *Client) getProjectPipelineJobsArtifacts(ctx context.Context, projectPath string, pipelineIid string, endCursor *string) ([]JobArtifactFields, error) {
	var jobArtifacts []JobArtifactFields

	for {
		resp, err := getProjectPipelineJobsArtifacts(
			ctx,
			c.client,
			projectPath,
			pipelineIid,
			endCursor,
		)
		if err != nil {
			return nil, err
		}

		project_ := resp.Project
		if project_ == nil {
			return nil, fmt.Errorf("project not found: %v", projectPath)
		}
		pipeline_ := project_.Pipeline
		if pipeline_ == nil {
			return nil, fmt.Errorf("project pipeline not found: %v (%v)", pipelineIid, projectPath)
		}

		for _, job_ := range pipeline_.Jobs.Nodes {
			if job_.Artifacts == nil {
				continue
			}

			for _, artifact_ := range job_.Artifacts.Nodes {
				jobArtifact := JobArtifactFields{
					Job:      job_.JobReferenceFields,
					Pipeline: pipeline_.PipelineReferenceFields,
					Project:  project_.ProjectReferenceFields,

					JobArtifactFieldsCore: artifact_.JobArtifactFieldsCore,
				}

				jobArtifacts = append(jobArtifacts, jobArtifact)
			}

			if job_.Artifacts.PageInfo.HasNextPage && job_.Id != nil {
				endCursor_ := job_.Artifacts.PageInfo.EndCursor
				jobArtifacts_, err := c.getProjectPipelineJobArtifacts(ctx, projectPath, pipelineIid, *job_.Id, endCursor_)
				if err != nil {
					return nil, err
				}
				jobArtifacts = append(jobArtifacts, jobArtifacts_...)
			}
		}

		if !pipeline_.Jobs.PageInfo.HasNextPage {
			break
		}

		endCursor = pipeline_.Jobs.PageInfo.EndCursor
	}

	return jobArtifacts, nil
}

func (c *Client) getProjectPipelineJobArtifacts(ctx context.Context, projectPath string, pipelineIid string, jobId string, endCursor *string) ([]JobArtifactFields, error) {
	var jobArtifacts []JobArtifactFields

	for {
		resp, err := getProjectPipelineJobArtifacts(
			ctx,
			c.client,
			projectPath,
			pipelineIid,
			jobId,
			endCursor,
		)
		if err != nil {
			return nil, err
		}

		project_ := resp.Project
		if project_ == nil {
			return nil, fmt.Errorf("project not found: %v", projectPath)
		}
		pipeline_ := project_.Pipeline
		if pipeline_ == nil {
			return nil, fmt.Errorf("project pipeline not found: %v (%v)", pipelineIid, projectPath)
		}
		job_ := pipeline_.Job
		if job_ == nil {
			return nil, fmt.Errorf("project pipeline job not found: %v [%v] (%v)", jobId, pipelineIid, projectPath)
		}

		if job_.Artifacts == nil {
			break
		}

		for _, artifact_ := range job_.Artifacts.Nodes {
			jobArtifact := JobArtifactFields{
				Job:      job_.JobReferenceFields,
				Pipeline: pipeline_.PipelineReferenceFields,
				Project:  project_.ProjectReferenceFields,

				JobArtifactFieldsCore: artifact_.JobArtifactFieldsCore,
			}

			jobArtifacts = append(jobArtifacts, jobArtifact)
		}

		if !job_.Artifacts.PageInfo.HasNextPage {
			break
		}

		endCursor = job_.Artifacts.PageInfo.EndCursor
	}

	return jobArtifacts, nil
}
