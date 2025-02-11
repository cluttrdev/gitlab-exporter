package rest

import (
	"context"
	"fmt"
	"strings"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type Project = gitlab.Project

func (c *Client) GetProject(ctx context.Context, pid interface{}) (*gitlab.Project, error) {
	opt := gitlab.GetProjectOptions{
		Statistics: gitlab.Ptr(true),
	}

	p, _, err := c.client.Projects.GetProject(pid, &opt, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return p, nil
}

type ListNamespaceProjectsOptions struct {
	Kind string

	// both
	gitlab.ListOptions

	Visibility *string

	// only groups
	WithShared       *bool
	IncludeSubgroups *bool
}

func (c *Client) ListNamespaceProjects(ctx context.Context, nid interface{}, opt ListNamespaceProjectsOptions, yield func(projects []*gitlab.Project) bool) error {
	kind := strings.ToLower(opt.Kind)
	if !(strings.EqualFold(kind, "user") || strings.EqualFold(kind, "group")) {
		n, _, err := c.client.Namespaces.GetNamespace(nid, gitlab.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("error determining namespace kind: %w", err)
		}
		kind = n.Kind
	}

	if kind == "user" {
		opts := gitlab.ListProjectsOptions{
			ListOptions: opt.ListOptions,

			Visibility: (*gitlab.VisibilityValue)(opt.Visibility),
		}
		return c.ListUserProjects(ctx, nid, opts, yield)
	} else if kind == "group" {
		opts := gitlab.ListGroupProjectsOptions{
			ListOptions: opt.ListOptions,

			Visibility: (*gitlab.VisibilityValue)(opt.Visibility),

			WithShared:       opt.WithShared,
			IncludeSubGroups: opt.IncludeSubgroups,
		}
		return c.ListGroupProjects(ctx, nid, opts, yield)
	}
	return fmt.Errorf("invalid namespace kind: %v", kind)
}

func (c *Client) ListUserProjects(ctx context.Context, uid interface{}, opt gitlab.ListProjectsOptions, yield func(projects []*gitlab.Project) bool) error {
	opt.ListOptions.Pagination = ""
	if opt.ListOptions.OrderBy == "" {
		opt.ListOptions.OrderBy = "updated_at"
	}
	if opt.ListOptions.Sort == "" {
		opt.ListOptions.Sort = "desc"
	}

	options := []gitlab.RequestOptionFunc{
		gitlab.WithContext(ctx),
	}

	for {
		ps, resp, err := c.client.Projects.ListUserProjects(uid, &opt, options...)
		if err != nil {
			return err
		}

		if !yield(ps) {
			break
		}

		if opt.ListOptions.Pagination == "keyset" {
			if resp.NextLink == "" {
				break
			}

			options = []gitlab.RequestOptionFunc{
				gitlab.WithContext(ctx),
				gitlab.WithKeysetPaginationParameters(resp.NextLink),
			}
		} else {
			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}
	}

	return nil
}

func (c *Client) ListGroupProjects(ctx context.Context, gid interface{}, opt gitlab.ListGroupProjectsOptions, yield func(projects []*gitlab.Project) bool) error {
	opt.ListOptions.Pagination = "keyset"
	if opt.ListOptions.OrderBy == "" {
		opt.ListOptions.OrderBy = "updated_at"
	}
	if opt.ListOptions.Sort == "" {
		opt.ListOptions.Sort = "desc"
	}

	options := []gitlab.RequestOptionFunc{
		gitlab.WithContext(ctx),
	}

	for {
		ps, resp, err := c.client.Groups.ListGroupProjects(gid, &opt, options...)
		if err != nil {
			return err
		}

		if !yield(ps) {
			break
		}

		if resp.NextLink == "" {
			break
		}

		options = []gitlab.RequestOptionFunc{
			gitlab.WithContext(ctx),
			gitlab.WithKeysetPaginationParameters(resp.NextLink),
		}
	}

	return nil
}
