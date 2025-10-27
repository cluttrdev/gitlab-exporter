package messages

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/internal/types"
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

		Participants: convertMergeRequestParticipants(mr.Participants),

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

func convertMergeRequestParticipants(p types.MergeRequestParticipants) *typespb.MergeRequestParticipants {
	pbp := &typespb.MergeRequestParticipants{
		Author: NewUserReference(p.Author),
		// Assignees: nil,
		// Reviewers: nil,
		// Approvers: nil,
		MergeUser: NewUserReference(p.MergeUser),
	}

	if l := len(p.Assignees); l > 0 {
		assignees := make([]*typespb.UserReference, 0, l)
		for _, assignee := range p.Assignees {
			assignees = append(assignees, NewUserReference(assignee))
		}
		pbp.Assignees = assignees
	}
	if l := len(p.Reviewers); l > 0 {
		reviewers := make([]*typespb.UserReference, 0, l)
		for _, reviewer := range p.Reviewers {
			reviewers = append(reviewers, NewUserReference(reviewer))
		}
		pbp.Reviewers = reviewers
	}

	return pbp
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
