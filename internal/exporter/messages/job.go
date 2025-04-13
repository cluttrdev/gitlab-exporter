package messages

import (
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewJobReference(job types.JobReference) *typespb.JobReference {
	return &typespb.JobReference{
		Id:       job.Id,
		Name:     job.Name,
		Pipeline: NewPipelineReference(job.Pipeline),
	}
}

func NewJob(job types.Job) *typespb.Job {
	j := &typespb.Job{
		Id:       job.Id,
		Name:     job.Name,
		Pipeline: NewPipelineReference(job.Pipeline),

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
		j.DownstreamPipeline = NewPipelineReference(*job.DownstreamPipeline)
	}

	return j
}

func convertJobKind(kind types.JobKind) typespb.JobKind {
	switch kind {
	case types.JobKindBuild:
		return typespb.JobKind_JOBKIND_BUILD
	case types.JobKindBridge:
		return typespb.JobKind_JOBKIND_BRIDGE
	default:
		return typespb.JobKind_JOBKIND_UNSPECIFIED
	}
}

func convertJobProperties(props []types.JobLogProperty) []*typespb.JobProperty {
	pbProps := make([]*typespb.JobProperty, 0, len(props))
	for _, p := range props {
		pbProps = append(pbProps, &typespb.JobProperty{
			Name:  p.Name,
			Value: p.Value,
		})
	}
	return pbProps
}
