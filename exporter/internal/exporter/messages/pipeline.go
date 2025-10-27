package messages

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewPipelineReference(pipeline types.PipelineReference) *typespb.PipelineReference {
	return &typespb.PipelineReference{
		Id:      pipeline.Id,
		Iid:     pipeline.Iid,
		Project: NewProjectReference(pipeline.Project),
	}
}

func NewPipeline(pipeline types.Pipeline) *typespb.Pipeline {
	pbPipeline := &typespb.Pipeline{
		Id:      pipeline.Id,
		Iid:     pipeline.Iid,
		Project: NewProjectReference(pipeline.Project),

		Name:    pipeline.Name,
		Ref:     pipeline.Ref,
		RefPath: pipeline.RefPath,
		Sha:     pipeline.Sha,
		Source:  pipeline.Source,
		Status:  pipeline.Status,

		Timestamps: &typespb.PipelineTimestamps{
			CommittedAt: timestamppb.New(valOrZero(pipeline.CommittedAt)),
			CreatedAt:   timestamppb.New(valOrZero(pipeline.CreatedAt)),
			UpdatedAt:   timestamppb.New(valOrZero(pipeline.UpdatedAt)),
			StartedAt:   timestamppb.New(valOrZero(pipeline.StartedAt)),
			FinishedAt:  timestamppb.New(valOrZero(pipeline.FinishedAt)),
		},

		QueuedDuration: durationpb.New(pipeline.QueuedDuration),
		Duration:       durationpb.New(pipeline.Duration),
		Coverage:       pipeline.Coverage,

		Warnings:   pipeline.Warnings,
		YamlErrors: pipeline.YamlErrors,

		Child: pipeline.Child,
		// UpstreamPipeline: nil,
		// DownstreamPipelines: nil,

		// MergeRequest: nil,

		User: NewUserReference(pipeline.User),
	}

	if pipeline.UpstreamPipeline != nil {
		pbPipeline.UpstreamPipeline = NewPipelineReference(*pipeline.UpstreamPipeline)
	}
	for _, dp := range pipeline.DownstreamPipelines {
		if dp == nil {
			continue
		}
		pbPipeline.DownstreamPipelines = append(pbPipeline.DownstreamPipelines, NewPipelineReference(*dp))
	}

	if pipeline.MergeRequest != nil {
		pbPipeline.MergeRequest = NewMergeRequestReference(*pipeline.MergeRequest)
	}

	return pbPipeline
}
