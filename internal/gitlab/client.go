package gitlab

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/time/rate"

	gitlab "github.com/xanzy/go-gitlab"

	"github.com/cluttrdev/gitlab-exporter/internal/types"
	"github.com/cluttrdev/gitlab-exporter/protobuf/typespb"
)

type Client struct {
	sync.RWMutex
	client *gitlab.Client
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

func (c *Client) Client() *gitlab.Client {
	return c.client
}

func (c *Client) Configure(cfg ClientConfig) error {
	opts := []gitlab.ClientOptionFunc{
		gitlab.WithBaseURL(cfg.URL),
	}

	if cfg.RateLimit > 0 {
		limit := rate.Limit(cfg.RateLimit * 0.66)
		burst := cfg.RateLimit * 0.33
		limiter := rate.NewLimiter(limit, int(burst))

		opts = append(opts, gitlab.WithCustomLimiter(limiter))
	}

	client, err := gitlab.NewOAuthClient(cfg.Token, opts...)
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
		[]gitlab.RequestOptionFunc{gitlab.WithContext(ctx)},
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
	PipelineHierarchy *PipelineHierarchy
	Metrics           []*typespb.Metric
}

func (c *Client) GetPipelineHierarchy(ctx context.Context, projectID int64, pipelineID int64, opt *GetPipelineHierarchyOptions) (*GetPipelineHierarchyResult, error) {
	pipeline, err := c.GetPipeline(ctx, projectID, pipelineID)
	if err != nil {
		return nil, err
	}

	jobs := []*typespb.Job{}
	sections := []*typespb.Section{}
	metrics := []*typespb.Metric{}

	js, err := c.GetPipelineJobs(ctx, projectID, pipelineID)
	if err != nil {
		return nil, fmt.Errorf("get pipeline jobs: %w", err)
	}

	jobs = append(jobs, js...)
	if opt.FetchSections || opt.FetchJobMetrics {
		for _, job := range js {
			r, err := c.GetJobLog(ctx, projectID, job.Id)
			if err != nil {
				return nil, fmt.Errorf("get job log: %w", err)
			}

			data, err := ParseJobLog(r)
			if err != nil {
				return nil, fmt.Errorf("parse job log: %w", err)
			}

			if opt.FetchSections {
				for secnum, secdat := range data.Sections {
					section := &typespb.Section{
						Id: job.Id*1000 + int64(secnum),
						Job: &typespb.JobReference{
							Id:     job.Id,
							Name:   job.Name,
							Status: job.Status,
						},
						Pipeline:   job.Pipeline,
						Name:       secdat.Name,
						StartedAt:  types.ConvertUnixSeconds(secdat.Start),
						FinishedAt: types.ConvertUnixSeconds(secdat.End),
						Duration:   types.ConvertDuration(float64(secdat.End - secdat.Start)),
					}

					sections = append(sections, section)
				}
			}

			if opt.FetchJobMetrics {
				var metricIID int = 0
				for _, m := range data.Metrics {
					metricIID++
					metric := &typespb.Metric{
						Id:        []byte(fmt.Sprintf("%d-%d", job.Id, metricIID)),
						Iid:       int64(metricIID),
						JobId:     job.Id,
						Name:      m.Name,
						Labels:    convertLabels(m.Labels),
						Value:     m.Value,
						Timestamp: types.ConvertUnixMilli(m.Timestamp),
					}
					metrics = append(metrics, metric)
				}
			}
		}
	}

	bridges := []*typespb.Bridge{}
	dps := []*PipelineHierarchy{}

	bs, err := c.GetPipelineBridges(ctx, projectID, pipelineID)
	if err != nil {
		return nil, fmt.Errorf("get pipeline bridges: %w", err)
	}

	bridges = append(bridges, bs...)
	for _, b := range bs {
		dp := b.DownstreamPipeline
		if dp == nil || dp.Id == 0 {
			continue
		}

		dpr, err := c.GetPipelineHierarchy(ctx, dp.ProjectId, dp.Id, opt)
		if err != nil {
			return nil, fmt.Errorf("get downstream pipeline hierarchy: %w", err)
		}
		dps = append(dps, dpr.PipelineHierarchy)
		metrics = append(metrics, dpr.Metrics...)
	}

	return &GetPipelineHierarchyResult{
		PipelineHierarchy: &PipelineHierarchy{
			Pipeline:            pipeline,
			Jobs:                jobs,
			Sections:            sections,
			Bridges:             bridges,
			DownstreamPipelines: dps,
		},
		Metrics: metrics,
	}, nil
}

func convertLabels(labels map[string]string) []*typespb.Metric_Label {
	list := make([]*typespb.Metric_Label, 0, len(labels))
	for name, value := range labels {
		list = append(list, &typespb.Metric_Label{
			Name:  name,
			Value: value,
		})
	}
	return list
}
