package messages

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewMergeRequestReference(mr types.MergeRequestReference) *typespb.MergeRequestReference {
	return &typespb.MergeRequestReference{
		Id:      mr.Id,
		Iid:     mr.Iid,
		Project: NewProjectReference(mr.Project),
	}
}

func NewMergeRequest(mr types.MergeRequest) *typespb.MergeRequest {
	pbMr := &typespb.MergeRequest{
		Id:      mr.Id,
		Iid:     mr.Iid,
		Project: NewProjectReference(mr.Project),

		Timestamps: &typespb.MergeRequestTimestamps{
			CreatedAt: timestamppb.New(valOrZero(mr.CreatedAt)),
			UpdatedAt: timestamppb.New(valOrZero(mr.UpdatedAt)),
			MergedAt:  timestamppb.New(valOrZero(mr.MergedAt)),
			ClosedAt:  timestamppb.New(valOrZero(mr.ClosedAt)),
		},

		Name:   mr.Name,
		Title:  mr.Title,
		Labels: mr.Labels,

		State:       mr.State,
		MergeStatus: mr.MergeStatus,
		MergeError:  mr.MergeError,

		SourceProjectId: mr.SourceProjectId,
		SourceBranch:    mr.SourceBranch,
		TargetProjectId: mr.TargetProjectId,
		TargetBranch:    mr.TargetBranch,

		DiffStats: &typespb.MergeRequestDiffStats{
			Additions:   mr.DiffStats.Additions,
			Changes:     mr.DiffStats.Changes,
			Deletions:   mr.DiffStats.Deletions,
			FileCount:   mr.DiffStats.FileCount,
			CommitCount: mr.DiffStats.CommitCount,
		},

		DiffRefs: &typespb.MergeRequestDiffRefs{
			BaseSha:  mr.DiffRefs.BaseSha,
			HeadSha:  mr.DiffRefs.HeadSha,
			StartSha: mr.DiffRefs.StartSha,

			MergeCommitSha:  mr.DiffRefs.MergeCommitSha,
			RebaseCommitSha: mr.DiffRefs.RebaseCommitSha,
		},

		Participants: &typespb.MergeRequestParticipants{
			Author:    NewUserReference(mr.Participants.Author),
			Assignees: NewUserReferences(mr.Participants.Assignees),
			Reviewers: NewUserReferences(mr.Participants.Reviewers),
			Approvers: NewUserReferences(mr.Participants.Approvers),
			MergeUser: NewUserReference(mr.Participants.MergeUser),
		},

		Flags: &typespb.MergeRequestFlags{
			Approved:  mr.Approved,
			Conflicts: mr.Conflicts,
			Draft:     mr.Draft,
			Mergeable: mr.Mergeable,
		},

		// Milestone: nil,
	}

	if mr.Milestone != nil {
		pbMr.Milestone = &typespb.MilestoneReference{
			Id:      mr.Milestone.Id,
			Iid:     mr.Milestone.Iid,
			Project: NewProjectReference(mr.Project),
		}
	}

	return pbMr
}

func NewMergeRequestNoteEvent(event types.MergeRequestNoteEvent) *typespb.MergeRequestNoteEvent {
	return &typespb.MergeRequestNoteEvent{
		Id:           int64(event.Id),
		MergeRequest: NewMergeRequestReference(event.MergeRequest),

		CreatedAt:  timestamppb.New(valOrZero(event.CreatedAt)),
		UpdatedAt:  timestamppb.New(valOrZero(event.UpdatedAt)),
		ResolvedAt: timestamppb.New(valOrZero(event.ResolvedAt)),

		Type:     event.Type,
		System:   event.System,
		Internal: event.Internal,

		Author: NewUserReference(event.Author),

		Resolveable: event.Resolvable,
		Resolved:    event.Resolved,
		Resolver:    NewUserReference(event.Resolver),
	}
}
