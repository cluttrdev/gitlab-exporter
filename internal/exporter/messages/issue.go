package messages

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewIssue(issue types.Issue) *typespb.Issue {
	pbIssue := &typespb.Issue{
		Id:      issue.Id,
		Iid:     issue.Iid,
		Project: NewProjectReference(issue.Project),

		Timestamps: &typespb.IssueTimestamps{
			CreatedAt: timestamppb.New(valOrZero(issue.CreatedAt)),
			UpdatedAt: timestamppb.New(valOrZero(issue.UpdatedAt)),
			ClosedAt:  timestamppb.New(valOrZero(issue.ClosedAt)),
		},

		Title:  issue.Title,
		Labels: issue.Labels,

		Type:     convertIssueType(issue.Type),
		Severity: convertIssueSeverity(issue.Severity),
		State:    convertIssueState(issue.State),
	}

	return pbIssue
}

func convertIssueType(t types.IssueType) typespb.IssueType {
	switch t {
	case types.IssueTypeEpic:
		return typespb.IssueType_ISSUE_TYPE_EPIC
	case types.IssueTypeIncident:
		return typespb.IssueType_ISSUE_TYPE_INCIDENT
	case types.IssueTypeIssue:
		return typespb.IssueType_ISSUE_TYPE_ISSUE
	case types.IssueTypeKeyResult:
		return typespb.IssueType_ISSUE_TYPE_KEY_RESULT
	case types.IssueTypeObjective:
		return typespb.IssueType_ISSUE_TYPE_OBJECTIVE
	case types.IssueTypeRequirement:
		return typespb.IssueType_ISSUE_TYPE_REQUIREMENT
	case types.IssueTypeTask:
		return typespb.IssueType_ISSUE_TYPE_TASK
	case types.IssueTypeTestCase:
		return typespb.IssueType_ISSUE_TYPE_TEST_CASE
	case types.IssueTypeTicket:
		return typespb.IssueType_ISSUE_TYPE_TICKET
	case types.IssueTypeUnknown:
		return typespb.IssueType_ISSUE_TYPE_UNKNOWN
	}

	return typespb.IssueType_ISSUE_TYPE_UNSPECIFIED
}

func convertIssueSeverity(s types.IssueSeverity) typespb.IssueSeverity {
	switch s {
	case types.IssueSeverityCritical:
		return typespb.IssueSeverity_ISSUE_SEVERITY_CRITICAL
	case types.IssueSeverityHigh:
		return typespb.IssueSeverity_ISSUE_SEVERITY_HIGH
	case types.IssueSeverityLow:
		return typespb.IssueSeverity_ISSUE_SEVERITY_LOW
	case types.IssueSeverityMedium:
		return typespb.IssueSeverity_ISSUE_SEVERITY_MEDIUM
	case types.IssueSeverityUnknown:
		return typespb.IssueSeverity_ISSUE_SEVERITY_UNKNOWN
	}

	return typespb.IssueSeverity_ISSUE_SEVERITY_UNSPECIFIED
}

func convertIssueState(s types.IssueState) typespb.IssueState {
	switch s {
	case types.IssueStateOpened:
		return typespb.IssueState_ISSUE_STATE_OPENED
	case types.IssueStateClosed:
		return typespb.IssueState_ISSUE_STATE_CLOSED
	case types.IssueStateLocked:
		return typespb.IssueState_ISSUE_STATE_LOCKED
	case types.IssueStateAll:
		return typespb.IssueState_ISSUE_STATE_ALL
	}

	return typespb.IssueState_ISSUE_STATE_UNKNOWN
}
