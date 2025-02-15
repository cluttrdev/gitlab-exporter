package types

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type MergeRequestReference struct {
	Id  int64
	Iid int64

	Project ProjectReference
}

type MergeRequest struct {
	Id      int64
	Iid     int64
	Project ProjectReference

	CreatedAt *time.Time
	UpdatedAt *time.Time
	MergedAt  *time.Time
	ClosedAt  *time.Time

	Name   string
	Title  string
	Labels []string

	State       string
	MergeStatus string
	MergeError  string

	SourceProjectId int64
	SourceBranch    string
	TargetProjectId int64
	TargetBranch    string

	DiffStats    MergeRequestDiffStats
	DiffRefs     MergeRequestDiffRefs
	Participants MergeRequestParticipants

	Approved  bool
	Conflicts bool
	Draft     bool
	Mergeable bool

	UserNotesCount int64

	Milestone *MilestoneReference
}

type MergeRequestParticipants struct {
	Author    UserReference
	Assignees []UserReference
	Reviewers []UserReference
	Approvers []UserReference
	MergeUser UserReference
}

type MergeRequestDiffRefs struct {
	BaseSha  string
	HeadSha  string
	StartSha string

	MergeCommitSha  string
	RebaseCommitSha string
}

type MergeRequestDiffStats struct {
	Additions   int64
	Changes     int64
	Deletions   int64
	FileCount   int64
	CommitCount int64
}

type MilestoneReference struct {
	Id      int64
	Iid     int64
	Project ProjectReference
}

func ConvertMergeRequestReference(mr MergeRequestReference) *typespb.MergeRequestReference {
	return &typespb.MergeRequestReference{
		Id:      mr.Id,
		Iid:     mr.Iid,
		Project: ConvertProjectReference(mr.Project),
	}
}

func ConvertMergeRequest(mr MergeRequest) *typespb.MergeRequest {
	pbMr := &typespb.MergeRequest{
		Id:      mr.Id,
		Iid:     mr.Iid,
		Project: ConvertProjectReference(mr.Project),

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
			Project: ConvertProjectReference(mr.Project),
		}
	}

	return pbMr
}

func convertMergeRequestParticipants(p MergeRequestParticipants) *typespb.MergeRequestParticipants {
	pbp := &typespb.MergeRequestParticipants{
		Author: convertUserReference(p.Author),
		// Assignees: nil,
		// Reviewers: nil,
		// Approvers: nil,
		MergeUser: convertUserReference(p.MergeUser),
	}

	if l := len(p.Assignees); l > 0 {
		assignees := make([]*typespb.UserReference, 0, l)
		for _, assignee := range p.Assignees {
			assignees = append(assignees, convertUserReference(assignee))
		}
		pbp.Assignees = assignees
	}
	if l := len(p.Reviewers); l > 0 {
		reviewers := make([]*typespb.UserReference, 0, l)
		for _, reviewer := range p.Reviewers {
			reviewers = append(reviewers, convertUserReference(reviewer))
		}
		pbp.Reviewers = reviewers
	}

	return pbp
}

type MergeRequestNoteEvent struct {
	Id           int64
	MergeRequest MergeRequestReference

	CreatedAt *time.Time
	UpdatedAt *time.Time

	Type     string
	System   bool
	Internal bool

	AuthorId int64

	Resolvable bool
	Resolved   bool
	ResolvedAt *time.Time
	ResolverId int64
}

func ConvertMergeRequestNoteEvent(event MergeRequestNoteEvent) *typespb.MergeRequestNoteEvent {
	return &typespb.MergeRequestNoteEvent{
		Id:           int64(event.Id),
		MergeRequest: ConvertMergeRequestReference(event.MergeRequest),

		CreatedAt:  timestamppb.New(valOrZero(event.CreatedAt)),
		UpdatedAt:  timestamppb.New(valOrZero(event.UpdatedAt)),
		ResolvedAt: timestamppb.New(valOrZero(event.ResolvedAt)),

		Type:     event.Type,
		System:   event.System,
		Internal: event.Internal,

		Author: convertUserReference(event.Author),

		Resolveable: event.Resolvable,
		Resolved:    event.Resolved,
		Resolver:    convertUserReference(event.Resolver),
	}
}

func convertUserReference(user UserReference) *typespb.UserReference {
	return &typespb.UserReference{
		Id:       user.Id,
		Username: user.Username,
		Name:     user.Name,
	}
}
