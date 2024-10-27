package graphql

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cluttrdev/gitlab-exporter/internal/types"
)

type JobFields struct {
	JobReferenceFields
	Pipeline PipelineReferenceFields
	Project  ProjectReferenceFields

	JobFieldsCore
	JobFieldsExtra
}

func ConvertJob(jf JobFields) (types.Job, error) {
	var (
		id, pipelineId, pipelineIid, projectId int64
		err                                    error
	)
	if id, err = parseJobId(valOrZero(jf.Id)); err != nil {
		return types.Job{}, fmt.Errorf("parse job id: %w", err)
	}
	if pipelineId, err = ParseId(jf.Pipeline.Id, GlobalIdPipelinePrefix); err != nil {
		return types.Job{}, fmt.Errorf("parse pipeline id: %w", err)
	}
	if pipelineIid, err = ParseId(jf.Pipeline.Iid, ""); err != nil {
		return types.Job{}, fmt.Errorf("parse pipeline iid: %w", err)
	}
	if projectId, err = ParseId(jf.Project.Id, GlobalIdProjectPrefix); err != nil {
		return types.Job{}, fmt.Errorf("parse job id: %w", err)
	}

	var stage string
	if jf.Stage != nil {
		stage = valOrZero(jf.Stage.Name)
	}

	job := types.Job{
		Id: id,
		Pipeline: types.PipelineReference{
			Id:  pipelineId,
			Iid: pipelineIid,
			Project: types.ProjectReference{
				Id:       projectId,
				FullPath: jf.Project.FullPath,
			},
		},

		Name:   valOrZero(jf.Name),
		Ref:    valOrZero(jf.RefName),
		Status: strings.ToLower(string(valOrZero(jf.Status))),

		CreatedAt:  &jf.CreatedAt,
		QueuedAt:   jf.QueuedAt,
		StartedAt:  jf.StartedAt,
		FinishedAt: jf.FinishedAt,
		ErasedAt:   jf.ErasedAt,

		Stage: stage,
		Tags:  jf.Tags,

		QueuedDuration: time.Duration(valOrZero(jf.QueuedDuration) * float64(time.Second)),
		Duration:       time.Duration(valOrZero(jf.Duration) * int(time.Second)),
		Coverage:       valOrZero(jf.Coverage),
		FailureReason:  valOrZero(jf.FailureMessage),

		AllowFailure: jf.AllowFailure,
		Manual:       valOrZero(jf.ManualJob),
		Retried:      valOrZero(jf.Retried),
		Retryable:    jf.Retryable,

		Kind: convertJobKind(jf.Kind),
		// DownstreamPipeline: nil

		RunnerId: valOrZero(jf.Runner).Id,
	}

	if jf.DownstreamPipeline != nil {
		downstreamId, _ := ParseId(jf.DownstreamPipeline.Id, GlobalIdPipelinePrefix)
		downstreamIid, _ := ParseId(jf.DownstreamPipeline.Iid, "")
		downstreamProjectId, _ := ParseId(jf.DownstreamPipeline.Project.Id, GlobalIdProjectPrefix)
		job.DownstreamPipeline = &types.PipelineReference{
			Id:  downstreamId,
			Iid: downstreamIid,
			Project: types.ProjectReference{
				Id:       downstreamProjectId,
				FullPath: jf.DownstreamPipeline.Project.FullPath,
			},
		}
	}

	return job, nil
}

func convertJobKind(kind CiJobKind) types.JobKind {
	switch kind {
	case CiJobKindBuild:
		return types.JobKindBuild
	case CiJobKindBridge:
		return types.JobKindBridge
	default:
		return types.JobKindUnkown
	}
}

func (c *Client) GetProjectsPipelinesJobs(ctx context.Context, ids []string, opts GetPipelinesOptions) ([]JobFields, error) {
	var jobs []JobFields

	pipelinesJobs, err := c.getProjectsPipelinesJobs(ctx, ids, getPipelinesOptions{
		GetPipelinesOptions: opts,
	})
	if err != nil {
		return nil, err
	}

	childPipelinesJobs, err := c.getProjectsPipelinesJobs(ctx, ids, getPipelinesOptions{
		GetPipelinesOptions: opts,
		Source:              ptr("parent_pipeline"),
	})
	if err != nil {
		return nil, err
	}

	downstreamPipelinesJobs, err := c.getProjectsPipelinesJobs(ctx, ids, getPipelinesOptions{
		GetPipelinesOptions: opts,
		Source:              ptr("pipeline"),
	})
	if err != nil {
		return nil, err
	}

	jobs = append(pipelinesJobs, append(childPipelinesJobs, downstreamPipelinesJobs...)...)

	return jobs, nil
}

func (c *Client) getProjectsPipelinesJobs(ctx context.Context, ids []string, opts getPipelinesOptions) ([]JobFields, error) {
	jobs_ := make(map[string]JobFields)

	jsCore, err := c.getProjectsPipelinesJobsPart(ctx, ids, getPipelinesJobsOptions{
		getPipelinesOptions: opts,

		core: true,
	})
	if err != nil {
		return nil, err
	}
	for _, j := range jsCore {
		if j.Id == nil {
			// TODO: what?
			continue
		}
		jobs_[*j.Id] = j
	}

	jsExtra, err := c.getProjectsPipelinesJobsPart(ctx, ids, getPipelinesJobsOptions{
		getPipelinesOptions: opts,

		extra: true,
	})
	if err != nil {
		return nil, err
	}
	for _, j := range jsExtra {
		if j.Id == nil {
			// TODO: what?
			continue
		}
		job_, ok := jobs_[*j.Id]
		if !ok {
			// TODO: what?
			continue
		}
		job_.JobFieldsExtra = j.JobFieldsExtra
		jobs_[*j.Id] = job_
	}

	jobs := make([]JobFields, 0, len(jobs_))
	for _, v := range jobs_ {
		jobs = append(jobs, v)
	}

	return jobs, nil
}

type getPipelinesJobsOptions struct {
	getPipelinesOptions

	core  bool
	extra bool
}

func (c *Client) getProjectsPipelinesJobsPart(ctx context.Context, ids []string, opts getPipelinesJobsOptions) ([]JobFields, error) {
	var jfs []JobFields

	for {
		resp, err := getProjectsPipelinesJobs(
			ctx,
			c.client,
			ids,
			opts.Source,
			opts.UpdatedAfter,
			opts.UpdatedBefore,
			opts.endCursor,
			//
			opts.core,
			opts.extra,
		)
		if err != nil {
			return nil, err
		}

		for _, project_ := range resp.Projects.Nodes {
			for _, pipeline_ := range project_.Pipelines.Nodes {
				for _, job_ := range pipeline_.Jobs.Nodes {
					jf := JobFields{
						JobReferenceFields: job_.JobReferenceFields,
						Pipeline:           pipeline_.PipelineReferenceFields,
						Project:            project_.ProjectReferenceFields,
					}
					if opts.core {
						jf.JobFieldsCore = job_.JobFieldsCore
					}
					if opts.extra {
						jf.JobFieldsExtra = job_.JobFieldsExtra
					}
					jfs = append(jfs, jf)
				}

				if pipeline_.Jobs.PageInfo.HasNextPage {
					opts_ := opts
					opts_.endCursor = pipeline_.Jobs.PageInfo.EndCursor
					jfs_, err := c.getProjectPipelineJobs(ctx, project_.FullPath, pipeline_.Iid, opts_)
					if err != nil {
						return nil, err
					}
					jfs = append(jfs, jfs_...)
				}
			}

			if project_.Pipelines.PageInfo.HasNextPage {
				opts_ := opts
				opts_.endCursor = project_.Pipelines.PageInfo.EndCursor
				jfs_, err := c.getProjectPipelinesJobs(ctx, project_.FullPath, opts_)
				if err != nil {
					return nil, err
				}
				jfs = append(jfs, jfs_...)
			}
		}

		if !resp.Projects.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = resp.Projects.PageInfo.EndCursor
	}

	return jfs, nil
}

func (c *Client) getProjectPipelinesJobs(ctx context.Context, projectPath string, opts getPipelinesJobsOptions) ([]JobFields, error) {
	var jfs []JobFields

	for {
		resp, err := getProjectPipelinesJobs(
			ctx,
			c.client,
			projectPath,
			opts.Source,
			opts.UpdatedAfter,
			opts.UpdatedBefore,
			opts.endCursor,
			//
			opts.core,
			opts.extra,
		)
		if err != nil {
			return nil, err
		}

		project_ := resp.Project
		if project_ == nil {
			return nil, fmt.Errorf("project not found: %v", projectPath)
		}

		for _, pipeline_ := range project_.Pipelines.Nodes {
			for _, job_ := range pipeline_.Jobs.Nodes {
				jf := JobFields{
					JobReferenceFields: job_.JobReferenceFields,
					Pipeline:           pipeline_.PipelineReferenceFields,
					Project:            project_.ProjectReferenceFields,
				}
				if opts.core {
					jf.JobFieldsCore = job_.JobFieldsCore
				}
				if opts.extra {
					jf.JobFieldsExtra = job_.JobFieldsExtra
				}
				jfs = append(jfs, jf)
			}

			if pipeline_.Jobs.PageInfo.HasNextPage {
				opts_ := opts
				opts_.endCursor = pipeline_.Jobs.PageInfo.EndCursor
				jfs_, err := c.getProjectPipelineJobs(ctx, projectPath, pipeline_.Iid, opts_)
				if err != nil {
					return nil, err
				}
				jfs = append(jfs, jfs_...)
			}
		}

		if !project_.Pipelines.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = project_.Pipelines.PageInfo.EndCursor
	}

	return jfs, nil
}

func (c *Client) getProjectPipelineJobs(ctx context.Context, projectPath string, pipelineIid string, opts getPipelinesJobsOptions) ([]JobFields, error) {
	var jfs []JobFields

	for {
		resp, err := getProjectPipelineJobs(
			ctx,
			c.client,
			projectPath,
			pipelineIid,
			opts.endCursor,
			//
			opts.core,
			opts.extra,
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
			jf := JobFields{
				JobReferenceFields: job_.JobReferenceFields,
				Pipeline:           pipeline_.PipelineReferenceFields,
				Project:            project_.ProjectReferenceFields,
			}
			if opts.core {
				jf.JobFieldsCore = job_.JobFieldsCore
			}
			if opts.extra {
				jf.JobFieldsExtra = job_.JobFieldsExtra
			}
			jfs = append(jfs, jf)
		}

		if !pipeline_.Jobs.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = pipeline_.Jobs.PageInfo.EndCursor
	}

	return jfs, nil
}

func parseJobId(s string) (int64, error) {
	var s_ string

	prefixes := []string{
		GlobalIdJobBuildPrefix,
		GlobalIdJobBridgePrefix,
		GlobalIdPrefix + "GenericCommitStatus/",
		GlobalIdPrefix + "CommitStatus/",
	}

	for _, prefix := range prefixes {
		s_ = strings.TrimPrefix(s, prefix)
		if len(s_) < len(s) {
			return strconv.ParseInt(s_, 10, 64)
		}
	}

	return strconv.ParseInt(s, 10, 64)
}
