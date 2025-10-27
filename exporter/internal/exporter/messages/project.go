package messages

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.cluttr.dev/gitlab-exporter/internal/types"
	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func NewProjectReference(project types.ProjectReference) *typespb.ProjectReference {
	return &typespb.ProjectReference{
		Id:       project.Id,
		FullPath: project.FullPath,
	}
}

func NewProject(p types.Project) *typespb.Project {
	return &typespb.Project{
		Id: p.Id,
		Namespace: &typespb.NamespaceReference{
			Id:       p.Namspace.Id,
			FullPath: p.Namspace.FullPath,
		},

		Name:     p.Name,
		FullName: p.FullName,
		Path:     p.Path,
		FullPath: p.FullPath,

		Timestamps: &typespb.ProjectTimestamps{
			CreatedAt:      timestamppb.New(valOrZero(p.CreatedAt)),
			UpdatedAt:      timestamppb.New(valOrZero(p.UpdatedAt)),
			LastActivityAt: timestamppb.New(valOrZero(p.LastActivityAt)),
		},

		Statistics: convertProjectStatistics(p.Statistics),

		Description: p.Description,

		Archived:   p.Archived,
		Visibility: string(p.Visibility),

		DefaultBranch: p.DefaultBranch,
	}
}

func convertProjectStatistics(stats types.ProjectStatistics) *typespb.ProjectStatistics {
	return &typespb.ProjectStatistics{
		JobArtifactsSize:      stats.JobArtifactsSize,
		ContainerRegistrySize: stats.ContainerRegistrySize,
		LfsObjectsSize:        stats.LfsObjectsSize,
		PackagesSize:          stats.PackagesSize,
		PipelineArtifactsSize: stats.PipelineArtifactsSize,
		RepositorySize:        stats.RepositorySize,
		SnippetsSize:          stats.SnippetsSize,
		StorageSize:           stats.StorageSize,
		UploadsSize:           stats.UploadsSize,
		WikiSize:              stats.WikiSize,

		StarsCount:  stats.StarCount,
		ForksCount:  stats.ForksCount,
		CommitCount: stats.CommitCount,
		// OpenIssuesCount: 0,
	}
}
