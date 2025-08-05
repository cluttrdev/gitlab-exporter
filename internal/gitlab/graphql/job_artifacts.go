package graphql

import (
	"context"
	"fmt"
	"log/slog"

	"go.cluttr.dev/gitlab-exporter/internal/types"
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
	var (
		jobArtifacts []JobArtifactFields

		data *getProjectPipelineJobsArtifactsResponse
		err  error
	)

	for {
		data, err = getProjectPipelineJobsArtifacts(
			ctx,
			c.client,
			projectPath,
			pipelineIid,
			endCursor,
		)
		err = handleError(err, "getProjectPipelineJobsArtifacts",
			slog.String("projectPath", projectPath),
			slog.String("pipelineIid", pipelineIid),
		)
		if err != nil {
			break
		}

		project_ := data.Project
		if project_ == nil {
			err = fmt.Errorf("project not found: %v", projectPath)
			break
		}
		pipeline_ := project_.Pipeline
		if pipeline_ == nil {
			err = fmt.Errorf("project pipeline not found: %v (%v)", pipelineIid, projectPath)
			break
		}

		if pipeline_.Jobs == nil {
			break
		}

	jobsLoop:
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
				jobArtifacts_, err_ := c.getProjectPipelineJobArtifacts(ctx, projectPath, pipelineIid, *job_.Id, endCursor_)
				if err_ != nil {
					err = err_
					break jobsLoop
				}
				jobArtifacts = append(jobArtifacts, jobArtifacts_...)
			}
		}
		if err != nil {
			break
		}

		if !pipeline_.Jobs.PageInfo.HasNextPage {
			break
		}

		endCursor = pipeline_.Jobs.PageInfo.EndCursor
	}

	return jobArtifacts, err
}

func (c *Client) getProjectPipelineJobArtifacts(ctx context.Context, projectPath string, pipelineIid string, jobId string, endCursor *string) ([]JobArtifactFields, error) {
	var (
		jobArtifacts []JobArtifactFields

		data *getProjectPipelineJobArtifactsResponse
		err  error
	)

	for {
		data, err = getProjectPipelineJobArtifacts(
			ctx,
			c.client,
			projectPath,
			pipelineIid,
			jobId,
			endCursor,
		)
		err = handleError(err, "getProjectPipelineJobArtifacts",
			slog.String("projectPath", projectPath),
			slog.String("pipelineIid", pipelineIid),
			slog.String("jobId", jobId),
		)
		if err != nil {
			break
		}

		project_ := data.Project
		if project_ == nil {
			err = fmt.Errorf("project not found: %v", projectPath)
			break
		}
		pipeline_ := project_.Pipeline
		if pipeline_ == nil {
			err = fmt.Errorf("project pipeline not found: %v (%v)", pipelineIid, projectPath)
			break
		}
		job_ := pipeline_.Job
		if job_ == nil {
			err = fmt.Errorf("project pipeline job not found: %v [%v] (%v)", jobId, pipelineIid, projectPath)
			break
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

	return jobArtifacts, err
}
