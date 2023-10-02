package gitlabclient

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	gogitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

func parseID(id interface{}) (string, error) {
	switch v := id.(type) {
	case int:
		return strconv.Itoa(v), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("invalid id type %#v", id)
	}
}

func pathEscape(path string) string {
	return strings.ReplaceAll(url.PathEscape(path), ".", "%2E")
}

type GetProjectResult struct {
	Project *models.Project
	Error   error
}

func (c *Client) GetProject(ctx context.Context, id string) <-chan GetProjectResult {
	ch := make(chan GetProjectResult)

	go func() {
		defer close(ch)

		opts := &gogitlab.GetProjectOptions{
			Statistics: &[]bool{true}[0],
		}

		c.RLock()
		p, _, err := c.client.Projects.GetProject(id, opts, gogitlab.WithContext(ctx))
		c.RUnlock()
		if err != nil {
			ch <- GetProjectResult{
				Error: err,
			}
			return
		}

		ch <- GetProjectResult{
			Project: models.NewProject(p),
		}
	}()

	return ch
}
