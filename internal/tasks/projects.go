package tasks

import (
	"context"

	"github.com/cluttrdev/gitlab-exporter/internal/exporter"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type ExportProjectOptions struct {
	ProjectID int64
}

func ExportProject(ctx context.Context, glab *gitlab.Client, exp *exporter.Exporter, opt ExportProjectOptions) error {
	p, err := glab.GetProject(ctx, opt.ProjectID)
	if err != nil {
		return err
	}

	if err := exp.ExportProjects(ctx, []*typespb.Project{p}); err != nil {
		return err
	}

	return nil
}
