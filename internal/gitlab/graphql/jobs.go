package graphql

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/metaerr"
	"go.cluttr.dev/gitlab-exporter/internal/types"
)

type JobFields struct {
	JobReferenceFields
	Pipeline PipelineReferenceFields
	Project  ProjectReferenceFields

	JobFieldsCore
	JobFieldsExtra
}

func ConvertJobReference(job JobReferenceFields, pipeline PipelineReferenceFields, project ProjectReferenceFields) (types.JobReference, error) {
	var (
		id, pipelineId, pipelineIid, projectId int64
		err                                    error
	)
	if id, err = parseJobId(valOrZero(job.Id)); err != nil {
		return types.JobReference{}, fmt.Errorf("parse job id: %w", err)
	}
	if pipelineId, err = ParseId(pipeline.Id, GlobalIdPipelinePrefix); err != nil {
		return types.JobReference{}, fmt.Errorf("parse pipeline id: %w", err)
	}
	if pipelineIid, err = ParseId(pipeline.Iid, ""); err != nil {
		return types.JobReference{}, fmt.Errorf("parse pipeline iid: %w", err)
	}
	if projectId, err = ParseId(project.Id, GlobalIdProjectPrefix); err != nil {
		return types.JobReference{}, fmt.Errorf("parse project id: %w", err)
	}

	return types.JobReference{
		Id: id,
		Pipeline: types.PipelineReference{
			Id:  pipelineId,
			Iid: pipelineIid,
			Project: types.ProjectReference{
				Id:       projectId,
				FullPath: project.FullPath,
			},
		},
	}, nil
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
		return types.Job{}, fmt.Errorf("parse project id: %w", err)
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

		Name:    valOrZero(jf.Name),
		Ref:     valOrZero(jf.RefName),
		RefPath: valOrZero(jf.RefPath),
		Status:  strings.ToLower(string(valOrZero(jf.Status))),

		CreatedAt:  &jf.CreatedAt,
		QueuedAt:   jf.QueuedAt,
		StartedAt:  jf.StartedAt,
		FinishedAt: jf.FinishedAt,
		ErasedAt:   jf.ErasedAt,

		Stage: stage,
		Tags:  jf.Tags,

		ExitCode: int64(valOr(jf.ExitCode, -1)),

		QueuedDuration: time.Duration(valOrZero(jf.QueuedDuration) * float64(time.Second)),
		Duration:       time.Duration(valOrZero(jf.Duration) * int(time.Second)),
		Coverage:       valOrZero(jf.Coverage),
		// FailureReason:  valOrZero(jf.FailureMessage),

		AllowFailure: jf.AllowFailure,
		Manual:       valOrZero(jf.ManualJob),
		Retried:      valOrZero(jf.Retried),
		Retryable:    jf.Retryable,

		Kind: convertJobKind(jf.Kind),
		// DownstreamPipeline: nil

		RunnerId: valOrZero(jf.Runner).Id,
	}

	if valOrZero(jf.Status) == CiJobStatusFailed {
		job.FailureReason = mapFailureMessage(valOrZero(jf.FailureMessage))
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
	var (
		jfs []JobFields

		data *getProjectsPipelinesJobsResponse
		err  error
	)

outerLoop:
	for {
		data, err = getProjectsPipelinesJobs(
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
		err = handleError(err, "getProjectsPipelinesJobs",
			slog.Any("projectIds", ids),
			slog.String("updatedAfter", opts.UpdatedAfter.Format(time.RFC3339)),
			slog.String("updatedBefore", opts.UpdatedBefore.Format(time.RFC3339)),
		)
		if err != nil {
			break
		}

		if data.Projects == nil {
			break
		}
		for _, project_ := range data.Projects.Nodes {
			if project_.Pipelines == nil {
				continue
			}
			for _, pipeline_ := range project_.Pipelines.Nodes {
				if pipeline_.Jobs == nil {
					continue
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

				if pipeline_.Jobs.PageInfo.HasNextPage {
					opts_ := opts
					opts_.endCursor = pipeline_.Jobs.PageInfo.EndCursor
					jfs_, err_ := c.getProjectPipelineJobs(ctx, project_.FullPath, pipeline_.Iid, opts_)
					if err_ != nil {
						err = metaerr.WithMetadata(err_, "projectPath", project_.FullPath, "pipelineiid", pipeline_.Iid)
						break outerLoop
					}
					jfs = append(jfs, jfs_...)
				}
			}

			if project_.Pipelines.PageInfo.HasNextPage {
				opts_ := opts
				opts_.endCursor = project_.Pipelines.PageInfo.EndCursor
				jfs_, err_ := c.getProjectPipelinesJobs(ctx, project_.FullPath, opts_)
				if err_ != nil {
					err = metaerr.WithMetadata(err_, "projectPath", project_.FullPath)
					break outerLoop
				}
				jfs = append(jfs, jfs_...)
			}
		}

		if !data.Projects.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = data.Projects.PageInfo.EndCursor
	}

	return jfs, err
}

func (c *Client) getProjectPipelinesJobs(ctx context.Context, projectPath string, opts getPipelinesJobsOptions) ([]JobFields, error) {
	var (
		jfs []JobFields

		data *getProjectPipelinesJobsResponse
		err  error
	)

outerLoop:
	for {
		data, err = getProjectPipelinesJobs(
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
		err = handleError(err, "getProjectPipelinesJobs",
			slog.String("projectPath", projectPath),
			slog.String("updatedAfter", opts.UpdatedAfter.Format(time.RFC3339)),
			slog.String("updatedBefore", opts.UpdatedBefore.Format(time.RFC3339)),
		)
		if err != nil {
			break
		}

		project_ := data.Project
		if project_ == nil {
			err = fmt.Errorf("project not found: %v", projectPath)
			break
		}

		if project_.Pipelines == nil {
			break
		}
		for _, pipeline_ := range project_.Pipelines.Nodes {
			if pipeline_.Jobs == nil {
				continue
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

			if pipeline_.Jobs.PageInfo.HasNextPage {
				opts_ := opts
				opts_.endCursor = pipeline_.Jobs.PageInfo.EndCursor
				jfs_, err_ := c.getProjectPipelineJobs(ctx, projectPath, pipeline_.Iid, opts_)
				if err_ != nil {
					err = metaerr.WithMetadata(err_, "projectPath", projectPath, "pipelineIid", pipeline_.Iid)
					break outerLoop
				}
				jfs = append(jfs, jfs_...)
			}
		}

		if !project_.Pipelines.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = project_.Pipelines.PageInfo.EndCursor
	}

	return jfs, err
}

func (c *Client) getProjectPipelineJobs(ctx context.Context, projectPath string, pipelineIid string, opts getPipelinesJobsOptions) ([]JobFields, error) {
	var (
		jfs []JobFields

		data *getProjectPipelineJobsResponse
		err  error
	)

	for {
		data, err = getProjectPipelineJobs(
			ctx,
			c.client,
			projectPath,
			pipelineIid,
			opts.endCursor,
			//
			opts.core,
			opts.extra,
		)
		err = handleError(err, "getProjectPipelineJobs",
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

	return jfs, err
}

// mapFailureMessage returns the failure reason corresponding to a job failure message
// based on the mapping configured in:
// https://gitlab.com/gitlab-org/gitlab/-/blob/master/app/presenters/commit_status_presenter.rb
//
// see also: https://gitlab.com/gitlab-org/gitlab/-/blob/master/lib/gitlab/ci/status/build/failed.rb
func mapFailureMessage(msg string) string {
	switch {
	// 'There is an unknown failure, please try again' (f494f271)
	case strings.HasPrefix(msg, "There is an unknown failure"):
		return "unknown_failure"
	// 'There has been a script failure. Check the job log for more information' (f494f271,0bc9e0b4)
	// nil (dddc8a61)
	case msg == "":
		return "script_failure"
	// 'There has been an API failure, please try again' (f494f271)
	case strings.HasPrefix(msg, "There has been an API failure"):
		return "api_failure"
	// 'There has been a timeout failure or the job got stuck. Check your timeout limits or try again' (f494f271)
	case strings.HasPrefix(msg, "There has been a timeout failure or the job got stuck"):
		return "stuck_or_timeout_failure"
	// 'There has been a runner system failure, please try again' (f494f271)
	case strings.HasPrefix(msg, "There has been a runner system failure"):
		return "runner_system_failure"
	// 'There has been a missing dependency failure, check the job log for more information' (f494f271,0bc9e0b4)
	// 'There has been a missing dependency failure' (0bc9e0b4)
	case strings.HasPrefix(msg, "There has been a missing dependency failure"):
		return "missing_dependency_failure"
	// 'No runners support the requirements to run this job.' (ae20052d)
	case strings.HasPrefix(msg, "No runners support the requirements to run this job"):
		return "runner_unsupported"
	// 'Your runner is outdated, please upgrade your runner' (dddc8a61)
	case strings.HasPrefix(msg, "Your runner is outdated"):
		return "runner_unsupported"
	// 'Scheduled job could not be executed by some reason, please try again' (af4b85ce,2d328886)
	case strings.HasPrefix(msg, "Scheduled job could not be executed by some reason"):
		return "schedule_expired"
	// 'Delayed job could not be executed by some reason, please try again' (2d328886)
	case strings.HasPrefix(msg, "Delayed job could not be executed by some reason"):
		return "stale_schedule"
	// 'The script exceeded the maximum execution time set for the job' (4611a599)
	case strings.HasPrefix(msg, "The script exceeded the maximum execution time set for the job"):
		return "job_execution_timeout"
	// 'The job is archived and cannot be run' (83e479b8)
	case strings.HasPrefix(msg, "The job is archived and cannot be run"):
		return "archived_failure"
	// 'The job failed to complete prerequisite tasks' (00f0d356)
	case strings.HasPrefix(msg, "The job failed to complete prerequisite tasks"):
		return "unmet_prerequisites"
	// 'The scheduler failed to assign job to the runner, please try again or contact system administrator' (fd24f258)
	case strings.HasPrefix(msg, "The scheduler failed to assign job to the runner"):
		return "scheduler_failure"
	// 'There has been a structural integrity problem detected, please contact system administrator' (fd24f258)
	case strings.HasPrefix(msg, "There has been a structural integrity problem detected"):
		return "data_integrity_failure"
	// 'There has been an unknown job problem, please contact your system administrator with the job ID to review the logs' (5c7cd53b)
	case strings.HasPrefix(msg, "There has been an unknown job problem"):
		return "data_integrity_failure"
	// 'The deployment job is older than the previously succeeded deployment job, and therefore cannot be run' (97b7355e)
	case strings.HasPrefix(msg, "The deployment job is older than the previously succeeded deployment job"):
		return "forward_deployment_failure"
	// 'This job could not be executed because it would create infinitely looping pipelines' (99a15b8f)
	case strings.HasPrefix(msg, "This job could not be executed because it would create infinitely looping pipelines"):
		return "pipeline_loop_detected"
	// 'This job could not be executed because of insufficient permissions to track the upstream project.' (e75e2d47)
	case strings.HasPrefix(msg, "This job could not be executed because of insufficient permissions to track the upstream project"):
		return "insufficient_upstream_permissions"
	// 'This job could not be executed because upstream bridge project could not be found.' (e75e2d47)
	case strings.HasPrefix(msg, "This job could not be executed because upstream bridge project could not be found"):
		return "upstream_bridge_project_not_found"
	// 'This job could not be executed because downstream pipeline trigger definition is invalid' (7add0464)
	case strings.HasPrefix(msg, "This job could not be executed because downstream pipeline trigger definition is invalid"):
		return "invalid_bridge_trigger"
	// 'This job could not be executed because downstream bridge project could not be found' (7add0464)
	case strings.HasPrefix(msg, "This job could not be executed because downstream bridge project could not be found"):
		return "downstream_bridge_project_not_found"
	// 'The environment this job is deploying to is protected. Only users with permission may successfully run this job.' (e75e2d47)
	case strings.HasPrefix(msg, "The environment this job is deploying to is protected"):
		return "protected_environment_failure"
	// 'This job could not be executed because of insufficient permissions to create a downstream pipeline' (7add0464)
	case strings.HasPrefix(msg, "This job could not be executed because of insufficient permissions to create a downstream pipeline"):
		return "insufficient_bridge_permissions"
	// 'This job belongs to a child pipeline and cannot create further child pipelines' (7add0464)
	case strings.HasPrefix(msg, "This job belongs to a child pipeline and cannot create further child pipelines"):
		return "bridge_pipeline_is_child_pipeline"
	// 'The downstream pipeline could not be created' (7add0464)
	case strings.HasPrefix(msg, "The downstream pipeline could not be created"):
		return "downstream_pipeline_creation_failed"
	// 'The secrets provider can not be found' (cdb300b6)
	// 'The secrets provider can not be found. Check your CI/CD variables and try again.' (af06901d)
	case strings.HasPrefix(msg, "The secrets provider can not be found"):
		return "secrets_provider_not_found"
	// 'Maximum child pipeline depth has been reached' (14c31b6f)
	case strings.HasPrefix(msg, "Maximum child pipeline depth has been reached"):
		return "reached_max_descendant_pipelines_depth"
	// 'You reached the maximum depth of child pipelines' (0178ad61)
	case strings.HasPrefix(msg, "You reached the maximum depth of child pipelines"):
		return "reached_max_descendant_pipelines_depth"
	// 'The downstream pipeline tree is too large' (8467414d)
	case strings.HasPrefix(msg, "The downstream pipeline tree is too large"):
		return "reached_max_pipeline_hierarchy_size"
	// 'The job belongs to a deleted project' (0178ad61)
	case strings.HasPrefix(msg, "The job belongs to a deleted project"):
		return "project_deleted"
	// 'The user who created this job is blocked' (0178ad61)
	case strings.HasPrefix(msg, "The user who created this job is blocked"):
		return "user_blocked"
	// 'No more CI minutes available' (ba9ab364)
	case strings.HasPrefix(msg, "No more CI minutes available"):
		return "ci_quota_exceeded"
	// 'No more compute minutes available' (06762051)
	case strings.HasPrefix(msg, "No more compute minutes available"):
		return "ci_quota_exceeded"
	// 'No matching runner available' (3168331b)
	case strings.HasPrefix(msg, "No matching runner available"):
		return "no_matching_runner"
	// 'The job log size limit was reached' (204b33cf)
	case strings.HasPrefix(msg, "The job log size limit was reached"):
		return "trace_size_exceeded"
	// 'The CI/CD is disabled for this project' (1deeb21e)
	case strings.HasPrefix(msg, "The CI/CD is disabled for this project"):
		return "builds_disabled"
	// 'This job could not be executed because it would create an environment with an invalid parameter.' (8c4f6e12)
	case strings.HasPrefix(msg, "This job could not be executed because it would create an environment with an invalid parameter"):
		return "environment_creation_failure"
	// 'This deployment job was rejected.' (c4d5ac49)
	case strings.HasPrefix(msg, "This deployment job was rejected"):
		return "deployment_rejected"
	// 'This job could not be executed because group IP address restrictions are enabled, and the runners IP address is not in the allowed range.' (4c673725)
	case strings.HasPrefix(msg, "This job could not be executed because group IP address restrictions"):
		return "ip_restriction_failure"
	// 'The deployment job is older than the latest deployment, and therefore failed.' (7eaf7b52)
	case strings.HasPrefix(msg, "The deployment job is older than the latest deployment"):
		return "failed_outdated_deployment_job"
	// 'Too many downstream pipelines triggered in the last minute. Try again later.' (38dc9857)
	case strings.HasPrefix(msg, "Too many downstream pipelines triggered in the last minute"):
		return "reached_downstream_pipeline_trigger_rate_limit"
	}

	return "_unknown_"
}
