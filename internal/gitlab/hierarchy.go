package gitlab

import (
	"fmt"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type PipelineHierarchy struct {
	Pipeline            *typespb.Pipeline    `json:"pipeline"`
	Jobs                []*typespb.Job       `json:"jobs"`
	Sections            []*typespb.Section   `json:"sections"`
	Bridges             []*typespb.Bridge    `json:"bridges"`
	DownstreamPipelines []*PipelineHierarchy `json:"downstream_pipelines"`
}

func (ph *PipelineHierarchy) GetAllPipelines() []*typespb.Pipeline {
	pipelines := []*typespb.Pipeline{ph.Pipeline}

	for _, dph := range ph.DownstreamPipelines {
		pipelines = append(pipelines, dph.GetAllPipelines()...)
	}

	return pipelines
}

func (ph *PipelineHierarchy) GetAllJobs() []*typespb.Job {
	jobs := ph.Jobs

	for _, dph := range ph.DownstreamPipelines {
		jobs = append(jobs, dph.GetAllJobs()...)
	}

	return jobs
}

func (ph *PipelineHierarchy) GetAllSections() []*typespb.Section {
	sections := ph.Sections

	for _, dph := range ph.DownstreamPipelines {
		sections = append(sections, dph.GetAllSections()...)
	}

	return sections
}

func (ph *PipelineHierarchy) GetAllBridges() []*typespb.Bridge {
	bridges := ph.Bridges

	for _, dph := range ph.DownstreamPipelines {
		bridges = append(bridges, dph.GetAllBridges()...)
	}

	return bridges
}

func (ph *PipelineHierarchy) GetTrace() *typespb.Trace {
	traceID := fmt.Sprintf("%d", ph.Pipeline.Id)
	parentID := ""
	return NewPipelineHierarchyTrace([]byte(traceID), []byte(parentID), ph)
}

func (ph *PipelineHierarchy) GetAllTraces() []*typespb.Trace {
	traces := []*typespb.Trace{ph.GetTrace()}
	for _, dph := range ph.DownstreamPipelines {
		traces = append(traces, dph.GetAllTraces()...)
	}
	return traces
}
