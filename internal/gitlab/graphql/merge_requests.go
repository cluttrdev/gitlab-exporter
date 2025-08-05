package graphql

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/metaerr"
	"go.cluttr.dev/gitlab-exporter/internal/types"
)

type MergeRequestFields struct {
	MergeRequestReferenceFields
	Project ProjectReferenceFields

	MergeRequestFieldsCore
	MergeRequestFieldsExtra
	MergeRequestFieldsParticipants
}

func ConvertMergeRequest(mrf MergeRequestFields) (types.MergeRequest, error) {
	var (
		id, iid, projectId int64
		err                error
	)
	if id, err = ParseId(mrf.Id, GlobalIdMergeRequestPrefix); err != nil {
		return types.MergeRequest{}, fmt.Errorf("parse merge request id: %w", err)
	}
	if iid, err = ParseId(mrf.Iid, ""); err != nil {
		return types.MergeRequest{}, fmt.Errorf("parse merge request iid: %w", err)
	}
	if projectId, err = ParseId(mrf.Project.Id, GlobalIdProjectPrefix); err != nil {
		return types.MergeRequest{}, fmt.Errorf("parse project id: %w", err)
	}

	var labels []string
	for _, label := range mrf.Labels.Nodes {
		labels = append(labels, label.Title)
	}

	mr := types.MergeRequest{
		Id:  id,
		Iid: iid,
		Project: types.ProjectReference{
			Id:       projectId,
			FullPath: mrf.Project.FullPath,
		},

		CreatedAt: &mrf.CreatedAt,
		UpdatedAt: &mrf.UpdatedAt,
		MergedAt:  mrf.MergedAt,
		ClosedAt:  mrf.ClosedAt,

		Name:   valOrZero(mrf.Name),
		Title:  mrf.Title,
		Labels: labels,

		State:       string(mrf.State),
		MergeStatus: strings.ToLower(string(valOrZero(mrf.DetailedMergeStatus))),
		MergeError:  valOrZero(mrf.MergeError),

		SourceProjectId: int64(valOrZero(mrf.SourceProjectId)),
		SourceBranch:    mrf.SourceBranch,
		TargetProjectId: int64(mrf.TargetProjectId),
		TargetBranch:    mrf.TargetBranch,

		DiffStats: types.MergeRequestDiffStats{
			// ...
			CommitCount: int64(valOrZero(mrf.CommitCount)),
		},
		DiffRefs: types.MergeRequestDiffRefs{
			// ...
			MergeCommitSha:  valOrZero(mrf.MergeCommitSha),
			RebaseCommitSha: valOrZero(mrf.RebaseCommitSha),
		},

		Participants: types.MergeRequestParticipants{},

		Approved:  mrf.Approved,
		Conflicts: mrf.Conflicts,
		Draft:     mrf.Draft,
		Mergeable: mrf.Mergeable,

		UserNotesCount: int64(valOrZero(mrf.UserNotesCount)),

		// Milestone: nil,

	}

	// DiffStats
	if mrf.DiffStatsSummary != nil {
		mr.DiffStats.Additions = int64(mrf.DiffStatsSummary.Additions)
		mr.DiffStats.Changes = int64(mrf.DiffStatsSummary.Changes)
		mr.DiffStats.Deletions = int64(mrf.DiffStatsSummary.Deletions)
		mr.DiffStats.FileCount = int64(mrf.DiffStatsSummary.FileCount)
	}

	// DiffRefs
	if mrf.DiffRefs != nil {
		mr.DiffRefs.BaseSha = valOrZero(mrf.DiffRefs.BaseSha)
		mr.DiffRefs.HeadSha = mrf.DiffRefs.HeadSha
		mr.DiffRefs.StartSha = mrf.DiffRefs.StartSha
	}

	// Participants
	if mrf.Author != nil {
		author, err := convertUserReference(mrf.Author)
		if err != nil {
			return types.MergeRequest{}, fmt.Errorf("convert author reference: %w", err)
		}
		mr.Participants.Author = author
	}
	for _, assignee := range valOrZero(mrf.Assignees).Nodes {
		if assignee == nil {
			continue
		}
		assignee, err := convertUserReference(assignee)
		if err != nil {
			return types.MergeRequest{}, fmt.Errorf("convert assignee reference: %w", err)
		}
		mr.Participants.Assignees = append(mr.Participants.Assignees, assignee)
	}
	for _, reviewer := range valOrZero(mrf.Reviewers).Nodes {
		if reviewer == nil {
			continue
		}
		reviewer, err := convertUserReference(reviewer)
		if err != nil {
			return types.MergeRequest{}, fmt.Errorf("convert reviewer reference: %w", err)
		}
		mr.Participants.Reviewers = append(mr.Participants.Reviewers, reviewer)
	}
	for _, approver := range valOrZero(mrf.ApprovedBy).Nodes {
		if approver == nil {
			continue
		}
		approver, err := convertUserReference(approver)
		if err != nil {
			return types.MergeRequest{}, fmt.Errorf("convert approver reference: %w", err)
		}
		mr.Participants.Approvers = append(mr.Participants.Approvers, approver)
	}
	if mrf.MergeUser != nil {
		mergeUser, err := convertUserReference(mrf.MergeUser)
		if err != nil {
			return types.MergeRequest{}, fmt.Errorf("convert merge user reference: %w", err)
		}
		mr.Participants.MergeUser = mergeUser
	}

	// Milestone
	if mrf.Milestone != nil {
		var (
			milestoneId, milestoneIid, milestoneProjectId int64
			err                                           error
		)
		if milestoneId, err = ParseId(mrf.Milestone.Id, GlobalIdMilestonePrefix); err != nil {
			return types.MergeRequest{}, fmt.Errorf("parse milestone id: %w", err)
		}
		if milestoneIid, err = ParseId(mrf.Milestone.Iid, ""); err != nil {
			return types.MergeRequest{}, fmt.Errorf("parse milestone iid: %w", err)
		}
		mr.Milestone = &types.MilestoneReference{
			Id:  milestoneId,
			Iid: milestoneIid,
		}
		if mrf.Milestone.Project != nil {
			if milestoneProjectId, err = ParseId(mrf.Milestone.Project.Id, GlobalIdProjectPrefix); err != nil {
				return types.MergeRequest{}, fmt.Errorf("parse project id: %w", err)
			}
			mr.Milestone.Project = types.ProjectReference{
				Id:       milestoneProjectId,
				FullPath: mrf.Milestone.Project.FullPath,
			}
		}
	}

	return mr, nil
}

func convertUserReference(user UserReferenceFields) (types.UserReference, error) {
	var (
		id  int64
		err error
	)
	if id, err = ParseId(user.GetId(), GlobalIdUserPrefix); err != nil {
		return types.UserReference{}, fmt.Errorf("parse user id: %w", err)
	}

	return types.UserReference{
		Id:       id,
		Username: user.GetUsername(),
		Name:     user.GetName(),
	}, nil
}

type GetMergeRequestsOptions struct {
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
}

func (c *Client) GetProjectsMergeRequests(ctx context.Context, ids []string, opts GetMergeRequestsOptions) ([]MergeRequestFields, error) {
	mergeRequestsMap := make(map[string]MergeRequestFields)

	mrsCore, err := c.getProjectsMergeRequests(ctx, ids, getMergeRequestsOptions{
		GetMergeRequestsOptions: opts,
		includeCore:             true,
	})
	if err != nil {
		return nil, err
	}
	for _, mr := range mrsCore {
		mergeRequestsMap[mr.Id] = mr
	}

	mrsExtra, err := c.getProjectsMergeRequests(ctx, ids, getMergeRequestsOptions{
		GetMergeRequestsOptions: opts,
		includeExtra:            true,
	})
	if err != nil {
		return nil, err
	}
	for _, mr := range mrsExtra {
		mr_, ok := mergeRequestsMap[mr.Id]
		if !ok {
			// TODO: what?
			continue
		}
		mr_.MergeRequestFieldsExtra = mr.MergeRequestFieldsExtra
		mergeRequestsMap[mr.Id] = mr_
	}

	mrsParticipants, err := c.getProjectsMergeRequests(ctx, ids, getMergeRequestsOptions{
		GetMergeRequestsOptions: opts,
		includeParticipants:     true,
	})
	if err != nil {
		return nil, err
	}
	for _, mr := range mrsParticipants {
		mr_, ok := mergeRequestsMap[mr.Id]
		if !ok {
			// TODO: what?
			continue
		}
		mr_.MergeRequestFieldsParticipants = mr.MergeRequestFieldsParticipants
		mergeRequestsMap[mr.Id] = mr_
	}

	mergeRequests := make([]MergeRequestFields, 0, len(mergeRequestsMap))
	for _, v := range mergeRequestsMap {
		mergeRequests = append(mergeRequests, v)
	}
	return mergeRequests, nil
}

type getMergeRequestsOptions struct {
	GetMergeRequestsOptions

	endCursor *string

	includeCore         bool
	includeExtra        bool
	includeParticipants bool
}

func (c *Client) getProjectsMergeRequests(ctx context.Context, ids []string, opts getMergeRequestsOptions) ([]MergeRequestFields, error) {
	var (
		mergeRequests []MergeRequestFields

		data *getProjectsMergeRequestsResponse
		err  error
	)

outerLoop:
	for {
		data, err = getProjectsMergeRequests(
			ctx,
			c.client,
			ids,
			opts.UpdatedAfter,
			opts.UpdatedBefore,
			opts.endCursor,
			//
			&opts.includeCore,
			&opts.includeExtra,
			&opts.includeParticipants,
		)
		err = handleError(err, "getProjectsMergeRequests",
			slog.Any("projectIds", ids),
			slog.String("updatedAfter", opts.UpdatedAfter.Format(time.RFC3339)),
			slog.String("updatedBefore", opts.UpdatedBefore.Format(time.RFC3339)),
		)
		if err != nil {
			break
		}

		for _, project_ := range data.Projects.Nodes {
			if project_.MergeRequests == nil {
				continue
			}
			for _, mergeRequest_ := range project_.MergeRequests.Nodes {
				mergeRequest := MergeRequestFields{
					MergeRequestReferenceFields: mergeRequest_.MergeRequestReferenceFields,
					Project:                     project_.ProjectReferenceFields,
				}
				if opts.includeCore {
					mergeRequest.MergeRequestFieldsCore = mergeRequest_.MergeRequestFieldsCore
				}
				if opts.includeExtra {
					mergeRequest.MergeRequestFieldsExtra = mergeRequest_.MergeRequestFieldsExtra
				}
				if opts.includeParticipants {
					mergeRequest.MergeRequestFieldsParticipants = mergeRequest_.MergeRequestFieldsParticipants
				}
				mergeRequests = append(mergeRequests, mergeRequest)
			}

			if project_.MergeRequests.PageInfo.HasNextPage {
				opts_ := opts
				opts_.endCursor = project_.MergeRequests.PageInfo.EndCursor
				mergeRequests_, err_ := c.getProjectMergeRequests(ctx, project_.FullPath, opts_)
				if err_ != nil {
					err = err_
					break outerLoop
				}
				mergeRequests = append(mergeRequests, mergeRequests_...)
			}
		}

		if !data.Projects.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = data.Projects.PageInfo.EndCursor
	}

	return mergeRequests, err
}

func (c *Client) getProjectMergeRequests(ctx context.Context, projectPath string, opts getMergeRequestsOptions) ([]MergeRequestFields, error) {
	var (
		mergeRequests []MergeRequestFields

		data *getProjectMergeRequestsResponse
		err  error
	)

	for {
		data, err = getProjectMergeRequests(
			ctx,
			c.client,
			projectPath,
			opts.UpdatedAfter,
			opts.UpdatedBefore,
			opts.endCursor,
			//
			&opts.includeCore,
			&opts.includeExtra,
			&opts.includeParticipants,
		)
		err = handleError(err, "getProjectsMergeRequests",
			slog.String("projectPath", projectPath),
			slog.String("updatedAfter", opts.UpdatedAfter.Format(time.RFC3339)),
			slog.String("updatedBefore", opts.UpdatedBefore.Format(time.RFC3339)),
		)
		if err != nil {
			break
		}

		project_ := data.Project
		if project_ == nil {
			err = fmt.Errorf("project not found: %v", projectPath)
			break
		}

		for _, mergeRequest_ := range project_.MergeRequests.Nodes {
			mergeRequest := MergeRequestFields{
				MergeRequestReferenceFields: mergeRequest_.MergeRequestReferenceFields,
				Project:                     project_.ProjectReferenceFields,
			}
			if opts.includeCore {
				mergeRequest.MergeRequestFieldsCore = mergeRequest_.MergeRequestFieldsCore
			}
			if opts.includeExtra {
				mergeRequest.MergeRequestFieldsExtra = mergeRequest_.MergeRequestFieldsExtra
			}
			if opts.includeParticipants {
				mergeRequest.MergeRequestFieldsParticipants = mergeRequest_.MergeRequestFieldsParticipants
			}
			mergeRequests = append(mergeRequests, mergeRequest)
		}

		if !project_.MergeRequests.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = project_.MergeRequests.PageInfo.EndCursor
	}

	return mergeRequests, err
}

type MergeRequestNoteFields struct {
	MergeRequest MergeRequestReferenceFields
	Project      ProjectReferenceFields

	MergeRequestNotesFieldsCore
}

func ConvertMergeRequestNoteEvent(nf MergeRequestNoteFields) (types.MergeRequestNoteEvent, error) {
	var (
		id, mrId, mrIid, projectId int64
		err                        error
	)
	if id, err = parseNoteId(nf.Id); err != nil {
		return types.MergeRequestNoteEvent{}, fmt.Errorf("parse note id: %w", err)
	}
	if mrId, err = ParseId(nf.MergeRequest.Id, GlobalIdMergeRequestPrefix); err != nil {
		return types.MergeRequestNoteEvent{}, fmt.Errorf("parse merge request id: %w", err)
	}
	if mrIid, err = ParseId(nf.MergeRequest.Iid, ""); err != nil {
		return types.MergeRequestNoteEvent{}, fmt.Errorf("parse merge request iid: %w", err)
	}
	if projectId, err = ParseId(nf.Project.Id, GlobalIdProjectPrefix); err != nil {
		return types.MergeRequestNoteEvent{}, fmt.Errorf("parse project id: %w", err)
	}

	t := getNoteEventType(nf)
	// if t == "" {
	// 	// TODO: what?
	// }

	ne := types.MergeRequestNoteEvent{
		Id: id,
		MergeRequest: types.MergeRequestReference{
			Id:  mrId,
			Iid: mrIid,
			Project: types.ProjectReference{
				Id:       projectId,
				FullPath: nf.Project.FullPath,
			},
		},

		CreatedAt: &nf.CreatedAt,
		UpdatedAt: &nf.UpdatedAt,

		Type:     t,
		System:   nf.System,
		Internal: valOrZero(nf.Internal),

		// Author: {},

		Resolvable: nf.Resolvable,
		Resolved:   nf.Resolved,
		ResolvedAt: nf.ResolvedAt,
		// Resolver: {},
	}

	if nf.Author != nil {
		author, err := convertUserReference(nf.Author)
		if err != nil {
			return types.MergeRequestNoteEvent{}, fmt.Errorf("convert author reference: %w", err)
		}
		ne.Author = author
	}
	if nf.ResolvedBy != nil {
		resolver, err := convertUserReference(nf.ResolvedBy)
		if err != nil {
			return types.MergeRequestNoteEvent{}, fmt.Errorf("convert resolver reference: %w", err)
		}
		ne.Resolver = resolver
	}

	return ne, nil
}

func parseNoteId(s string) (int64, error) {
	var (
		id  int64
		err error
	)

	if id, err = ParseId(s, GlobalIdNotePrefix); err == nil {
		return id, nil
	}

	if id, err = strconv.ParseInt(s, 10, 64); err == nil {
		return id, nil
	}

	re := regexp.MustCompile(`.*/([0-9]+)$`)
	matches := re.FindStringSubmatch(s)
	if len(matches) != 2 {
		return 0, fmt.Errorf("failed to match %q", s)
	}

	return strconv.ParseInt(matches[1], 10, 64)
}

func getNoteEventType(note MergeRequestNoteFields) string {
	if !note.System {
		return ""
	}

	switch {
	case note.Body == "resolved all threads":
		return "AllThreadsResolved"

	case note.Body == "approved this merge request":
		return "Approved"
	case note.Body == "unapproved this merge request":
		return "Unapproved"

	case note.Body == "changed the description":
		return "DescriptionChanged"

	case note.Body == "marked this merge request as **draft**":
		return "MarkedDraft"
	case note.Body == "marked this merge request as **ready**":
		return "MarkedReady"

	case strings.HasPrefix(note.Body, "assigned to"):
		return "Assigned"
	case strings.HasPrefix(note.Body, "unassigned"):
		return "Unassigned"

	case strings.HasPrefix(note.Body, "requested review"):
		return "ReviewRequested"
	case strings.HasPrefix(note.Body, "removed review requested"):
		return "ReviewRequestRemoved"
	case note.Body == "requested changes":
		return "ChangesRequested"
	}

	return ""
}

func (c *Client) GetProjectsMergeRequestsNotes(ctx context.Context, projectGids []string, opts GetMergeRequestsOptions) ([]MergeRequestNoteFields, error) {
	return c.getProjectsMergeRequestsNotes(ctx, projectGids, getMergeRequestsNotesOptions{
		GetMergeRequestsOptions: opts,
	})
}

type getMergeRequestsNotesOptions struct {
	GetMergeRequestsOptions

	endCursor *string
}

func (c *Client) getProjectsMergeRequestsNotes(ctx context.Context, projectGids []string, opts getMergeRequestsNotesOptions) ([]MergeRequestNoteFields, error) {
	var (
		notes []MergeRequestNoteFields

		data *getProjectsMergeRequestNotesResponse
		err  error
	)

outerLoop:
	for {
		data, err = getProjectsMergeRequestNotes(
			ctx,
			c.client,
			projectGids,
			opts.UpdatedAfter,
			opts.UpdatedBefore,
			opts.endCursor,
		)
		err = handleError(err, "getProjectsMergeRequestsNotes",
			slog.Any("projectIds", projectGids),
			slog.String("updatedAfter", opts.UpdatedAfter.Format(time.RFC3339)),
			slog.String("updatedBefore", opts.UpdatedBefore.Format(time.RFC3339)),
		)
		if err != nil {
			break
		}

		for _, project_ := range data.Projects.Nodes {
			if project_.MergeRequests == nil {
				continue
			}
			for _, mr_ := range project_.MergeRequests.Nodes {
				for _, note_ := range mr_.Notes.Nodes {
					note := MergeRequestNoteFields{
						MergeRequest:                mr_.MergeRequestReferenceFields,
						Project:                     project_.ProjectReferenceFields,
						MergeRequestNotesFieldsCore: note_.MergeRequestNotesFieldsCore,
					}

					notes = append(notes, note)
				}

				if mr_.Notes.PageInfo.HasNextPage {
					opts_ := opts
					opts_.endCursor = mr_.Notes.PageInfo.EndCursor
					notes_, err_ := c.getProjectMergeRequestNotes(ctx, project_.FullPath, mr_.Iid, opts_)
					if err_ != nil {
						err = metaerr.WithMetadata(err_, "projectPath", project_.FullPath, "iid", mr_.Iid)
						break outerLoop
					}
					notes = append(notes, notes_...)
				}
			}

			if project_.MergeRequests.PageInfo.HasNextPage {
				opts_ := opts
				opts_.endCursor = project_.MergeRequests.PageInfo.EndCursor
				notes_, err_ := c.getProjectMergeRequestsNotes(ctx, project_.FullPath, opts_)
				if err_ != nil {
					err = metaerr.WithMetadata(err_, "projectPath", project_.FullPath)
					break outerLoop
				}
				notes = append(notes, notes_...)
			}
		}

		if !data.Projects.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = data.Projects.PageInfo.EndCursor
	}

	return notes, err
}

func (c *Client) getProjectMergeRequestsNotes(ctx context.Context, projectPath string, opts getMergeRequestsNotesOptions) ([]MergeRequestNoteFields, error) {
	var (
		notes []MergeRequestNoteFields

		data *getProjectMergeRequestsNotesResponse
		err  error
	)

outerLoop:
	for {
		data, err = getProjectMergeRequestsNotes(
			ctx,
			c.client,
			projectPath,
			opts.UpdatedAfter,
			opts.UpdatedBefore,
			opts.endCursor,
		)
		err = handleError(err, "getProjectMergeRequestsNotes",
			slog.String("projectPath", projectPath),
			slog.String("updatedAfter", opts.UpdatedAfter.Format(time.RFC3339)),
			slog.String("updatedBefore", opts.UpdatedBefore.Format(time.RFC3339)),
		)
		if err != nil {
			break
		}

		project_ := data.Project
		if project_ == nil {
			err = fmt.Errorf("project not found: %v", projectPath)
			break
		}
		if project_.MergeRequests == nil {
			break
		}

		for _, mr_ := range project_.MergeRequests.Nodes {
			for _, note_ := range mr_.Notes.Nodes {
				note := MergeRequestNoteFields{
					MergeRequest: mr_.MergeRequestReferenceFields,
					Project:      project_.ProjectReferenceFields,

					MergeRequestNotesFieldsCore: note_.MergeRequestNotesFieldsCore,
				}

				notes = append(notes, note)
			}

			if mr_.Notes.PageInfo.HasNextPage {
				opts_ := opts
				opts_.endCursor = mr_.Notes.PageInfo.EndCursor
				notes_, err_ := c.getProjectMergeRequestNotes(ctx, project_.FullPath, mr_.Iid, opts_)
				if err_ != nil {
					err = metaerr.WithMetadata(err_, "projectPath", project_.FullPath, "iid", mr_.Iid)
					break outerLoop
				}
				notes = append(notes, notes_...)
			}
		}

		if !project_.MergeRequests.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = project_.MergeRequests.PageInfo.EndCursor
	}

	return notes, err
}

func (c *Client) getProjectMergeRequestNotes(ctx context.Context, projectPath string, mergeRequestIid string, opts getMergeRequestsNotesOptions) ([]MergeRequestNoteFields, error) {
	var (
		notes []MergeRequestNoteFields

		data *getProjectMergeRequestNotesResponse
		err  error
	)

	for {
		data, err = getProjectMergeRequestNotes(
			ctx,
			c.client,
			projectPath,
			mergeRequestIid,
			opts.endCursor,
		)
		err = handleError(err, "getProjectMergeRequestNotes",
			slog.String("projectPath", projectPath),
			slog.String("iid", mergeRequestIid),
		)
		if err != nil {
			break
		}

		project_ := data.Project
		if project_ == nil {
			err = fmt.Errorf("project not found: %v", projectPath)
			break
		}
		mr_ := project_.MergeRequest
		if mr_ == nil {
			err = fmt.Errorf("project merge request not found: %v/%v", projectPath, mergeRequestIid)
			break
		}

		for _, note_ := range mr_.Notes.Nodes {
			note := MergeRequestNoteFields{
				MergeRequest:                mr_.MergeRequestReferenceFields,
				Project:                     project_.ProjectReferenceFields,
				MergeRequestNotesFieldsCore: note_.MergeRequestNotesFieldsCore,
			}

			notes = append(notes, note)
		}

		if !mr_.Notes.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = mr_.Notes.PageInfo.EndCursor
	}

	return notes, err
}
