package tasks

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/cluttrdev/gitlab-exporter/internal/gitlab"
	"github.com/cluttrdev/gitlab-exporter/internal/gitlab/rest"
	"github.com/cluttrdev/gitlab-exporter/internal/types"
)

func FetchProjectsDeployments(ctx context.Context, glab *gitlab.Client, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) ([]types.Deployment, error) {
	type result struct {
		deployments []types.Deployment
		err         error
	}

	var (
		deployments []types.Deployment

		wg      sync.WaitGroup
		results = make(chan result)
	)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, projectId := range projectIds {
			if err := glab.Acquire(ctx, 1); err != nil {
				slog.Error("failed to acquire gitlab client", "error", err)
				continue
			}
			wg.Add(1)
			go func() {
				defer glab.Release(1)
				defer wg.Done()

				ds, err := FetchProjectDeployments(ctx, glab, projectId, updatedAfter, updatedBefore)
				results <- result{
					deployments: ds,
					err:         err,
				}
			}()
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	var errs error
loop:
	for {
		select {
		case <-done:
			break loop
		case r := <-results:
			if r.err != nil {
				errs = errors.Join(errs, r.err)
			} else {
				deployments = append(deployments, r.deployments...)
			}
		}
	}

	return deployments, errs
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
