package types

import (
	"time"
)

type JobKind string

const (
	JobKindUnkown JobKind = "unknown"
	JobKindBuild  JobKind = "build"
	JobKindBridge JobKind = "bridge"
)

type JobReference struct {
	Id   int64
	Name string

	Pipeline PipelineReference
}

type Job struct {
	Id       int64
	Name     string
	Pipeline PipelineReference

	Ref           string
	RefPath       string
	Status        string
	FailureReason string

	CreatedAt  *time.Time
	QueuedAt   *time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
	ErasedAt   *time.Time

	QueuedDuration time.Duration
	Duration       time.Duration
	Coverage       float64

	Stage      string
	Tags       []string
	Properties []JobLogProperty

	AllowFailure bool
	Manual       bool
	Retried      bool
	Retryable    bool

	Kind               JobKind
	DownstreamPipeline *PipelineReference

	RunnerId string
	Trace    *JobTrace
}

type JobLogProperty struct {
	Name  string
	Value string
}

type JobArtifact struct {
	Job JobReference

	FileType     string
	Name         string
	DownloadPath string
}

type JobTrace struct {}
