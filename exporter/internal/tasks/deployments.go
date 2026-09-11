package tasks

import (
	"context"
	"slices"
	"time"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/rest"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/types"
)

func FetchProjectsDeployments(ctx context.Context, glab *gitlab.Client, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) ([]types.Deployment, error) {
	bundles, err := fanOut(ctx, glab, projectIds,
		func(ctx context.Context, glab *gitlab.Client, pid int64) ([]types.Deployment, error) {
			return FetchProjectDeployments(ctx, glab, pid, updatedAfter, updatedBefore)
		},
		fanOutOptions{},
	)

	// Flatten per-project slices returned by fanOut into a single slice.
	return slices.Concat(bundles...), err
}

func FetchProjectDeployments(ctx context.Context, glab *gitlab.Client, projectId int64, updatedAfter *time.Time, updatedBefore *time.Time) ([]types.Deployment, error) {
	opt := rest.GetProjectDeploymentsOptions{
		UpdatedAfter:  updatedAfter,
		UpdatedBefore: updatedBefore,
	}

	ds, err := glab.Rest.GetProjectDeployments(ctx, projectId, opt)
	if err != nil {
		return nil, err
	}

	deployments := make([]types.Deployment, 0, len(ds))
	for _, d := range ds {
		deployments = append(deployments, rest.ConvertDeployment(d))
	}

	return deployments, nil
}
