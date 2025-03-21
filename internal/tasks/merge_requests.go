package tasks

import (
	"context"
	"fmt"
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
	if err != nil {
		return nil, fmt.Errorf("get projects merge requests: %w", err)
	}

	mergeRequests := make([]types.MergeRequest, 0, len(mergeRequestsFields))
	for _, mrf := range mergeRequestsFields {
		mr, err := graphql.ConvertMergeRequest(mrf)
		if err != nil {
			return nil, fmt.Errorf("convert merge request fields: %w", err)
		}
		mergeRequests = append(mergeRequests, mr)
	}

	return mergeRequests, nil
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
	if err != nil {
		return nil, fmt.Errorf("get projects merge requests notes: %w", err)
	}

	mergeRequestNoteEvents := make([]types.MergeRequestNoteEvent, 0, len(mergeRequestsNotesFields))
	for _, nf := range mergeRequestsNotesFields {
		ne, err := graphql.ConvertMergeRequestNoteEvent(nf)
		if err != nil {
			return nil, fmt.Errorf("convert merge request note fields: %w", err)
		}
		mergeRequestNoteEvents = append(mergeRequestNoteEvents, ne)
	}

	return mergeRequestNoteEvents, nil
}
