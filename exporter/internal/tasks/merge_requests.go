package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/exporter/internal/types"
)

func FetchProjectsMergeRequests(ctx context.Context, glab *gitlab.Client, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) ([]types.MergeRequest, []types.MergeRequestCommit, error) {
	gids := make([]string, 0, len(projectIds))
	for _, id := range projectIds {
		gids = append(gids, graphql.FormatId(id, graphql.GlobalIdProjectPrefix))
	}

	opts := graphql.GetMergeRequestsOptions{
		UpdatedAfter:  updatedAfter,
		UpdatedBefore: updatedBefore,
	}

	mergeRequestsFields, err := glab.GraphQL.GetProjectsMergeRequests(ctx, gids, opts)
	if errors.Is(err, context.Canceled) {
		return nil, nil, err
	} else if err != nil {
		err = fmt.Errorf("get projects merge requests: %w", err)
	}

	mergeRequests := make([]types.MergeRequest, 0, len(mergeRequestsFields))
	mergeRequestCommits := []types.MergeRequestCommit{}
	for _, mrf := range mergeRequestsFields {
		mr, err := graphql.ConvertMergeRequest(mrf)
		if err != nil {
			slog.Error("error converting merge request fields",
				slog.String("id", mrf.Id),
				slog.String("projectId", mrf.Project.Id),
				slog.String("error", err.Error()),
			)
			continue
		}
		mergeRequests = append(mergeRequests, mr)

		if len(mrf.Commits) > 0 {
			mrRef := types.MergeRequestReference{
				Id:      mr.Id,
				Iid:     mr.Iid,
				Project: mr.Project,
			}
			for _, cf := range mrf.Commits {
				mrc, err := graphql.ConvertMergeRequestCommit(mrRef, cf)
				if err != nil {
					slog.Error("error converting merge request commit fields",
						slog.String("commit.sha", cf.Sha),
						slog.String("merge_request.iid", mrf.Iid),
						slog.String("merge_request.project.id", mrf.Project.Id),
						slog.String("error", err.Error()),
					)
					continue
				}
				mergeRequestCommits = append(mergeRequestCommits, mrc)
			}
		}
	}

	return mergeRequests, mergeRequestCommits, err
}

func FetchProjectsMergeRequestsNotes(ctx context.Context, glab *gitlab.Client, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) ([]types.MergeRequestNoteEvent, error) {
	gids := make([]string, 0, len(projectIds))
	for _, id := range projectIds {
		gids = append(gids, graphql.FormatId(id, graphql.GlobalIdProjectPrefix))
	}

	opts := graphql.GetMergeRequestsOptions{
		UpdatedAfter:  updatedAfter,
		UpdatedBefore: updatedBefore,
	}

	mergeRequestsNotesFields, err := glab.GraphQL.GetProjectsMergeRequestsNotes(ctx, gids, opts)
	if errors.Is(err, context.Canceled) {
		return nil, err
	} else if err != nil {
		err = fmt.Errorf("get projects merge requests notes: %w", err)
	}

	mergeRequestNoteEvents := make([]types.MergeRequestNoteEvent, 0, len(mergeRequestsNotesFields))
	for _, nf := range mergeRequestsNotesFields {
		ne, err := graphql.ConvertMergeRequestNoteEvent(nf)
		if err != nil {
			slog.Error("error converting merge request note fields",
				slog.String("id", nf.Id),
				slog.String("mrIid", nf.MergeRequest.Iid),
				slog.String("projectId", nf.Project.Id),
				slog.String("error", err.Error()),
			)
			continue
		}
		mergeRequestNoteEvents = append(mergeRequestNoteEvents, ne)
	}

	return mergeRequestNoteEvents, err
}
