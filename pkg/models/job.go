package models

import (
	"time"

	_gitlab "github.com/xanzy/go-gitlab"
)

type Job struct {
	// Commit            *Commit    `json:"commit"`
	Coverage       float64    `json:"coverage"`
	AllowFailure   bool       `json:"allow_failure"`
	CreatedAt      *time.Time `json:"created_at"`
	StartedAt      *time.Time `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
	ErasedAt       *time.Time `json:"erased_at"`
	Duration       float64    `json:"duration"`
	QueuedDuration float64    `json:"queued_duration"`
	// ArtifactsExpireAt *time.Time `json:"artifacts_expire_at"`
	TagList  []string `json:"tag_list"`
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Pipeline struct {
		ID        int64  `json:"id"`
		ProjectID int64  `json:"project_id"`
		Ref       string `json:"ref"`
		Sha       string `json:"sha"`
		Status    string `json:"status"`
	} `json:"pipeline"`
	Ref string `json:"ref"`
	// Artifacts []struct {
	//  FileType   string `json:"file_type"`
	//  Filename   string `json:"filename"`
	//  Size       int    `json:"size"`
	//  FileFormat string `json:"file_format"`
	// } `json:"artifacts"`
	// ArtifactsFile struct {
	//  Filename string `json:"filename"`
	//  Size     int    `json:"size"`
	// } `json:"artifacts_file"`
	// Runner struct {
	//  ID          int    `json:"id"`
	//  Description string `json:"description"`
	//  Active      bool   `json:"active"`
	//  IsShared    bool   `json:"is_shared"`
	//  Name        string `json:"name"`
	// } `json:"runner"`
	Stage         string `json:"stage"`
	Status        string `json:"status"`
	FailureReason string `json:"failure_reason"`
	Tag           bool   `json:"tag"`
	WebURL        string `json:"web_url"`
	// Project       *Project `json:"project"`
	// User          *User    `json:"user"`
}

type Bridge struct {
	// Commit             *Commit       `json:"commit"`
	Coverage       float64      `json:"coverage"`
	AllowFailure   bool         `json:"allow_failure"`
	CreatedAt      *time.Time   `json:"created_at"`
	StartedAt      *time.Time   `json:"started_at"`
	FinishedAt     *time.Time   `json:"finished_at"`
	ErasedAt       *time.Time   `json:"erased_at"`
	Duration       float64      `json:"duration"`
	QueuedDuration float64      `json:"queued_duration"`
	ID             int64        `json:"id"`
	Name           string       `json:"name"`
	Pipeline       PipelineInfo `json:"pipeline"`
	Ref            string       `json:"ref"`
	Stage          string       `json:"stage"`
	Status         string       `json:"status"`
	FailureReason  string       `json:"failure_reason"`
	Tag            bool         `json:"tag"`
	WebURL         string       `json:"web_url"`
	// User               *User         `json:"user"`
	DownstreamPipeline *PipelineInfo `json:"downstream_pipeline"`
}

func NewJob(j *_gitlab.Job) *Job {
	return &Job{
		Coverage:       j.Coverage,
		AllowFailure:   j.AllowFailure,
		CreatedAt:      j.CreatedAt,
		StartedAt:      j.StartedAt,
		FinishedAt:     j.FinishedAt,
		ErasedAt:       j.ErasedAt,
		Duration:       j.Duration,
		QueuedDuration: j.QueuedDuration,
		TagList:        j.TagList,
		ID:             int64(j.ID),
		Name:           j.Name,
		Pipeline: struct {
			ID        int64  `json:"id"`
			ProjectID int64  `json:"project_id"`
			Ref       string `json:"ref"`
			Sha       string `json:"sha"`
			Status    string `json:"status"`
		}{
			ID:        int64(j.Pipeline.ID),
			ProjectID: int64(j.Pipeline.ProjectID),
			Ref:       j.Pipeline.Ref,
			Sha:       j.Pipeline.Sha,
			Status:    j.Pipeline.Status,
		},
		Ref:           j.Ref,
		Stage:         j.Stage,
		Status:        j.Status,
		FailureReason: j.FailureReason,
		Tag:           j.Tag,
		WebURL:        j.WebURL,
	}
}

func NewBridge(b *_gitlab.Bridge) *Bridge {
	// account for downstream pipeline creation failures
	dp := nullPipelineInfo()
	if b.DownstreamPipeline != nil {
		dp = NewPipelineInfo(b.DownstreamPipeline)
	}
	return &Bridge{
		Coverage:           b.Coverage,
		AllowFailure:       b.AllowFailure,
		CreatedAt:          b.CreatedAt,
		StartedAt:          b.StartedAt,
		FinishedAt:         b.FinishedAt,
		ErasedAt:           b.ErasedAt,
		Duration:           b.Duration,
		QueuedDuration:     b.QueuedDuration,
		ID:                 int64(b.ID),
		Name:               b.Name,
		Pipeline:           *NewPipelineInfo(&b.Pipeline),
		Ref:                b.Ref,
		Stage:              b.Stage,
		Status:             b.Status,
		FailureReason:      b.FailureReason,
		Tag:                b.Tag,
		WebURL:             b.WebURL,
		DownstreamPipeline: dp,
	}
}
