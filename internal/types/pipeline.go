package types

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

type PipelineReference struct {
	Id  int64
	Iid int64

	Project ProjectReference
}

type Pipeline struct {
	Id  int64
	Iid int64

	Project ProjectReference

	Name          string
	Ref           string
	RefPath       string
	Sha           string
	Source        string
	Status        string
	FailureReason string

	CommittedAt *time.Time
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time

	QueuedDuration time.Duration
	Duration       time.Duration
	Coverage       float64

	Warnings   bool
	YamlErrors bool

	Child               bool
	UpstreamPipeline    *PipelineReference
	DownstreamPipelines []*PipelineReference

	MergeRequest *MergeRequestReference

	User UserReference
}

func ConvertPipelineReference(pipeline PipelineReference) *typespb.PipelineReference {
	return &typespb.PipelineReference{
		Id:      pipeline.Id,
		Iid:     pipeline.Iid,
		Project: ConvertProjectReference(pipeline.Project),
	}
}

func ConvertPipeline(pipeline Pipeline) *typespb.Pipeline {
	pbPipeline := &typespb.Pipeline{
		Id:      pipeline.Id,
		Iid:     pipeline.Iid,
		Project: ConvertProjectReference(pipeline.Project),

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

		User: convertUserReference(pipeline.User),
	}

	if pipeline.UpstreamPipeline != nil {
		pbPipeline.UpstreamPipeline = ConvertPipelineReference(*pipeline.UpstreamPipeline)
	}
	for _, dp := range pipeline.DownstreamPipelines {
		if dp == nil {
			continue
		}
		pbPipeline.DownstreamPipelines = append(pbPipeline.DownstreamPipelines, ConvertPipelineReference(*dp))
	}

	if pipeline.MergeRequest != nil {
		pbPipeline.MergeRequest = ConvertMergeRequestReference(*pipeline.MergeRequest)
	}

	return pbPipeline
}
