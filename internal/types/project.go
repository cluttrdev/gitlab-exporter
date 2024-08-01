package types

import (
	"github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

func ConvertProject(p *gitlab.Project) *typespb.Project {
	return &typespb.Project{
		Id:                int64(p.ID),
		Name:              p.Name,
		NameWithNamespace: p.NameWithNamespace,
		Path:              p.Path,
		PathWithNamespace: p.PathWithNamespace,

		CreatedAt:      ConvertTime(p.CreatedAt),
		LastActivityAt: ConvertTime(p.LastActivityAt),

		Namespace: convertProjectNamespace(p.Namespace),
		Owner:     convertUser(p.Owner),
		CreatorId: int64(p.CreatorID),

		Topics:          p.Topics,
		ForksCount:      int64(p.ForksCount),
		StarsCount:      int64(p.StarCount),
		Statistics:      convertProjectStatistics(p.Statistics),
		OpenIssuesCount: int64(p.OpenIssuesCount),

		Description: p.Description,

		EmptyRepo: p.EmptyRepo,
		Archived:  p.Archived,

		DefaultBranch: p.DefaultBranch,
		Visibility:    string(p.Visibility),
		WebUrl:        p.WebURL,
	}
}

func convertProjectNamespace(n *gitlab.ProjectNamespace) *typespb.ProjectNamespace {
	if n == nil {
		return nil
	}
	return &typespb.ProjectNamespace{
		Id:       int64(n.ID),
		Name:     n.Name,
		Kind:     n.Kind,
		Path:     n.Path,
		FullPath: n.FullPath,
		ParentId: int64(n.ParentID),

		AvatarUrl: n.AvatarURL,
		WebUrl:    n.WebURL,
	}
}

func convertProjectStatistics(stats *gitlab.Statistics) *typespb.ProjectStatistics {
	if stats == nil {
		return nil
	}
	return &typespb.ProjectStatistics{
		CommitCount:      stats.CommitCount,
		StorageSize:      stats.StorageSize,
		RepositorySize:   stats.RepositorySize,
		WikiSize:         stats.WikiSize,
		LfsObjectsSize:   stats.LFSObjectsSize,
		JobArtifactsSize: stats.JobArtifactsSize,
		PackagesSize:     stats.PackagesSize,
		SnippetsSize:     stats.SnippetsSize,
		UploadsSize:      stats.UploadsSize,
	}
}
