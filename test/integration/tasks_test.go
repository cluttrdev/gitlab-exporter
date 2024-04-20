package integration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/tasks"
)

func Test_ExportPipelineHierarchy(t *testing.T) {
	env, err := GetTestEnvironment(testSet)
	if err != nil {
		t.Error(err)
	}

	glc, err := gitlab.NewGitLabClient(gitlab.ClientConfig{
		URL: fmt.Sprintf("%s/api/v4", env.URL),
	})
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	exp, rec := setupExporter(t)

	var (
		projectID  int64 = 50817395 // cluttrdev/gitlab-exporter
		pipelineID int64 = 1252248442
	)

	opts := tasks.ExportPipelineHierarchyOptions{
		ProjectID:  projectID,
		PipelineID: pipelineID,

		ExportSections:    false,
		ExportTestReports: true,
		ExportTraces:      true,
		ExportMetrics:     false,
	}

	if err := tasks.ExportPipelineHierarchy(context.Background(), glc, exp, opts); err != nil {
		t.Error(err)
	}

	assert.Equal(t, 1, len(rec.Datastore().ListProjectPipelines(projectID)))

	p := rec.Datastore().GetPipeline(pipelineID)
	if p == nil {
		t.Fatalf("pipeline not recorded: %v", pipelineID)
	}

	assert.Equal(t, 4, len(rec.Datastore().ListPipelineJobs(projectID, pipelineID)))
	assert.Equal(t, int64(13), rec.Datastore().GetPipelineTestReport(pipelineID).TotalCount)
}
