package clickhouseclient

import (
	"context"
	"fmt"
	"time"

	"github.com/cluttrdev/gitlab-clickhouse-exporter/pkg/models"
)

func timestamp(t *time.Time) float64 {
	const msPerS float64 = 1000.0
	if t == nil {
		return 0.0
	}
	return float64(t.UnixMilli()) / msPerS
}

func InsertPipelines(ctx context.Context, pipelines []*models.Pipeline, c *Client) error {
	batch, err := c.PrepareBatch(ctx, "INSERT INTO gitlab_ci.pipelines")
	if err != nil {
		return fmt.Errorf("[clickhouse.Client.InsertPipelines] %w", err)
	}

	for _, p := range pipelines {
		err = batch.Append(
			p.ID,
			p.IID,
			p.ProjectID,
			p.Status,
			p.Source,
			p.Ref,
			p.SHA,
			p.BeforeSHA,
			p.Tag,
			p.YamlErrors,
			timestamp(p.CreatedAt),
			timestamp(p.UpdatedAt),
			timestamp(p.StartedAt),
			timestamp(p.FinishedAt),
			timestamp(p.CommittedAt),
			p.Duration,
			p.QueuedDuration,
			p.Coverage,
			p.WebURL,
		)
		if err != nil {
			return fmt.Errorf("[clickhouse.Client.InsertPipelines] %w", err)
		}
	}

	return batch.Send()
}

func InsertJobs(ctx context.Context, jobs []*models.Job, c *Client) error {
	batch, err := c.PrepareBatch(ctx, "INSERT INTO gitlab_ci.jobs")
	if err != nil {
		return fmt.Errorf("[clickhouse.Client.InsertJobs] %w", err)
	}

	for _, j := range jobs {
		err = batch.Append(
			j.Coverage,
			j.AllowFailure,
			timestamp(j.CreatedAt),
			timestamp(j.StartedAt),
			timestamp(j.FinishedAt),
			timestamp(j.ErasedAt),
			j.Duration,
			j.QueuedDuration,
			j.TagList,
			j.ID,
			j.Name,
			map[string]interface{}{
				"id":         j.Pipeline.ID,
				"project_id": j.Pipeline.ProjectID,
				"ref":        j.Pipeline.Ref,
				"sha":        j.Pipeline.Sha,
				"status":     j.Pipeline.Status,
			},
			j.Ref,
			j.Stage,
			j.Status,
			j.FailureReason,
			j.Tag,
			j.WebURL,
		)
		if err != nil {
			return fmt.Errorf("[clickhouse.Client.InsertJobs] %w", err)
		}
	}

	return batch.Send()
}

func InsertBridges(ctx context.Context, bridges []*models.Bridge, c *Client) error {
	batch, err := c.PrepareBatch(ctx, "INSERT INTO gitlab_ci.bridges")
	if err != nil {
		return fmt.Errorf("[clickhouse.Client.InsertBridges] %w", err)
	}

	for _, b := range bridges {
		err = batch.Append(
			b.Coverage,
			b.AllowFailure,
			timestamp(b.CreatedAt),
			timestamp(b.StartedAt),
			timestamp(b.FinishedAt),
			timestamp(b.ErasedAt),
			b.Duration,
			b.QueuedDuration,
			b.ID,
			b.Name,
			map[string]interface{}{
				"id":         b.Pipeline.ID,
				"iid":        b.Pipeline.IID,
				"project_id": b.Pipeline.ProjectID,
				"status":     b.Pipeline.Status,
				"source":     b.Pipeline.Source,
				"ref":        b.Pipeline.Source,
				"sha":        b.Pipeline.SHA,
				"web_url":    b.Pipeline.WebURL,
				"created_at": timestamp(b.Pipeline.CreatedAt),
				"updated_at": timestamp(b.Pipeline.UpdatedAt),
			},
			b.Ref,
			b.Stage,
			b.Status,
			b.FailureReason,
			b.Tag,
			b.WebURL,
			map[string]interface{}{
				"id":         b.DownstreamPipeline.ID,
				"iid":        b.DownstreamPipeline.IID,
				"project_id": b.DownstreamPipeline.ProjectID,
				"status":     b.DownstreamPipeline.Status,
				"source":     b.DownstreamPipeline.Source,
				"ref":        b.DownstreamPipeline.Source,
				"sha":        b.DownstreamPipeline.SHA,
				"web_url":    b.DownstreamPipeline.WebURL,
				"created_at": timestamp(b.DownstreamPipeline.CreatedAt),
				"updated_at": timestamp(b.DownstreamPipeline.UpdatedAt),
			},
		)
		if err != nil {
			return fmt.Errorf("[clickhouse.Client.InsertBridges] %w", err)
		}
	}

	return batch.Send()
}

func InsertSections(ctx context.Context, sections []*models.Section, client *Client) error {
	batch, err := client.PrepareBatch(ctx, "INSERT INTO gitlab_ci.sections")
	if err != nil {
		return fmt.Errorf("[clickhouse.Client.InsertSections] %w", err)
	}

	for _, s := range sections {
		err = batch.Append(
			s.ID,
			s.Name,
			map[string]interface{}{
				"id":     s.Job.ID,
				"name":   s.Job.Name,
				"status": s.Job.Status,
			},
			map[string]interface{}{
				"id":         s.Pipeline.ID,
				"project_id": s.Pipeline.ProjectID,
				"ref":        s.Pipeline.Ref,
				"sha":        s.Pipeline.Sha,
				"status":     s.Pipeline.Status,
			},
			timestamp(s.StartedAt),
			timestamp(s.FinishedAt),
			s.Duration,
		)
		if err != nil {
			return fmt.Errorf("[clickhouse.Client.InsertSections] %w", err)
		}
	}

	return batch.Send()
}

func InsertPipelineHierarchy(ctx context.Context, hierarchy *models.PipelineHierarchy, client *Client) error {
	if err := InsertPipelines(ctx, hierarchy.GetAllPipelines(), client); err != nil {
		return fmt.Errorf("[InsertPipelineHierarchy] %w", err)
	}

	if err := InsertJobs(ctx, hierarchy.GetAllJobs(), client); err != nil {
		return fmt.Errorf("[InsertPipelineHierarchy] %w", err)
	}

	if err := InsertSections(ctx, hierarchy.GetAllSections(), client); err != nil {
		return fmt.Errorf("[InsertPipelineHierarchy] %w", err)
	}

	if err := InsertBridges(ctx, hierarchy.GetAllBridges(), client); err != nil {
		return fmt.Errorf("[InsertPipelineHierarchy] %w", err)
	}

	return nil
}

func timeFromUnixNano(ts int64) time.Time {
	const nsecPerSecond int64 = 1e09
	sec := ts / nsecPerSecond
	nsec := ts - (sec * nsecPerSecond)
	return time.Unix(sec, nsec)
}

func convertEvents(events []models.SpanEvent) ([]time.Time, []string, []map[string]string) {
	var (
		times []time.Time
		names []string
		attrs []map[string]string
	)
	for _, event := range events {
		times = append(times, timeFromUnixNano(int64(event.Time)))
		names = append(names, event.Name)
		attrs = append(attrs, event.Attributes)
	}
	return times, names, attrs
}

func convertLinks(links []models.SpanLink) ([]string, []string, []string, []map[string]string) {
	var (
		traceIDs []string
		spanIDs  []string
		states   []string
		attrs    []map[string]string
	)
	for _, link := range links {
		traceIDs = append(traceIDs, link.TraceID)
		spanIDs = append(spanIDs, link.SpanID)
		states = append(states, link.TraceState)
		attrs = append(attrs, link.Attributes)
	}
	return traceIDs, spanIDs, states, attrs
}

func InsertTraces(ctx context.Context, traces [][]*models.Span, client *Client) error {
	batch, err := client.PrepareBatch(ctx, "INSERT INTO gitlab_ci.traces")
	if err != nil {
		return fmt.Errorf("[clickhouse.Client.InsertTraces] %w", err)
	}

	scopeName := ""
	scopeVersion := ""

	for _, trace := range traces {
		for _, span := range trace {
			serviceName := ""
			if sn, ok := span.Resource.Attributes["service.name"]; ok {
				serviceName = sn
			}

			eventTimes, eventNames, eventAttrs := convertEvents(span.Events)
			linkTraceIDs, linkSpanIDs, linkStates, linkAttrs := convertLinks(span.Links)

			err = batch.Append(
				timeFromUnixNano(int64(span.StartTime)),
				span.TraceID,
				span.SpanID,
				span.ParentSpanID,
				span.TraceState,
				span.Name,
				span.Kind.Name(),
				serviceName,
				span.Resource.Attributes,
				scopeName,
				scopeVersion,
				span.Attributes,
				int64(span.EndTime-span.StartTime),
				span.Status.Code.Name(),
				span.Status.Message,
				eventTimes,
				eventNames,
				eventAttrs,
				linkTraceIDs,
				linkSpanIDs,
				linkStates,
				linkAttrs,
			)

			if err != nil {
				return fmt.Errorf("[clickhouse.Client.InsertTraces] %w", err)
			}
		}
	}

	return batch.Send()
}

func InsertTrace(ctx context.Context, trace []*models.Span, client *Client) error {
	return InsertTraces(ctx, [][]*models.Span{trace}, client)
}
