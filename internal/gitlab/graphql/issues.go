package graphql

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/types"
)

type IssueFields struct {
	IssueReferenceFields
	Project ProjectReferenceFields

	IssueFieldsCore
}

func ConvertIssue(isf IssueFields) (types.Issue, error) {
	var (
		id, iid, projectId int64
		err                error
	)
	if id, err = ParseId(isf.Id, GlobalIdIssuePrefix); err != nil {
		return types.Issue{}, fmt.Errorf("parse issue id: %w", err)
	}
	if iid, err = ParseId(isf.Iid, ""); err != nil {
		return types.Issue{}, fmt.Errorf("parse issue iid: %w", err)
	}
	if projectId, err = ParseId(isf.Project.Id, GlobalIdProjectPrefix); err != nil {
		return types.Issue{}, fmt.Errorf("parse project id: %w", err)
	}

	var labels []string
	for _, label := range isf.Labels.Nodes {
		labels = append(labels, label.Title)
	}

	return types.Issue{
		Id:  id,
		Iid: iid,
		Project: types.ProjectReference{
			Id:       projectId,
			FullPath: isf.Project.FullPath,
		},

		CreatedAt: &isf.CreatedAt,
		UpdatedAt: &isf.UpdatedAt,
		ClosedAt:  isf.ClosedAt,

		Title:  isf.Title,
		Labels: labels,

		Type:     convertIssueType(isf.Type),
		Severity: convertIssueSeverity(isf.Severity),
		State:    convertIssueState(isf.State),
	}, nil
}

func convertIssueType(t *IssueType) types.IssueType {
	if t == nil {
		return types.IssueTypeUnspecified
	}

	switch *t {
	case IssueTypeEpic:
		return types.IssueTypeEpic
	case IssueTypeIncident:
		return types.IssueTypeIncident
	case IssueTypeIssue:
		return types.IssueTypeIssue
	case IssueTypeKeyResult:
		return types.IssueTypeKeyResult
	case IssueTypeObjective:
		return types.IssueTypeObjective
	case IssueTypeRequirement:
		return types.IssueTypeRequirement
	case IssueTypeTask:
		return types.IssueTypeTask
	case IssueTypeTestCase:
		return types.IssueTypeTestCase
		// case IssueTypeTicket:
		//     return types.IssueTypeTicket
	}

	if string(*t) == "TICKET" {
		return types.IssueTypeTicket
	}

	return types.IssueTypeUnknown
}

func convertIssueSeverity(s *IssuableSeverity) types.IssueSeverity {
	if s == nil {
		return types.IssueSeverityUnspecified
	}

	switch *s {
	case IssuableSeverityCritical:
		return types.IssueSeverityCritical
	case IssuableSeverityHigh:
		return types.IssueSeverityHigh
	case IssuableSeverityLow:
		return types.IssueSeverityLow
	case IssuableSeverityMedium:
		return types.IssueSeverityMedium
	case IssuableSeverityUnknown:
		return types.IssueSeverityUnknown
	}

	return types.IssueSeverityUnknown
}

func convertIssueState(s IssueState) types.IssueState {
	switch s {
	case IssueStateAll:
		return types.IssueStateAll
	case IssueStateClosed:
		return types.IssueStateClosed
	case IssueStateLocked:
		return types.IssueStateLocked
	case IssueStateOpened:
		return types.IssueStateOpened
	}

	return types.IssueStateUnknown
}

type GetProjectsIssuesOptions struct {
	TimeRangeOptions
}

func (c *Client) GetProjectsIssues(ctx context.Context, projectIds []string, opts GetProjectsIssuesOptions) ([]IssueFields, error) {
	return c.getProjectsIssues(ctx, projectIds, getProjectsIssuesOptions{GetProjectsIssuesOptions: opts})
}

type getProjectsIssuesOptions struct {
	GetProjectsIssuesOptions

	endCursor *string
}

func (c *Client) getProjectsIssues(ctx context.Context, projectIds []string, opts getProjectsIssuesOptions) ([]IssueFields, error) {
	var (
		issues []IssueFields

		data *getProjectsIssuesResponse
		err  error
	)

	for {
		data, err = getProjectsIssues(ctx, c.client, projectIds, opts.UpdatedAfter, opts.UpdatedBefore, opts.endCursor)
		err = handleError(err, "getProjectsIssues", slog.Any("ids", projectIds))
		if err != nil {
			break
		}

		for _, project_ := range data.GetProjects().GetNodes() {
			for _, issue_ := range project_.GetIssues().GetNodes() {
				issue := IssueFields{
					IssueReferenceFields: issue_.IssueReferenceFields,
					Project:              project_.ProjectReferenceFields,

					IssueFieldsCore: issue_.IssueFieldsCore,
				}

				issues = append(issues, issue)
			}

			if project_.GetIssues().GetPageInfo().HasNextPage {
				opts_ := opts
				opts_.endCursor = project_.GetIssues().GetPageInfo().EndCursor
				issues_, err := c.getProjectIssues(ctx, project_.GetFullPath(), opts_)
				if err != nil {
					return nil, err
				}
				issues = append(issues, issues_...)
			}
		}

		if !data.GetProjects().GetPageInfo().HasNextPage {
			break
		}

		opts.endCursor = data.GetProjects().GetPageInfo().EndCursor
	}

	return issues, err
}

func (c *Client) getProjectIssues(ctx context.Context, projectPath string, opts getProjectsIssuesOptions) ([]IssueFields, error) {
	var (
		issues []IssueFields

		data *getProjectIssuesResponse
		err  error
	)

	for {
		data, err = getProjectIssues(ctx, c.client, projectPath, opts.UpdatedAfter, opts.UpdatedBefore, opts.endCursor)
		err = handleError(err, "getProjectIssues",
			slog.String("projectPath", projectPath),
			slog.String("updatedAfter", opts.UpdatedAfter.Format(time.RFC3339)),
			slog.String("updatedBefore", opts.UpdatedBefore.Format(time.RFC3339)),
		)
		if err != nil {
			break
		}

		project_ := data.GetProject()
		if project_ == nil {
			err = fmt.Errorf("project not found: %v", projectPath)
			break
		}

		for _, issue_ := range project_.GetIssues().GetNodes() {
			issue := IssueFields{
				IssueReferenceFields: issue_.IssueReferenceFields,
				Project:              project_.ProjectReferenceFields,

				IssueFieldsCore: issue_.IssueFieldsCore,
			}

			issues = append(issues, issue)
		}

		if !project_.GetIssues().GetPageInfo().HasNextPage {
			break
		}

		opts.endCursor = project_.GetIssues().GetPageInfo().EndCursor
	}

	return issues, err
}
