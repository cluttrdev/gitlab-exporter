package rest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	gitlab "gitlab.com/gitlab-org/api/client-go"

	"github.com/cluttrdev/gitlab-exporter/internal/types"
)

type GitLabEnvironment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Tier string `json:"tier"`

	Project struct {
		ID                int    `json:"id"`
		NameWithNamespace string `json:"name_with_namespace"`
	} `json:"project"`
}

type GitLabDeployment struct {
	ID        int        `json:"id"`
	IID       int        `json:"iid"`
	Ref       string     `json:"ref"`
	SHA       string     `json:"sha"`
	Status    string     `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	User      struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Username string `json:"username"`
	} `json:"user"`
	Environment struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"environment"`
	Deployable struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Pipeline struct {
			ID        int `json:"id"`
			IID       int `json:"iid"`
			ProjectID int `json:"project_id"`
		} `json:"pipeline"`
		WebURL string `json:"web_url"`
	} `json:"deployable"`

	// --- internal
	environment *GitLabEnvironment
}

func (d *GitLabDeployment) getEnvironment() *GitLabEnvironment {
	if d.environment != nil {
		return d.environment
	}

	env := &GitLabEnvironment{
		ID:   d.Environment.ID,
		Name: d.Environment.Name,
		Tier: "",
	}
	env.Project.ID = d.Deployable.Pipeline.ProjectID
	env.Project.NameWithNamespace = d.deployableProjectFullPath()

	return env
}

func (d *GitLabDeployment) deployableProjectFullPath() string {
	u, err := url.Parse(d.Deployable.WebURL)
	if err != nil {
		return ""
	}

	fullPath, _, ok := strings.Cut(u.Path, "/-/jobs/")
	if !ok {
		return ""
	}

	return strings.TrimPrefix(fullPath, "/")
}

func ConvertDeployment(deployment *GitLabDeployment) types.Deployment {
	env := deployment.getEnvironment()

	project := types.ProjectReference{
		Id:       int64(env.Project.ID),
		FullPath: env.Project.NameWithNamespace,
	}

	return types.Deployment{
		Id:  int64(deployment.ID),
		Iid: int64(deployment.IID),

		Job: types.JobReference{
			Id:   int64(deployment.Deployable.ID),
			Name: deployment.Deployable.Name,

			Pipeline: types.PipelineReference{
				Id:  int64(deployment.Deployable.Pipeline.ID),
				Iid: int64(deployment.Deployable.Pipeline.IID),

				Project: project,
			},
		},

		Triggerer: types.UserReference{
			Id:       int64(deployment.User.ID),
			Username: deployment.User.Username,
			Name:     deployment.User.Name,
		},

		Environment: types.EnvironmentReference{
			Id:   int64(env.ID),
			Name: env.Name,
			Tier: env.Tier,

			Project: project,
		},

		CreatedAt:  deployment.CreatedAt,
		FinishedAt: deployment.UpdatedAt,
		UpdatedAt:  deployment.UpdatedAt,

		Status: deployment.Status,
		Ref:    deployment.Ref,
		Sha:    deployment.SHA,
	}
}

type GetProjectDeploymentsOptions struct {
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time
}

func (c *Client) GetProjectDeployments(ctx context.Context, projectId int64, opt GetProjectDeploymentsOptions) ([]*GitLabDeployment, error) {
	environments, err := getProjectEnvironments(c, ctx, projectId)
	if err != nil {
		return nil, err
	}
	envMap := make(map[int64]*GitLabEnvironment, len(environments))
	for _, env := range environments {
		envMap[int64(env.ID)] = env
	}

	var deployments []*GitLabDeployment

	opts := &gitlab.ListProjectDeploymentsOptions{
		ListOptions: gitlab.ListOptions{
			Pagination: "keyset",
			PerPage:    100,
			OrderBy:    "updated_at",
			Sort:       "desc",
		},

		UpdatedAfter:  opt.UpdatedAfter,
		UpdatedBefore: opt.UpdatedBefore,
	}

	options := []gitlab.RequestOptionFunc{
		gitlab.WithContext(ctx),
	}

	for {
		ds, resp, err := listProjectDeployments(c.client, int(projectId), opts, options...)
		if err != nil {
			return nil, err
		}

		deployments = append(deployments, ds...)

		if resp.NextLink == "" {
			break
		}

		options = []gitlab.RequestOptionFunc{
			gitlab.WithContext(ctx),
			gitlab.WithKeysetPaginationParameters(resp.NextLink),
		}
	}

	for i := 0; i < len(deployments); i++ {
		env, ok := envMap[int64(deployments[i].Environment.ID)]
		if !ok {
			continue
		}
		deployments[i].environment = env
	}

	return deployments, nil
}

func getProjectEnvironments(c *Client, ctx context.Context, projectId int64) ([]*GitLabEnvironment, error) {
	var environments []*GitLabEnvironment

	opts := &gitlab.ListEnvironmentsOptions{
		ListOptions: gitlab.ListOptions{
			Pagination: "keyset",
			PerPage:    100,
		},
	}

	options := []gitlab.RequestOptionFunc{
		gitlab.WithContext(ctx),
	}

	for {
		envs, resp, err := listEnvironments(c.client, int(projectId), opts, options...)
		if err != nil {
			return nil, err
		}

		environments = append(environments, envs...)

		if resp.NextLink == "" {
			break
		}

		options = []gitlab.RequestOptionFunc{
			gitlab.WithContext(ctx),
			gitlab.WithKeysetPaginationParameters(resp.NextLink),
		}
	}
	return environments, nil
}

func listEnvironments(client *gitlab.Client, pid interface{}, opts *gitlab.ListEnvironmentsOptions, options ...gitlab.RequestOptionFunc) ([]*GitLabEnvironment, *gitlab.Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/environments", gitlab.PathEscape(project))

	req, err := client.NewRequest(http.MethodGet, u, opts, options)
	if err != nil {
		return nil, nil, err
	}

	var envs []*GitLabEnvironment
	resp, err := client.Do(req, &envs)
	if err != nil {
		return nil, resp, err
	}

	return envs, resp, nil
}

func listProjectDeployments(client *gitlab.Client, pid interface{}, opts *gitlab.ListProjectDeploymentsOptions, options ...gitlab.RequestOptionFunc) ([]*GitLabDeployment, *gitlab.Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/deployments", gitlab.PathEscape(project))

	req, err := client.NewRequest(http.MethodGet, u, opts, options)
	if err != nil {
		return nil, nil, err
	}

	var ds []*GitLabDeployment
	resp, err := client.Do(req, &ds)
	if err != nil {
		return nil, resp, err
	}

	return ds, resp, nil
}

// Helper function to accept and format both the project ID or name as project
// identifier for all API calls.
func parseID(id interface{}) (string, error) {
	switch v := id.(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("invalid ID type %#v, the ID must be an int or a string", id)
	}
}
