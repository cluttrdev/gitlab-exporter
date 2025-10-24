package messages

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewRunnerReference(runner types.RunnerReference) *typespb.RunnerReference {
	return &typespb.RunnerReference{
		Id:       runner.Id,
		ShortSha: runner.ShortSha,
	}
}

func NewRunner(runner types.Runner) *typespb.Runner {
	return &typespb.Runner{
		Id:          runner.Id,
		ShortSha:    runner.ShortSha,
		Description: runner.Description,

		RunnerType: convertRunnerType(runner.RunnerType),
		TagList:    runner.TagList,

		Status: convertRunnerStatus(runner.Status),

		Flags: &typespb.RunnerFlags{
			Locked: runner.Locked,
			Paused: runner.Paused,

			RunProtected: runner.RunProtected,
			RunUntagged:  runner.RunUntagged,
		},

		Timestamps: &typespb.RunnerTimestamps{
			CreatedAt:   timestamppb.New(valOrZero(runner.CreatedAt)),
			ContactedAt: timestamppb.New(valOrZero(runner.ContactedAt)),
		},

		CreatedBy: NewUserReference(runner.CreatedBy),
	}
}

func convertRunnerType(rt types.RunnerType) typespb.RunnerType {
	switch rt {
	case types.RunnerTypeInstance:
		return typespb.RunnerType_RUNNER_TYPE_INSTANCE
	case types.RunnerTypeGroup:
		return typespb.RunnerType_RUNNER_TYPE_GROUP
	case types.RunnerTypeProject:
		return typespb.RunnerType_RUNNER_TYPE_PROJECT
	case types.RunnerTypeUnknown:
		return typespb.RunnerType_RUNNER_TYPE_UNKNOWN
	default:
		return typespb.RunnerType_RUNNER_TYPE_UNSPECIFIED
	}
}

func convertRunnerStatus(rs types.RunnerStatus) typespb.RunnerStatus {
	switch rs {
	case types.RunnerStatusOnline:
		return typespb.RunnerStatus_RUNNER_STATUS_ONLINE
	case types.RunnerStatusOffline:
		return typespb.RunnerStatus_RUNNER_STATUS_OFFLINE
	case types.RunnerStatusStale:
		return typespb.RunnerStatus_RUNNER_STATUS_STALE
	case types.RunnerStatusNeverContacted:
		return typespb.RunnerStatus_RUNNER_STATUS_NEVER_CONTACTED
	case types.RunnerStatusUnknown:
		return typespb.RunnerStatus_RUNNER_STATUS_UNKNOWN
	default:
		return typespb.RunnerStatus_RUNNER_STATUS_UNSPECIFIED
	}
}
