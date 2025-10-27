package tasks

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.cluttr.dev/gitlab-exporter/internal/gitlab"
	"go.cluttr.dev/gitlab-exporter/internal/gitlab/graphql"
	"go.cluttr.dev/gitlab-exporter/internal/types"
)

func FetchRunners(ctx context.Context, glab *gitlab.Client) ([]types.Runner, error) {
	runnerFields, err := glab.GraphQL.GetRunners(ctx)
	if errors.Is(err, context.Canceled) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("get runners: %w", err)
	}

	runners := make([]types.Runner, 0, len(runnerFields))
	for _, rf := range runnerFields {
		r, err := graphql.ConvertRunner(rf)
		if err != nil {
			slog.Error("error converting runner fields",
				slog.String("error", err.Error()),
				slog.String("id", rf.Id),
			)
			continue
		}
		runners = append(runners, r)
	}

	return runners, nil
}
