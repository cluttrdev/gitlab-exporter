package gitlab

import (
	"context"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

type GetProjectResult struct {
	Project *models.Project
	Error   error
}

func (c *Client) GetProject(ctx context.Context, id interface{}) <-chan GetProjectResult {
	ch := make(chan GetProjectResult)

	go func() {
		defer close(ch)

		opts := &_gitlab.GetProjectOptions{
			Statistics: &[]bool{true}[0],
		}

		c.RLock()
		p, _, err := c.client.Projects.GetProject(id, opts, _gitlab.WithContext(ctx))
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
