package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/time/rate"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

type Client struct {
	sync.RWMutex
	client *_gitlab.Client
}

type ClientConfig struct {
	URL   string
	Token string

	RateLimit float64
}

func NewGitLabClient(cfg ClientConfig) (*Client, error) {
	var client Client

	if err := client.Configure(cfg); err != nil {
		return nil, err
	}

	return &client, nil
}

func (c *Client) Configure(cfg ClientConfig) error {
	opts := []_gitlab.ClientOptionFunc{
		_gitlab.WithBaseURL(cfg.URL),
	}

	if cfg.RateLimit > 0 {
		limit := rate.Limit(cfg.RateLimit * 0.66)
		burst := cfg.RateLimit * 0.33
		limiter := rate.NewLimiter(limit, int(burst))

		opts = append(opts, _gitlab.WithCustomLimiter(limiter))
	}

	client, err := _gitlab.NewOAuthClient(cfg.Token, opts...)
	if err != nil {
		return err
	}

	c.Lock()
	c.client = client
	c.Unlock()
	return nil
}

func (c *Client) CheckReadiness(ctx context.Context) error {
	const readinessEndpoint string = "version"

	req, err := c.client.NewRequest(
		http.MethodGet,
		readinessEndpoint,
		nil,
		[]_gitlab.RequestOptionFunc{_gitlab.WithContext(ctx)},
	)
	if err != nil {
		return err
	}

	res, err := c.client.Do(req, nil)
	if err != nil {
		return err
	}

	if res == nil {
		return fmt.Errorf("http error: empty response")
	} else if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http error: %d", res.StatusCode)
	}

	return nil
}

type GetPipelineHierarchyOptions struct {
	FetchSections bool
}

type GetPipelineHierarchyResult struct {
	PipelineHierarchy *models.PipelineHierarchy
	Error             error
}

func (c *Client) GetPipelineHierarchy(ctx context.Context, projectID int64, pipelineID int64, opt *GetPipelineHierarchyOptions) <-chan GetPipelineHierarchyResult {
	ch := make(chan GetPipelineHierarchyResult)

	go func() {
		defer close(ch)

		pipeline, err := c.GetPipeline(ctx, projectID, pipelineID)
		if err != nil {
			ch <- GetPipelineHierarchyResult{
				Error: err,
			}
			return
		}

		jobs := []*models.Job{}
		sections := []*models.Section{}
		for jr := range c.ListPipelineJobs(ctx, projectID, pipelineID) {
			if jr.Error != nil {
				ch <- GetPipelineHierarchyResult{
					Error: fmt.Errorf("[ListPipelineJobs] %w", jr.Error),
				}
				return
			}
			jobs = append(jobs, jr.Job)

			if opt.FetchSections {
				jobID := jr.Job.ID
				for sr := range c.ListJobSections(ctx, projectID, jobID) {
					if sr.Error != nil {
						ch <- GetPipelineHierarchyResult{
							Error: fmt.Errorf("[ListJobSections] %w", sr.Error),
						}
						return
					}
					sections = append(sections, sr.Section)
				}
			}
		}

		bridges := []*models.Bridge{}
		dps := []*models.PipelineHierarchy{}
		for br := range c.ListPipelineBridges(ctx, projectID, pipelineID) {
			if br.Error != nil {
				ch <- GetPipelineHierarchyResult{
					Error: fmt.Errorf("[ListPipelineBridges] %w", br.Error),
				}
				return
			}
			bridges = append(bridges, br.Bridge)

			dp := br.Bridge.DownstreamPipeline
			if dp == nil || dp.ID == 0 {
				continue
			}

			dpr := <-c.GetPipelineHierarchy(ctx, dp.ProjectID, dp.ID, opt)
			if dpr.Error != nil {
				ch <- GetPipelineHierarchyResult{
					Error: fmt.Errorf("[GetPipelineHierarchy] %w", dpr.Error),
				}
				return
			}
			dps = append(dps, dpr.PipelineHierarchy)
		}

		ch <- GetPipelineHierarchyResult{
			PipelineHierarchy: &models.PipelineHierarchy{
				Pipeline:            pipeline,
				Jobs:                jobs,
				Sections:            sections,
				Bridges:             bridges,
				DownstreamPipelines: dps,
			},
		}
	}()

	return ch
}
