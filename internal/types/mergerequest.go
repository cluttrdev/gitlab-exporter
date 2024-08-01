package types

import (
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

func ConvertMergeRequest(mr *gitlab.MergeRequest) *typespb.MergeRequest {
	return &typespb.MergeRequest{
		Id:        int64(mr.ID),
		Iid:       int64(mr.IID),
		ProjectId: int64(mr.ProjectID),

		CreatedAt: ConvertTime(mr.CreatedAt),
		UpdatedAt: ConvertTime(mr.UpdatedAt),
		MergedAt:  ConvertTime(mr.MergedAt),
		ClosedAt:  ConvertTime(mr.ClosedAt),

		SourceProjectId: int64(mr.SourceProjectID),
		TargetProjectId: int64(mr.TargetProjectID),
		SourceBranch:    mr.SourceBranch,
		TargetBranch:    mr.TargetBranch,

		Title:               mr.Title,
		State:               mr.State,
		DetailedMergeStatus: mr.DetailedMergeStatus,
		Draft:               mr.Draft,
		HasConflicts:        mr.HasConflicts,
		MergeError:          mr.MergeError,

		DiffRefs: &typespb.MergeRequestDiffRefs{
			BaseSha:  mr.DiffRefs.BaseSha,
			HeadSha:  mr.DiffRefs.HeadSha,
			StartSha: mr.DiffRefs.StartSha,
		},

		Author:    convertBasicUser(mr.Author),
		Assignee:  convertBasicUser(mr.Assignee),
		Assignees: convertUsers(mr.Assignees),
		Reviewers: convertUsers(mr.Reviewers),
		MergeUser: convertBasicUser(mr.MergedBy),
		CloseUser: convertBasicUser(mr.ClosedBy),

		Labels: mr.Labels,

		Sha:             mr.SHA,
		MergeCommitSha:  mr.MergeCommitSHA,
		SquashCommitSha: mr.SquashCommitSHA,

		ChangesCount:   mr.ChangesCount,
		UserNotesCount: int64(mr.UserNotesCount),
		Upvotes:        int64(mr.Upvotes),
		Downvotes:      int64(mr.Downvotes),

		Pipeline: ConvertPipelineInfo(mr.Pipeline),

		Milestone: convertMilestone(mr.Milestone),

		WebUrl: mr.WebURL,
	}
}

func convertBasicUser(u *gitlab.BasicUser) *typespb.User {
	if u == nil {
		return nil
	}
	return &typespb.User{
		Id:        int64(u.ID),
		Username:  u.Username,
		Name:      u.Name,
		State:     u.State,
		CreatedAt: ConvertTime(u.CreatedAt),
	}
}

func convertUser(u *gitlab.User) *typespb.User {
	if u == nil {
		return nil
	}
	return &typespb.User{
		Id:        int64(u.ID),
		Username:  u.Username,
		Name:      u.Name,
		State:     u.State,
		CreatedAt: ConvertTime(u.CreatedAt),
	}
}

func convertUsers(us []*gitlab.BasicUser) []*typespb.User {
	users := make([]*typespb.User, 0, len(us))
	for _, u := range us {
		users = append(users, convertBasicUser(u))
	}
	return users
}

func convertMilestone(m *gitlab.Milestone) *typespb.Milestone {
	if m == nil {
		return nil
	}
	return &typespb.Milestone{
		Id:        int64(m.ID),
		Iid:       int64(m.IID),
		ProjectId: int64(m.ProjectID),
		GroupId:   int64(m.GroupID),
		CreatedAt: ConvertTime(m.CreatedAt),
		UpdatedAt: ConvertTime(m.UpdatedAt),
		StartDate: ConvertTime((*time.Time)(m.StartDate)),
		DueDate:   ConvertTime((*time.Time)(m.DueDate)),
	}
}