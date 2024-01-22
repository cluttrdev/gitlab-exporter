package models

import (
	"fmt"

	pb "github.com/cluttrdev/gitlab-exporter/grpc/exporterpb"
)

type PipelineHierarchy struct {
	Pipeline            *pb.Pipeline         `json:"pipeline"`
	Jobs                []*pb.Job            `json:"jobs"`
	Sections            []*pb.Section        `json:"sections"`
	Bridges             []*pb.Bridge         `json:"bridges"`
	DownstreamPipelines []*PipelineHierarchy `json:"downstream_pipelines"`
}

func (ph *PipelineHierarchy) GetAllPipelines() []*pb.Pipeline {
	pipelines := []*pb.Pipeline{ph.Pipeline}

	for _, dph := range ph.DownstreamPipelines {
		pipelines = append(pipelines, dph.GetAllPipelines()...)
	}

	return pipelines
}

func (ph *PipelineHierarchy) GetAllJobs() []*pb.Job {
	jobs := ph.Jobs

	for _, dph := range ph.DownstreamPipelines {
		jobs = append(jobs, dph.GetAllJobs()...)
	}

	return jobs
}

func (ph *PipelineHierarchy) GetAllSections() []*pb.Section {
	sections := ph.Sections

	for _, dph := range ph.DownstreamPipelines {
		sections = append(sections, dph.GetAllSections()...)
	}

	return sections
}

func (ph *PipelineHierarchy) GetAllBridges() []*pb.Bridge {
	bridges := ph.Bridges

	for _, dph := range ph.DownstreamPipelines {
		bridges = append(bridges, dph.GetAllBridges()...)
	}

	return bridges
}

func (ph *PipelineHierarchy) GetTrace() *pb.Trace {
	traceID := fmt.Sprintf("%d", ph.Pipeline.Id)
	parentID := ""
	return NewPipelineHierarchyTrace([]byte(traceID), []byte(parentID), ph)
}

func (ph *PipelineHierarchy) GetAllTraces() []*pb.Trace {
	traces := []*pb.Trace{ph.GetTrace()}
	for _, dph := range ph.DownstreamPipelines {
		traces = append(traces, dph.GetAllTraces()...)
	}
	return traces
}
