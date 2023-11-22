package gitlab

import (
	"context"
	"fmt"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

func (c *Client) GetJobSections(ctx context.Context, projectID int64, jobID int64) ([]*models.Section, error) {
	sections := []*models.Section{}
	for r := range c.ListJobSections(ctx, projectID, jobID) {
		if r.Error != nil {
			return nil, fmt.Errorf("[gitlab.Client.GetSections] %w", r.Error)
		}
		sections = append(sections, r.Section)
	}
	return sections, nil
}

type ListJobSectionsResult struct {
	Section *models.Section
	Error   error
}

func (c *Client) ListJobSections(ctx context.Context, projectID int64, jobID int64) <-chan ListJobSectionsResult {
	ch := make(chan ListJobSectionsResult)

	go func() {
		defer close(ch)

		c.RLock()
		job, _, err := c.client.Jobs.GetJob(int(projectID), int(jobID), _gitlab.WithContext(ctx))
		c.RUnlock()
		if err != nil {
			ch <- ListJobSectionsResult{
				Error: err,
			}
			return
		}

		c.RLock()
		trace, _, err := c.client.Jobs.GetTraceFile(int(projectID), int(jobID), _gitlab.WithContext(ctx))
		c.RUnlock()
		if err != nil {
			ch <- ListJobSectionsResult{
				Error: err,
			}
			return
		}

		sections, err := models.ParseSections(trace)
		if err != nil {
			ch <- ListJobSectionsResult{
				Error: err,
			}
			return
		}

		for secnum, section := range sections {
			section.ID = int64(job.ID*1000 + secnum)
			section.Job.ID = int64(job.ID)
			section.Job.Name = job.Name
			section.Job.Status = job.Status
			section.Pipeline.ID = int64(job.Pipeline.ID)
			section.Pipeline.ProjectID = int64(job.Pipeline.ProjectID)
			section.Pipeline.Ref = job.Pipeline.Ref
			section.Pipeline.Sha = job.Pipeline.Sha
			section.Pipeline.Status = job.Pipeline.Status

			ch <- ListJobSectionsResult{
				Section: section,
			}
		}
	}()

	return ch
}
