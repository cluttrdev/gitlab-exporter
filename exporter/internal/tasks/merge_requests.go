package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/internal/types"
)

func FetchProjectsMergeRequests(ctx context.Context, glab *gitlab.Client, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) ([]types.MergeRequest, error) {
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
		return nil, err
	} else if err != nil {
		err = fmt.Errorf("get projects merge requests: %w", err)
	}

	mergeRequests := make([]types.MergeRequest, 0, len(mergeRequestsFields))
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
	}

	return mergeRequests, err
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
