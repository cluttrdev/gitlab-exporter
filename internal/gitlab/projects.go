package gitlab

import (
	"context"
	"fmt"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type ListProjectsOptions = gitlab.ListProjectsOptions
type ListGroupProjectsOptions = gitlab.ListGroupProjectsOptions
type VisibilityValue = gitlab.VisibilityValue

func (c *Client) ListProjects(ctx context.Context, opt ListProjectsOptions) ([]*typespb.Project, error) {
	var projects []*typespb.Project

	for {
		ps, resp, err := c.client.Projects.ListProjects(&opt, gitlab.WithContext(ctx))
		if err != nil {
			return projects, err
		}

		for _, p := range ps {
			projects = append(projects, convertProject(p))
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return projects, nil
}

func (c *Client) GetProject(ctx context.Context, id int64) (*typespb.Project, error) {
	opt := gitlab.GetProjectOptions{
		Statistics: Ptr(true),
	}

	p, _, err := c.client.Projects.GetProject(int(id), &opt, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return convertProject(p), nil
}

type ListNamespaceProjectsOptions struct {
	gitlab.ListProjectsOptions

	Kind             string
	WithShared       bool
	IncludeSubgroups bool
}

func (c *Client) ListNamespaceProjects(ctx context.Context, id interface{}, opt ListNamespaceProjectsOptions) ([]*typespb.Project, error) {
	kind := strings.ToLower(opt.Kind)
	if !(strings.EqualFold(kind, "user") || strings.EqualFold(kind, "group")) {
		n, _, err := c.client.Namespaces.GetNamespace(id, gitlab.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("error determining namespace kind: %w", err)
		}
		kind = n.Kind
	}

	if kind == "user" {
		return c.ListUserProjects(ctx, id, opt.ListProjectsOptions)
	} else if kind == "group" {
		return c.ListGroupProjects(ctx, id, gitlab.ListGroupProjectsOptions{
			ListOptions:      opt.ListOptions,
			WithShared:       &opt.WithShared,
			IncludeSubGroups: &opt.IncludeSubgroups,
		})
	}
	return nil, fmt.Errorf("invalid namespace kind: %v", kind)
}

func (c *Client) ListUserProjects(ctx context.Context, uid interface{}, opt ListProjectsOptions) ([]*typespb.Project, error) {
	var projects []*typespb.Project

	for {
		ps, resp, err := c.client.Projects.ListUserProjects(uid, &opt, gitlab.WithContext(ctx))
		if err != nil {
			return projects, err
		}

		for _, p := range ps {
			projects = append(projects, convertProject(p))
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return projects, nil
}

func (c *Client) ListGroupProjects(ctx context.Context, gid interface{}, opt ListGroupProjectsOptions) ([]*typespb.Project, error) {
	var projects []*typespb.Project

	for {
		ps, resp, err := c.client.Groups.ListGroupProjects(gid, &opt, gitlab.WithContext(ctx))
		if err != nil {
			return projects, err
		}

		for _, p := range ps {
			projects = append(projects, convertProject(p))
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return projects, nil
}

func convertProject(p *gitlab.Project) *typespb.Project {
	return &typespb.Project{
		Id:                int64(p.ID),
		Name:              p.Name,
		NameWithNamespace: p.NameWithNamespace,
		Path:              p.Path,
		PathWithNamespace: p.PathWithNamespace,

		CreatedAt:      convertTime(p.CreatedAt),
		LastActivityAt: convertTime(p.LastActivityAt),

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
