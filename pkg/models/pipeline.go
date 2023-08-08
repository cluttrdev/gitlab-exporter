package models

import (
	"fmt"
	"strconv"
	"time"

	gogitlab "github.com/xanzy/go-gitlab"
)

type PipelineInfo struct {
	ID        int64
	IID       int64
	ProjectID int64
	Status    string
	Source    string
	Ref       string
	SHA       string
	WebURL    string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type Pipeline struct {
	ID         int64  `json:"id"`
	IID        int64  `json:"iid"`
	ProjectID  int64  `json:"project_id"`
	Status     string `json:"status"`
	Source     string `json:"source"`
	Ref        string `json:"ref"`
	SHA        string `json:"sha"`
	BeforeSHA  string `json:"before_sha"`
	Tag        bool   `json:"tag"`
	YamlErrors string `json:"yaml_errors"`
	// User           *User           `json:"user"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
	StartedAt      *time.Time `json:"started_at"`
	FinishedAt     *time.Time `json:"finished_at"`
	CommittedAt    *time.Time `json:"committed_at"`
	Duration       float64    `json:"duration"`
	QueuedDuration float64    `json:"queued_duration"`
	Coverage       float64    `json:"coverage,string"`
	WebURL         string     `json:"web_url"`
	// DetailedStatus *DetailedStatus `json:"detailed_status"`
}

func NewPipelineInfo(p *gogitlab.PipelineInfo) *PipelineInfo {
	return &PipelineInfo{
		ID:        int64(p.ID),
		IID:       int64(p.IID),
		ProjectID: int64(p.ProjectID),
		Status:    p.Status,
		Source:    p.Source,
		Ref:       p.Ref,
		SHA:       p.SHA,
		WebURL:    p.WebURL,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func nullPipelineInfo() *PipelineInfo {
	nullTime := time.UnixMilli(0).UTC()
	return &PipelineInfo{
		CreatedAt: &nullTime,
		UpdatedAt: &nullTime,
	}
}

func NewPipeline(p *gogitlab.Pipeline) *Pipeline {
	cov, err := strconv.ParseFloat(p.Coverage, 64)
	if err != nil {
		cov = 0.0
	}
	return &Pipeline{
		ID:             int64(p.ID),
		IID:            int64(p.IID),
		ProjectID:      int64(p.ProjectID),
		Status:         p.Status,
		Source:         p.Source,
		Ref:            p.Ref,
		SHA:            p.SHA,
		BeforeSHA:      p.BeforeSHA,
		Tag:            p.Tag,
		YamlErrors:     p.YamlErrors,
		CreatedAt:      p.CreatedAt,
		UpdatedAt:      p.UpdatedAt,
		StartedAt:      p.StartedAt,
		FinishedAt:     p.FinishedAt,
		CommittedAt:    p.CommittedAt,
		Duration:       float64(p.Duration),
		QueuedDuration: float64(p.QueuedDuration),
		Coverage:       cov,
		WebURL:         p.WebURL,
	}
}

type PipelineHierarchy struct {
    Pipeline            *Pipeline `json:"pipeline"`
    Jobs                []*Job `json:"jobs"`
    Sections            []*Section `json:"sections"`
    Bridges             []*Bridge `json:"bridges"`
    DownstreamPipelines []*PipelineHierarchy `json:"downstream_pipelines"`
}

func (ph *PipelineHierarchy) GetAllPipelines() []*Pipeline {
	pipelines := []*Pipeline{ph.Pipeline}

	for _, dph := range ph.DownstreamPipelines {
		pipelines = append(pipelines, dph.GetAllPipelines()...)
	}

	return pipelines
}

func (ph *PipelineHierarchy) GetAllJobs() []*Job {
	jobs := ph.Jobs

	for _, dph := range ph.DownstreamPipelines {
		jobs = append(jobs, dph.GetAllJobs()...)
	}

	return jobs
}

func (ph *PipelineHierarchy) GetAllSections() []*Section {
	sections := ph.Sections

	for _, dph := range ph.DownstreamPipelines {
		sections = append(sections, dph.GetAllSections()...)
	}

	return sections
}

func (ph *PipelineHierarchy) GetAllBridges() []*Bridge {
	bridges := ph.Bridges

	for _, dph := range ph.DownstreamPipelines {
		bridges = append(bridges, dph.GetAllBridges()...)
	}

	return bridges
}

func (ph *PipelineHierarchy) GetTrace() []*Span {
	traceID := fmt.Sprintf("%d", ph.Pipeline.ID)
	parentID := ""
	return NewPipelineHierarchyTrace(traceID, parentID, ph)
}

func (ph *PipelineHierarchy) GetAllTraces() [][]*Span {
	traces := [][]*Span{ph.GetTrace()}
	for _, dph := range ph.DownstreamPipelines {
		traces = append(traces, dph.GetAllTraces()...)
	}
	return traces
}
