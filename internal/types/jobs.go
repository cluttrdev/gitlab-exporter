package types

import (
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
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
	Name     string
	Pipeline PipelineReference

	Ref           string
	RefPath       string
	Status        string
	FailureReason string

	CreatedAt  *time.Time
	QueuedAt   *time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
	ErasedAt   *time.Time

	QueuedDuration time.Duration
	Duration       time.Duration
	Coverage       float64

	Stage      string
	Tags       []string
	Properties []JobLogProperty

	AllowFailure bool
	Manual       bool
	Retried      bool
	Retryable    bool

	Kind               JobKind
	DownstreamPipeline *PipelineReference

	RunnerId string
}

type JobLogProperty struct {
	Name  string
	Value string
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
		Name:     job.Name,
		Pipeline: ConvertPipelineReference(job.Pipeline),

		Ref:           job.Ref,
		RefPath:       job.RefPath,
		Status:        job.Status,
		FailureReason: job.FailureReason,

		Timestamps: &typespb.JobTimestamps{
			CreatedAt:  timestamppb.New(valOrZero(job.CreatedAt)),
			QueuedAt:   timestamppb.New(valOrZero(job.QueuedAt)),
			StartedAt:  timestamppb.New(valOrZero(job.StartedAt)),
			FinishedAt: timestamppb.New(valOrZero(job.FinishedAt)),
			ErasedAt:   timestamppb.New(valOrZero(job.ErasedAt)),
		},

		QueuedDuration: durationpb.New(job.QueuedDuration),
		Duration:       durationpb.New(job.Duration),
		Coverage:       job.Coverage,

		Stage:      job.Stage,
		Tags:       job.Tags,
		Properties: convertJobProperties(job.Properties),

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

func convertJobProperties(props []JobLogProperty) []*typespb.JobProperty {
	pbProps := make([]*typespb.JobProperty, 0, len(props))
	for _, p := range props {
		pbProps = append(pbProps, &typespb.JobProperty{
			Name:  p.Name,
			Value: p.Value,
		})
	}
	return pbProps
}
