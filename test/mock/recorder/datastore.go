package recorder_mock

import (
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type Datastore interface {
	ListProjects() []*typespb.Project
	GetProject(id int64) *typespb.Project

	ListProjectPipelines(projectID int64) []*typespb.Pipeline
	GetPipeline(id int64) *typespb.Pipeline

	ListPipelineJobs(projectID int64, pipelineID int64) []*typespb.Job

	GetPipelineTestReport(pipelineID int64) *typespb.TestReport
}

type datastore struct {
	projects []*typespb.Project

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
		if tr.Pipeline.Id == pipelineID {
			return tr
		}
	}
	return nil
}
