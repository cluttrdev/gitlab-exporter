package graphql

import (
	"context"
	"fmt"

	"go.cluttr.dev/gitlab-exporter/internal/types"
)

type RunnerFields struct {
	RunnerReferenceFields
	RunnerFieldsCore
}

func ConvertRunner(rf RunnerFields) (types.Runner, error) {
	var (
		id  int64
		err error
	)
	if id, err = ParseId(rf.Id, GlobalIdRunnerPrefix); err != nil {
		return types.Runner{}, fmt.Errorf("parse runner id: %w", err)
	}

	runner := types.Runner{
		Id:       id,
		ShortSha: valOrZero(rf.ShortSha),

		Description: valOrZero(rf.Description),
		RunnerType:  convertRunnerType(rf.RunnerFieldsCore.RunnerType),
		TagList:     rf.TagList,

		Status: convertRunnerStatus(rf.Status),

		Locked: valOrZero(rf.Locked),
		Paused: rf.Paused,

		RunProtected: convertRunnerAccessLevel(rf.AccessLevel) == types.RunnerAccessLevelRefProtected,
		RunUntagged:  rf.RunUntagged,

		ContactedAt: rf.ContactedAt,
		CreatedAt:   rf.CreatedAt,

		// CreatedBy: nil,
	}

	if rf.CreatedBy != nil {
		createdBy, err := convertUserReference(rf.CreatedBy)
		if err != nil {
			return types.Runner{}, fmt.Errorf("convert createBy user reference: %w", err)
		}
		runner.CreatedBy = createdBy
	}

	return runner, nil
}

func convertRunnerType(rt CiRunnerType) types.RunnerType {
	switch rt {
	case CiRunnerTypeInstanceType:
		return types.RunnerTypeInstance
	case CiRunnerTypeGroupType:
		return types.RunnerTypeGroup
	case CiRunnerTypeProjectType:
		return types.RunnerTypeProject
	}

	return types.RunnerTypeUnknown
}

func convertRunnerStatus(s CiRunnerStatus) types.RunnerStatus {
	switch s {
	case CiRunnerStatusOnline:
		return types.RunnerStatusOnline
	case CiRunnerStatusOffline:
		return types.RunnerStatusOffline
	case CiRunnerStatusStale:
		return types.RunnerStatusStale
	case CiRunnerStatusNeverContacted:
		return types.RunnerStatusNeverContacted
	}

	return types.RunnerStatusUnknown
}

func convertRunnerAccessLevel(al CiRunnerAccessLevel) types.RunnerAccessLevel {
	switch al {
	case CiRunnerAccessLevelNotProtected:
		return types.RunnerAccessLevelNotProtected
	case CiRunnerAccessLevelRefProtected:
		return types.RunnerAccessLevelRefProtected
	}

	return types.RunnerAccessLevelUnknown
}

func (c *Client) GetRunners(ctx context.Context) ([]RunnerFields, error) {
	return c.getRunners(ctx, getRunnersOptions{})
}

type getRunnersOptions struct {
	endCursor *string
}

func (c *Client) getRunners(ctx context.Context, opts getRunnersOptions) ([]RunnerFields, error) {
	var (
		runners []RunnerFields

		data *getRunnersResponse
		err  error
	)

	for {
		data, err = getRunners(ctx, c.client, opts.endCursor)
		err = handleError(err, "getRunners")
		if err != nil {
			break
		}

		if data.Runners == nil {
			break
		}

		for _, runner_ := range data.Runners.Nodes {
			runner := RunnerFields{
				RunnerReferenceFields: runner_.RunnerReferenceFields,
				RunnerFieldsCore:      runner_.RunnerFieldsCore,
			}

			runners = append(runners, runner)
		}

		if !data.Runners.PageInfo.HasNextPage {
			break
		}

		opts.endCursor = data.Runners.PageInfo.EndCursor
	}

	return runners, err
}
