package models

import (
	gitlab "github.com/xanzy/go-gitlab"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

func ConvertJob(job *gitlab.Job) *typespb.Job {
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
		CreatedAt:      ConvertTime(job.CreatedAt),
		StartedAt:      ConvertTime(job.StartedAt),
		FinishedAt:     ConvertTime(job.FinishedAt),
		ErasedAt:       ConvertTime(job.ErasedAt),
		Duration:       ConvertDuration(job.Duration),
		QueuedDuration: ConvertDuration(job.QueuedDuration),
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

func ConvertBridge(bridge *gitlab.Bridge) *typespb.Bridge {
	// account for downstream pipeline creation failures
	downstreamPipeline := &typespb.PipelineInfo{
		CreatedAt: &timestamppb.Timestamp{},
		UpdatedAt: &timestamppb.Timestamp{},
	}
	if bridge.DownstreamPipeline != nil {
		downstreamPipeline = ConvertPipelineInfo(bridge.DownstreamPipeline)
	}
	return &typespb.Bridge{
		// Commit: ConvertCommit(bridge.Commit),
		Id:             int64(bridge.ID),
		Name:           bridge.Name,
		Pipeline:       ConvertPipelineInfo(&bridge.Pipeline),
		Ref:            bridge.Ref,
		CreatedAt:      ConvertTime(bridge.CreatedAt),
		StartedAt:      ConvertTime(bridge.StartedAt),
		FinishedAt:     ConvertTime(bridge.FinishedAt),
		ErasedAt:       ConvertTime(bridge.ErasedAt),
		Duration:       ConvertDuration(bridge.Duration),
		QueuedDuration: ConvertDuration(bridge.QueuedDuration),
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
