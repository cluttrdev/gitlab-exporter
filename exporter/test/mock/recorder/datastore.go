package recorder_mock

import (
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

type Datastore interface {
	ListProjects() []*typespb.Project
	GetProject(id int64) *typespb.Project

	ListRunners() []*typespb.Runner
	GetRunner(id int64) *typespb.Runner

	ListProjectPipelines(projectID int64) []*typespb.Pipeline
	GetPipeline(id int64) *typespb.Pipeline

	ListPipelineJobs(projectID int64, pipelineID int64) []*typespb.Job

	GetPipelineTestReport(pipelineID int64) *typespb.TestReport
}

type datastore struct {
	projects []*typespb.Project
	runners  []*typespb.Runner

	pipelines []*typespb.Pipeline
	jobs      []*typespb.Job
	sections  []*typespb.Section

	testreports []*typespb.TestReport
	testsuites  []*typespb.TestSuite
	testcases   []*typespb.TestCase

	traces  []*typespb.Trace
	metrics []*typespb.Metric
}

func (d *datastore) ListProjects() []*typespb.Project {
	return d.projects
}

func (d *datastore) GetProject(id int64) *typespb.Project {
	for _, p := range d.projects {
		if p.Id == id {
			return p
		}
	}
	return nil
}

func (d *datastore) ListRunners() []*typespb.Runner {
	return d.runners
}

func (d *datastore) GetRunner(id int64) *typespb.Runner {
	for _, r := range d.runners {
		if r.Id == id {
			return r
		}
	}
	return nil
}

func (d *datastore) ListProjectPipelines(projectID int64) []*typespb.Pipeline {
	var ps []*typespb.Pipeline
	for _, p := range d.pipelines {
		if p.Project.Id == projectID {
			ps = append(ps, p)
		}
	}
	return ps
}

func (d *datastore) GetPipeline(id int64) *typespb.Pipeline {
	for _, p := range d.pipelines {
		if p.GetId() == id {
			return p
		}
	}
	return nil
}

func (d *datastore) ListPipelineJobs(projectID int64, pipelineID int64) []*typespb.Job {
	var js []*typespb.Job
	for _, j := range d.jobs {
		if j.Pipeline.Project.Id == projectID && j.Pipeline.Id == pipelineID {
			js = append(js, j)
		}
	}
	return js
}

func (d *datastore) GetPipelineTestReport(pipelineID int64) *typespb.TestReport {
	for _, tr := range d.testreports {
		if tr.GetJob().GetPipeline().GetId() == pipelineID {
			return tr
		}
	}
	return nil
}
