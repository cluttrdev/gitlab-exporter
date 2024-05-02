package gitlab

import (
	"context"
	"time"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type ListProjectMergeRequestsOptions struct {
	gitlab.ListProjectMergeRequestsOptions

	Paginate bool
}

func (c *Client) ListProjectMergeRequests(ctx context.Context, id int64, opt ListProjectMergeRequestsOptions) ([]*typespb.MergeRequest, error) {
	var mergerequests []*typespb.MergeRequest

	options := []gitlab.RequestOptionFunc{
		gitlab.WithContext(ctx),
	}

	for {
		mrs, resp, err := c.client.MergeRequests.ListProjectMergeRequests(int(id), &opt.ListProjectMergeRequestsOptions, options...)
		if err != nil {
			return mergerequests, err
		}

		for _, mr := range mrs {
			mergerequests = append(mergerequests, ConvertMergeRequest(mr))
		}

		if !opt.Paginate {
			break
		}

		if opt.ListOptions.Pagination == "keyset" {
			if resp.NextLink == "" {
				break
			}
			options = []gitlab.RequestOptionFunc{
				gitlab.WithContext(ctx),
				gitlab.WithKeysetPaginationParameters(resp.NextLink),
			}
		} else {
			if resp.NextPage == 0 {
				break
			}
			opt.Page = resp.NextPage
		}
	}

	return mergerequests, nil
}

func ConvertMergeRequest(mr *gitlab.MergeRequest) *typespb.MergeRequest {
	return &typespb.MergeRequest{
		Id:        int64(mr.ID),
		Iid:       int64(mr.IID),
		ProjectId: int64(mr.ProjectID),

		CreatedAt: convertTime(mr.CreatedAt),
		UpdatedAt: convertTime(mr.UpdatedAt),
		MergedAt:  convertTime(mr.MergedAt),
		ClosedAt:  convertTime(mr.ClosedAt),

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

		Pipeline: convertPipelineInfo(mr.Pipeline),

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
		CreatedAt: convertTime(u.CreatedAt),
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
		CreatedAt: convertTime(u.CreatedAt),
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
		CreatedAt: convertTime(m.CreatedAt),
		UpdatedAt: convertTime(m.UpdatedAt),
		StartDate: convertTime((*time.Time)(m.StartDate)),
		DueDate:   convertTime((*time.Time)(m.DueDate)),
	}
}
