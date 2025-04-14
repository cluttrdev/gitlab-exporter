package tasks

import (
	"context"
	"fmt"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/internal/types"
)

func FetchProjectsIssues(ctx context.Context, glab *gitlab.Client, projectIds []int64, updatedAfter *time.Time, updatedBefore *time.Time) ([]types.Issue, error) {
	gids := make([]string, 0, len(projectIds))
	for _, id := range projectIds {
		gids = append(gids, graphql.FormatId(id, graphql.GlobalIdProjectPrefix))
	}

	opts := graphql.GetProjectsIssuesOptions{
		TimeRangeOptions: graphql.TimeRangeOptions{
			UpdatedAfter:  updatedAfter,
			UpdatedBefore: updatedBefore,
		},
	}

	issueFields, err := glab.GraphQL.GetProjectsIssues(ctx, gids, opts)
	if err != nil {
		return nil, fmt.Errorf("get projects issues: %w", err)
	}

	issues := make([]types.Issue, 0, len(issueFields))
	for _, isf := range issueFields {
		iss, err := graphql.ConvertIssue(isf)
		if err != nil {
			return nil, fmt.Errorf("convert issue fields: %w", err)
		}
		issues = append(issues, iss)
	}

	return issues, nil
}
