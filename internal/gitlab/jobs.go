package gitlab

import (
	"context"

	gitlab "github.com/xanzy/go-gitlab"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type ListPipelineJobsResult struct {
	Job   *typespb.Job
	Error error
}

func (c *Client) ListPipelineJobs(ctx context.Context, projectID int64, pipelineID int64) <-chan ListPipelineJobsResult {
	ch := make(chan ListPipelineJobsResult)

	go func() {
		defer close(ch)

		opts := &gitlab.ListJobsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    1,
			},
			IncludeRetried: &[]bool{false}[0],
		}

		for {
			c.RLock()
			jobs, res, err := c.client.Jobs.ListPipelineJobs(int(projectID), int(pipelineID), opts, gitlab.WithContext(ctx))
			c.RUnlock()
			if err != nil {
				ch <- ListPipelineJobsResult{
					Error: err,
				}
				return
			}

			for _, j := range jobs {
				ch <- ListPipelineJobsResult{
					Job: convertJob(j),
				}
			}

			if res.NextPage == 0 {
				break
			}

			opts.Page = res.NextPage
		}
	}()

	return ch
}

type ListPipelineBridgesResult struct {
	Bridge *typespb.Bridge
	Error  error
}

func (c *Client) ListPipelineBridges(ctx context.Context, projectID int64, pipelineID int64) <-chan ListPipelineBridgesResult {
	ch := make(chan ListPipelineBridgesResult)

	go func() {
		defer close(ch)

		opts := &gitlab.ListJobsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 100,
				Page:    1,
			},
		}

		for {
			c.RLock()
			bridges, res, err := c.client.Jobs.ListPipelineBridges(int(projectID), int(pipelineID), opts, gitlab.WithContext(ctx))
			c.RUnlock()
			if err != nil {
				ch <- ListPipelineBridgesResult{
					Error: err,
				}
				return
			}

			for _, b := range bridges {
				ch <- ListPipelineBridgesResult{
					Bridge: convertBridge(b),
				}
			}

			if res.NextPage == 0 {
				break
			}

			opts.Page = res.NextPage
		}
	}()

	return ch
}

func convertJob(job *gitlab.Job) *typespb.Job {
	return &typespb.Job{
		// Commit: ConvertCommit(job.Commit),
		Id:   int64(job.ID),
		Name: job.Name,
		Pipeline: &typespb.PipelineReference{
			Id:        int64(job.Pipeline.ID),
			ProjectId: int64(job.Pipeline.ProjectID),
			Ref:       job.Pipeline.Ref,
			Sha:       job.Pipeline.Sha,
			Status:    job.Pipeline.Status,
		},
		Ref:            job.Ref,
		CreatedAt:      convertTime(job.CreatedAt),
		StartedAt:      convertTime(job.StartedAt),
		FinishedAt:     convertTime(job.FinishedAt),
		ErasedAt:       convertTime(job.ErasedAt),
		Duration:       convertDuration(job.Duration),
		QueuedDuration: convertDuration(job.QueuedDuration),
		Coverage:       job.Coverage,
		// Artifacts: ?,
		// ArtifactsFile: ?,
		// Runner: ConvertRunner(job.Runner),
		Stage:         job.Stage,
		Status:        job.Status,
		AllowFailure:  job.AllowFailure,
		FailureReason: job.FailureReason,
		Tag:           job.Tag,
		WebUrl:        job.WebURL,
		TagList:       job.TagList,
		// Project: ConvertProject(job.Project),
		// User: ConvertUser(job.User),
	}
}

func convertBridge(bridge *gitlab.Bridge) *typespb.Bridge {
	// account for downstream pipeline creation failures
	downstreamPipeline := &typespb.PipelineInfo{
		CreatedAt: &timestamppb.Timestamp{},
		UpdatedAt: &timestamppb.Timestamp{},
	}
	if bridge.DownstreamPipeline != nil {
		downstreamPipeline = convertPipelineInfo(bridge.DownstreamPipeline)
	}
	return &typespb.Bridge{
		// Commit: ConvertCommit(bridge.Commit),
		Id:             int64(bridge.ID),
		Name:           bridge.Name,
		Pipeline:       convertPipelineInfo(&bridge.Pipeline),
		Ref:            bridge.Ref,
		CreatedAt:      convertTime(bridge.CreatedAt),
		StartedAt:      convertTime(bridge.StartedAt),
		FinishedAt:     convertTime(bridge.FinishedAt),
		ErasedAt:       convertTime(bridge.ErasedAt),
		Duration:       convertDuration(bridge.Duration),
		QueuedDuration: convertDuration(bridge.QueuedDuration),
		Coverage:       bridge.Coverage,
		Stage:          bridge.Stage,
		Status:         bridge.Status,
		AllowFailure:   bridge.AllowFailure,
		FailureReason:  bridge.FailureReason,
		Tag:            bridge.Tag,
		WebUrl:         bridge.WebURL,
		// User: ConvertUser(bridge.User),
		DownstreamPipeline: downstreamPipeline,
	}
}
