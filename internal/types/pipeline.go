package types

import (
	"time"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
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

	Child            bool
	UpstreamPipeline *PipelineReference

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

		Name:   pipeline.Name,
		Ref:    pipeline.Ref,
		Sha:    pipeline.Sha,
		Source: pipeline.Source,
		Status: pipeline.Status,

		Timestamps: &typespb.PipelineTimestamps{
			CommittedAt: ConvertTime(pipeline.CommittedAt),
			CreatedAt:   ConvertTime(pipeline.CreatedAt),
			UpdatedAt:   ConvertTime(pipeline.UpdatedAt),
			StartedAt:   ConvertTime(pipeline.StartedAt),
			FinishedAt:  ConvertTime(pipeline.FinishedAt),
		},

		QueuedDuration: ConvertDuration(float64(pipeline.QueuedDuration)),
		Duration:       ConvertDuration(float64(pipeline.Duration)),
		Coverage:       pipeline.Coverage,

		Warnings:   pipeline.Warnings,
		YamlErrors: pipeline.YamlErrors,

		Child: pipeline.Child,
		// UpstreamPipeline: nil,

		// MergeRequest: nil,

		User: convertUserReference(pipeline.User),
	}

	if pipeline.UpstreamPipeline != nil {
		pbPipeline.UpstreamPipeline = ConvertPipelineReference(*pipeline.UpstreamPipeline)
	}

	if pipeline.MergeRequest != nil {
		pbPipeline.MergeRequest = ConvertMergeRequestReference(*pipeline.MergeRequest)
	}

	return pbPipeline
}
