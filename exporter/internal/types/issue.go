package types

import "time"

type Issue struct {
	Id      int64
	Iid     int64
	Project ProjectReference

	CreatedAt *time.Time
	UpdatedAt *time.Time
	ClosedAt  *time.Time

	Title  string
	Labels []string

	Type     IssueType
	Severity IssueSeverity
	State    IssueState
}

type IssueType string

const (
	IssueTypeUnspecified IssueType = ""
	IssueTypeEpic        IssueType = "EPIC" // (experimental) introduced in 16.7
	IssueTypeIncident    IssueType = "INCIDENT"
	IssueTypeIssue       IssueType = "ISSUE"
	IssueTypeKeyResult   IssueType = "KEY_RESULT" // (experimental) introduced in 15.7
	IssueTypeObjective   IssueType = "OBJECTIVE"  // (experimental) introduced in 15.6
	IssueTypeRequirement IssueType = "REQUIREMENT"
	IssueTypeTask        IssueType = "TASK"
	IssueTypeTestCase    IssueType = "TEST_CASE"
	IssueTypeTicket      IssueType = "TICKET"
	IssueTypeUnknown     IssueType = "UNKNOWN"
)

type IssueSeverity string

const (
	IssueSeverityUnspecified IssueSeverity = ""
	IssueSeverityCritical    IssueSeverity = "CRITICAL"
	IssueSeverityHigh        IssueSeverity = "HIGH"
	IssueSeverityLow         IssueSeverity = "LOW"
	IssueSeverityMedium      IssueSeverity = "MEDIUM"
	IssueSeverityUnknown     IssueSeverity = "UNKNOWN"
)

type IssueState string

const (
	IssueStateUnknown IssueState = ""
	IssueStateOpened  IssueState = "opened"
	IssueStateClosed  IssueState = "closed"
	IssueStateLocked  IssueState = "locked"
	IssueStateAll     IssueState = "all"
)
