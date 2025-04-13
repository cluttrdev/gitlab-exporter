package types

import (
	"time"
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

	DefaultBranch string
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
