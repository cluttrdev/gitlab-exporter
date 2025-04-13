package types

import (
	"time"
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

type MergeRequestNoteEvent struct {
	Id           int64
	MergeRequest MergeRequestReference

	CreatedAt *time.Time
	UpdatedAt *time.Time

	Type     string
	System   bool
	Internal bool

	Author UserReference

	Resolvable bool
	Resolved   bool
	ResolvedAt *time.Time
	Resolver   UserReference
}
