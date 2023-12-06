package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"

	_gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/pkg/models"
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
	FetchSections   bool
	FetchJobMetrics bool
}

type GetPipelineHierarchyResult struct {
	PipelineHierarchy *models.PipelineHierarchy
	JobMetrics        []*models.JobMetric
	Error             error
}

func (c *Client) GetPipelineHierarchy(ctx context.Context, projectID int64, pipelineID int64, opt *GetPipelineHierarchyOptions) <-chan GetPipelineHierarchyResult {
	ch := make(chan GetPipelineHierarchyResult)

	go func() {
		defer close(ch)

		unixTime := func(ts int64) *time.Time {
			const nsec int64 = 0
			t := time.Unix(ts, nsec)
			return &t
		}

		pipeline, err := c.GetPipeline(ctx, projectID, pipelineID)
		if err != nil {
			ch <- GetPipelineHierarchyResult{
				Error: err,
			}
			return
		}

		jobs := []*models.Job{}
		sections := []*models.Section{}
		metrics := []*models.JobMetric{}
		for jr := range c.ListPipelineJobs(ctx, projectID, pipelineID) {
			if jr.Error != nil {
				ch <- GetPipelineHierarchyResult{
					Error: fmt.Errorf("[ListPipelineJobs] %w", jr.Error),
				}
				return
			}
			jobs = append(jobs, jr.Job)

			if opt.FetchSections || opt.FetchJobMetrics {
				job := jr.Job
				r, err := c.GetJobLog(ctx, projectID, job.ID)
				if err != nil {
					ch <- GetPipelineHierarchyResult{
						Error: fmt.Errorf("get job log: %w", err),
					}
					return
				}

				data, err := ParseJobLog(r)
				if err != nil {
					ch <- GetPipelineHierarchyResult{
						Error: fmt.Errorf("parse job log: %w", err),
					}
					return
				}

				if opt.FetchSections {
					for secnum, secdat := range data.Sections {
						section := &models.Section{
							Name:       secdat.Name,
							StartedAt:  unixTime(secdat.Start),
							FinishedAt: unixTime(secdat.End),
							Duration:   float64(secdat.End - secdat.Start),
						}

						section.ID = job.ID*1000 + int64(secnum)
						section.Job.ID = int64(job.ID)
						section.Job.Name = job.Name
						section.Job.Status = job.Status
						section.Pipeline.ID = int64(job.Pipeline.ID)
						section.Pipeline.ProjectID = int64(job.Pipeline.ProjectID)
						section.Pipeline.Ref = job.Pipeline.Ref
						section.Pipeline.Sha = job.Pipeline.Sha
						section.Pipeline.Status = job.Pipeline.Status

						sections = append(sections, section)
					}
				}

				if opt.FetchJobMetrics {
					for _, m := range data.Metrics {
						metric := &models.JobMetric{
							Name:      m.Name,
							Labels:    m.Labels,
							Value:     m.Value,
							Timestamp: m.Timestamp,
						}

						metric.Job.ID = job.ID
						metric.Job.Name = job.Name

						metrics = append(metrics, metric)
					}
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
			JobMetrics: metrics,
		}
	}()

	return ch
}
