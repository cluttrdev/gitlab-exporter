package gitlab

import (
	"context"
	"fmt"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/types"
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
			projects = append(projects, types.ConvertProject(p))
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

	return types.ConvertProject(p), nil
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
			projects = append(projects, types.ConvertProject(p))
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
			projects = append(projects, types.ConvertProject(p))
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return projects, nil
}
