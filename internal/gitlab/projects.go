package gitlab

import (
	"context"
	"fmt"
	"strings"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/types"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

func (c *Client) GetProject(ctx context.Context, id int64) (*typespb.Project, error) {
	opt := gitlab.GetProjectOptions{
		Statistics: gitlab.Ptr(true),
	}

	p, _, err := c.client.Projects.GetProject(int(id), &opt, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return types.ConvertProject(p), nil
}

type ListNamespaceProjectsOptions struct {
	Kind string

	// both
	gitlab.ListOptions

	Visibility *gitlab.VisibilityValue

	// only groups
	WithShared       *bool
	IncludeSubgroups *bool
}

func (c *Client) ListNamespaceProjects(ctx context.Context, id interface{}, opt ListNamespaceProjectsOptions, yield func(projects []*gitlab.Project) bool) error {
	kind := strings.ToLower(opt.Kind)
	if !(strings.EqualFold(kind, "user") || strings.EqualFold(kind, "group")) {
		n, _, err := c.client.Namespaces.GetNamespace(id, gitlab.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("error determining namespace kind: %w", err)
		}
		kind = n.Kind
	}

	if kind == "user" {
		opts := gitlab.ListProjectsOptions{
			ListOptions: opt.ListOptions,

			Visibility: opt.Visibility,
		}
		return c.ListUserProjects(ctx, id, opts, yield)
	} else if kind == "group" {
		opts := gitlab.ListGroupProjectsOptions{
			ListOptions: opt.ListOptions,

			Visibility: opt.Visibility,

			WithShared:       opt.WithShared,
			IncludeSubGroups: opt.IncludeSubgroups,
		}
		return c.ListGroupProjects(ctx, id, opts, yield)
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
