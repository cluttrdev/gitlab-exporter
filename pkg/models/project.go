package models

import (
	"time"

	gogitlab "github.com/xanzy/go-gitlab"
)

type Project struct {
	ID                int64             `json:"id"`
	Name              string            `json:"name"`
	NameWithNamespace string            `json:"name_with_namespace"`
	Path              string            `json:"path"`
	PathWithNamespace string            `json:"path_with_namespace"`
	CreatedAt         *time.Time        `json:"created_at"`
	LastActivityAt    *time.Time        `json:"last_activity_at"`
	WebURL            string            `json:"web_url"`
	Statistics        ProjectStatistics `json:"statistics"`
}

type ProjectStatistics struct {
	CommitCount           int64 `json:"commit_count"`
	StorageSize           int64 `json:"storage_size"`
	RepositorySize        int64 `json:"repository_size"`
	WikiSize              int64 `json:"wiki_size"`
	LFSObjectsSize        int64 `json:"lfs_objects_size"`
	JobArtifactsSize      int64 `json:"job_artifacts_size"`
	PipelineArtifactsSize int64 `json:"pipeline_artifacts_size"`
	PackagesSize          int64 `json:"packages_size"`
	SnippetsSize          int64 `json:"snippets_size"`
	UploadsSize           int64 `json:"uploads_size"`
}

func NewProject(p *gogitlab.Project) *Project {
	return &Project{
		ID:                int64(p.ID),
		Name:              p.Name,
		NameWithNamespace: p.NameWithNamespace,
		Path:              p.Path,
		PathWithNamespace: p.PathWithNamespace,
		CreatedAt:         p.CreatedAt,
		LastActivityAt:    p.LastActivityAt,
		WebURL:            p.WebURL,
		Statistics: ProjectStatistics{
			CommitCount:           p.Statistics.CommitCount,
			StorageSize:           p.Statistics.StorageSize,
			RepositorySize:        p.Statistics.RepositorySize,
			WikiSize:              p.Statistics.WikiSize,
			LFSObjectsSize:        p.Statistics.LFSObjectsSize,
			JobArtifactsSize:      p.Statistics.JobArtifactsSize,
			PipelineArtifactsSize: p.Statistics.PipelineArtifactsSize,
			PackagesSize:          p.Statistics.PackagesSize,
			SnippetsSize:          p.Statistics.SnippetsSize,
			UploadsSize:           p.Statistics.UploadsSize,
		},
	}
}
