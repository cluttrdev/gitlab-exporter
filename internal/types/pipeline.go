package types

import (
	"time"
)

type PipelineReference struct {
	Id  int64
	Iid int64

	Project ProjectReference
}

type Pipeline struct {
	Id  int64
	Iid int64

	Project ProjectReference

	Name          string
	Ref           string
	RefPath       string
	Sha           string
	Source        string
	Status        string
	FailureReason string

	CommittedAt *time.Time
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	StartedAt   *time.Time
	FinishedAt  *time.Time

	QueuedDuration time.Duration
	Duration       time.Duration
	Coverage       float64

	Warnings   bool
	YamlErrors bool

	Child               bool
	UpstreamPipeline    *PipelineReference
	DownstreamPipelines []*PipelineReference

	MergeRequest *MergeRequestReference

	User UserReference
}
