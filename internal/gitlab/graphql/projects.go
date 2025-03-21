package graphql

import (
	"context"
	"fmt"
	"time"

	"go.cluttr.dev/gitlab-exporter/internal/types"
)

type ListProjectsResult struct {
	ProjectFields ProjectFields

	PipelinesCount     int
	MergeRequestsCount int
}

type ProjectFields struct {
	ProjectReferenceFields
	Namespace NamespaceReferenceFields

	ProjectFieldsCore
}

func ConvertProject(pf ProjectFields) (types.Project, error) {
	var (
		id, namespaceId int64
		err             error
	)
	if id, err = ParseId(pf.Id, GlobalIdProjectPrefix); err != nil {
		return types.Project{}, fmt.Errorf("parse project id: %w", err)
	}
	if namespaceId, err = parseNamespaceId(pf.Namespace.Id); err != nil {
		return types.Project{}, fmt.Errorf("parse namespace id: %w", err)
	}

	p := types.Project{
		Id: id,
		Namspace: types.NamespaceReference{
			Id:       namespaceId,
			FullPath: pf.Namespace.FullPath,
		},

		Name:        pf.Name,
		FullName:    pf.NameWithNamespace,
		Path:        pf.Path,
		FullPath:    pf.FullPath,
		Description: valOrZero(pf.Description),

		CreatedAt:      pf.CreatedAt,
		UpdatedAt:      pf.UpdatedAt,
		LastActivityAt: pf.LastActivityAt,

		Statistics: types.ProjectStatistics{
			// ...
			StarCount:  int64(pf.StarCount),
			ForksCount: int64(pf.ForksCount),
		},

		Archived:   valOrZero(pf.Archived),
		Visibility: valOrZero(pf.Visibility),

		DefaultBranch: valOrZero(valOrZero(pf.Repository).RootRef),
	}

	if pf.Statistics != nil {
		p.Statistics.ContainerRegistrySize = int64(valOrZero(pf.Statistics.ContainerRegistrySize))
		p.Statistics.JobArtifactsSize = int64(pf.Statistics.BuildArtifactsSize)
		p.Statistics.LfsObjectsSize = int64(pf.Statistics.LfsObjectsSize)
		p.Statistics.PackagesSize = int64(pf.Statistics.PackagesSize)
		p.Statistics.PipelineArtifactsSize = int64(valOrZero(pf.Statistics.PipelineArtifactsSize))
		p.Statistics.RepositorySize = int64(pf.Statistics.RepositorySize)
		p.Statistics.SnippetsSize = int64(valOrZero(pf.Statistics.SnippetsSize))
		p.Statistics.StorageSize = int64(pf.Statistics.StorageSize)
		p.Statistics.UploadsSize = int64(valOrZero(pf.Statistics.UploadsSize))
		p.Statistics.WikiSize = int64(valOrZero(pf.Statistics.WikiSize))

		p.Statistics.CommitCount = int64(pf.Statistics.CommitCount)
	}

	return p, nil
}

func (c *Client) ListProjects(
	ctx context.Context,
	ids []string,
	updatedAfter *time.Time,
	updatedBefore *time.Time,
	yield func([]ListProjectsResult) bool,
) error {
	var endCursor *string

	for {
		resp, err := getProjects(ctx, c.client, ids, updatedAfter, updatedBefore, endCursor)
		if err != nil {
			return err
		}

		results := make([]ListProjectsResult, 0, len(resp.Projects.Nodes))
		for _, p := range resp.Projects.Nodes {
			results = append(results, ListProjectsResult{
				ProjectFields: ProjectFields{
					ProjectReferenceFields: p.ProjectReferenceFields,
					Namespace:              valOrZero(p.Namespace).NamespaceReferenceFields,

					ProjectFieldsCore: p.ProjectFieldsCore,
				},

				PipelinesCount:     p.Pipelines.Count,
				MergeRequestsCount: p.MergeRequests.Count,
			})
		}

		if !yield(results) {
			break
		}

		if !resp.Projects.PageInfo.HasNextPage {
			break
		}

		endCursor = resp.Projects.PageInfo.EndCursor
	}

	return nil
}

func (c *Client) GetProjects(
	ctx context.Context,
	ids []string,
	updatedAfter *time.Time,
	updatedBefore *time.Time,
) ([]ListProjectsResult, error) {
	var projects []ListProjectsResult

	err := c.ListProjects(ctx, ids, updatedAfter, updatedBefore, func(ps []ListProjectsResult) bool {
		projects = append(projects, ps...)
		return true
	})
	if err != nil {
		return nil, err
	}

	return projects, nil
}
