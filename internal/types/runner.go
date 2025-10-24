package types

import "time"

type RunnerReference struct {
	Id       int64
	ShortSha string
}

type Runner struct {
	Id          int64
	ShortSha    string
	Description string

	RunnerType RunnerType
	TagList    []string

	Status RunnerStatus

	Locked bool
	Paused bool

	RunProtected bool
	RunUntagged  bool

	CreatedAt   *time.Time
	ContactedAt *time.Time

	CreatedBy UserReference
}

type RunnerType string

const (
	RunnerTypeInstance RunnerType = "INSTANCE"
	RunnerTypeGroup    RunnerType = "GROUP"
	RunnerTypeProject  RunnerType = "PROJECT"
	RunnerTypeUnknown  RunnerType = "UNKNOWN"
)

type RunnerStatus string

const (
	// Runner that contacted the instance within the last 2 hours.
	RunnerStatusOnline RunnerStatus = "ONLINE"
	// Runner that has not contacted the instance within the last 2 hours.
	RunnerStatusOffline RunnerStatus = "OFFLINE"
	// Runner that has not contacted the instance within the last 7 days.
	RunnerStatusStale RunnerStatus = "STALE"
	// Runner that has never contacted the instance.
	RunnerStatusNeverContacted RunnerStatus = "NEVER_CONTACTED"
	// Unknown runner status.
	RunnerStatusUnknown RunnerStatus = "UNKNOWN"
)

type RunnerAccessLevel string

const (
	RunnerAccessLevelNotProtected RunnerAccessLevel = "NOT_PROTECTED"
	RunnerAccessLevelRefProtected RunnerAccessLevel = "REF_PROTECTED"
	RunnerAccessLevelUnknown      RunnerAccessLevel = "UNKNOWN"
)
