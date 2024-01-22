package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/time/rate"

	_gitlab "github.com/xanzy/go-gitlab"

	pb "github.com/cluttrdev/gitlab-exporter/grpc/exporterpb"
	"github.com/cluttrdev/gitlab-exporter/internal/models"
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
	PipelineHierarchy  *models.PipelineHierarchy
	LogEmbeddedMetrics []*pb.LogEmbeddedMetric
	Error              error
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

		jobs := []*pb.Job{}
		sections := []*pb.Section{}
		metrics := []*pb.LogEmbeddedMetric{}
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
				r, err := c.GetJobLog(ctx, projectID, job.Id)
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
						section := &pb.Section{
							Id: job.Id*1000 + int64(secnum),
							Job: &pb.JobReference{
								Id:     job.Id,
								Name:   job.Name,
								Status: job.Status,
							},
							Pipeline:   job.Pipeline,
							Name:       secdat.Name,
							StartedAt:  models.ConvertUnixSeconds(secdat.Start),
							FinishedAt: models.ConvertUnixSeconds(secdat.End),
							Duration:   models.ConvertDuration(float64(secdat.End - secdat.Start)),
						}

						sections = append(sections, section)
					}
				}

				if opt.FetchJobMetrics {
					for _, m := range data.Metrics {
						metric := &pb.LogEmbeddedMetric{
							Name:      m.Name,
							Labels:    models.ConvertLabels(m.Labels),
							Value:     m.Value,
							Timestamp: models.ConvertUnixMilli(m.Timestamp),
							Job: &pb.LogEmbeddedMetric_JobReference{
								Id:   job.Id,
								Name: job.Name,
							},
						}

						metrics = append(metrics, metric)
					}
				}
			}
		}

		bridges := []*pb.Bridge{}
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
			if dp == nil || dp.Id == 0 {
				continue
			}

			dpr := <-c.GetPipelineHierarchy(ctx, dp.ProjectId, dp.Id, opt)
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
			LogEmbeddedMetrics: metrics,
		}
	}()

	return ch
}
