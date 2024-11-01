package types

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type NamespaceReference struct {
	Id       int64
	FullPath string
}

type ProjectReference struct {
	Id       int64
	FullPath string
}

type Project struct {
	Id       int64
	Namspace NamespaceReference

	Name        string
	FullName    string
	Path        string
	FullPath    string
	Description string

	CreatedAt      *time.Time
	UpdatedAt      *time.Time
	LastActivityAt *time.Time

	Statistics ProjectStatistics

	Archived   bool
	Visibility string
}

type ProjectStatistics struct {
	JobArtifactsSize      int64
	ContainerRegistrySize int64
	LfsObjectsSize        int64
	PackagesSize          int64
	PipelineArtifactsSize int64
	RepositorySize        int64
	SnippetsSize          int64
	StorageSize           int64
	UploadsSize           int64
	WikiSize              int64

	CommitCount int64
	StarCount   int64
	ForksCount  int64
}

func ConvertProjectReference(project ProjectReference) *typespb.ProjectReference {
	return &typespb.ProjectReference{
		Id:       project.Id,
		FullPath: project.FullPath,
	}
}

func ConvertProject(p Project) *typespb.Project {
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
	}
}

func convertProjectStatistics(stats ProjectStatistics) *typespb.ProjectStatistics {
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
