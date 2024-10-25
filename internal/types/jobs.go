package types

import (
	"time"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type JobKind string

const (
	JobKindUnkown JobKind = "unknown"
	JobKindBuild  JobKind = "build"
	JobKindBridge JobKind = "bridge"
)

type JobReference struct {
	Id   int64
	Name string

	Pipeline PipelineReference
}

type Job struct {
	Id       int64
	Pipeline PipelineReference

	Name          string
	Ref           string
	Status        string
	FailureReason string

	CreatedAt  *time.Time
	QueuedAt   *time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
	ErasedAt   *time.Time

	Stage string
	Tags  []string

	QueuedDuration time.Duration
	Duration       time.Duration
	Coverage       float64

	AllowFailure bool
	Manual       bool
	Retried      bool
	Retryable    bool

	Kind               JobKind
	DownstreamPipeline *PipelineReference

	RunnerId string
}

func ConvertJobReference(job JobReference) *typespb.JobReference {
	return &typespb.JobReference{
		Id:       job.Id,
		Name:     job.Name,
		Pipeline: ConvertPipelineReference(job.Pipeline),
	}
}

func ConvertJob(job Job) *typespb.Job {
	j := &typespb.Job{
		Id:       job.Id,
		Pipeline: ConvertPipelineReference(job.Pipeline),

		Name:          job.Name,
		Ref:           job.Ref,
		Status:        job.Status,
		FailureReason: job.FailureReason,

		Timestamps: &typespb.JobTimestamps{
			CreatedAt:  ConvertTime(job.CreatedAt),
			QueuedAt:   ConvertTime(job.QueuedAt),
			StartedAt:  ConvertTime(job.StartedAt),
			FinishedAt: ConvertTime(job.FinishedAt),
			ErasedAt:   ConvertTime(job.ErasedAt),
		},

		QueuedDuration: ConvertDuration(job.Duration.Seconds()),
		Duration:       ConvertDuration(job.Duration.Seconds()),
		Coverage:       job.Coverage,

		Stage: job.Stage,
		Tags:  job.Tags,

		AllowFailure: job.AllowFailure,
		Manual:       job.Manual,
		Retried:      job.Retried,
		Retryable:    job.Retryable,

		Kind: convertJobKind(job.Kind),
		// DownstreamPipeline: nil,

		Runner: &typespb.RunnerReference{
			Id: job.RunnerId,
		},
	}

	if job.DownstreamPipeline != nil {
		j.DownstreamPipeline = ConvertPipelineReference(*job.DownstreamPipeline)
	}

	return j
}

func convertJobKind(kind JobKind) typespb.JobKind {
	switch kind {
	case JobKindBuild:
		return typespb.JobKind_JOBKIND_BUILD
	case JobKindBridge:
		return typespb.JobKind_JOBKIND_BRIDGE
	default:
		return typespb.JobKind_JOBKIND_UNSPECIFIED
	}
}
