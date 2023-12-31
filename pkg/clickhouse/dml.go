package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/cluttrdev/gitlab-exporter/pkg/models"
)

func timestamp(t *time.Time) float64 {
	const msPerS float64 = 1000.0
	if t == nil {
		return 0.0
	}
	return float64(t.UnixMilli()) / msPerS
}

func InsertPipelines(ctx context.Context, pipelines []*models.Pipeline, c *Client) error {
	const query string = `INSERT INTO {db: Identifier}.pipelines`
	var params = map[string]string{
		"db": c.dbName,
	}

	ctx = WithParameters(ctx, params)

	batch, err := c.PrepareBatch(ctx, query)
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
	const query string = `INSERT INTO {db: Identifier}.jobs`
	var params = map[string]string{
		"db": c.dbName,
	}

	ctx = WithParameters(ctx, params)

	batch, err := c.PrepareBatch(ctx, query)
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
	const query string = `INSERT INTO {db: Identifier}.bridges`
	var params = map[string]string{
		"db": c.dbName,
	}

	ctx = WithParameters(ctx, params)

	batch, err := c.PrepareBatch(ctx, query)
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
	const query string = `INSERT INTO {db: Identifier}.sections`
	var params = map[string]string{
		"db": client.dbName,
	}

	ctx = WithParameters(ctx, params)

	batch, err := client.PrepareBatch(ctx, query)
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

func InsertTestReports(ctx context.Context, reports []*models.PipelineTestReport, client *Client) error {
	const query string = `INSERT INTO {db: Identifier}.testreports`
	var params = map[string]string{
		"db": client.dbName,
	}

	ctx = WithParameters(ctx, params)

	batch, err := client.PrepareBatch(ctx, query)
	if err != nil {
		return fmt.Errorf("[clickhouse.Client.InsertTestReports] %w", err)
	}

	for _, tr := range reports {
		ids, names, times, counts := convertTestSuitesSummary(tr.TestSuites)

		err = batch.Append(
			tr.ID,
			tr.PipelineID,
			tr.TotalTime,
			tr.TotalCount,
			tr.SuccessCount,
			tr.FailedCount,
			tr.SkippedCount,
			tr.ErrorCount,
			ids,
			names,
			times,
			counts,
		)
		if err != nil {
			return fmt.Errorf("[clickhouse.Client.InsertTestReports] %w", err)
		}
	}

	return batch.Send()
}

func convertTestSuitesSummary(suites []*models.PipelineTestSuite) (ids []int64, names []string, times []float64, counts []int64) {
	for _, ts := range suites {
		ids = append(ids, ts.ID)
		names = append(names, ts.Name)
		times = append(times, ts.TotalTime)
		counts = append(counts, ts.TotalCount)
	}

	return
}

func InsertTestSuites(ctx context.Context, suites []*models.PipelineTestSuite, client *Client) error {
	const query string = `INSERT INTO {db: Identifier}.testsuites`
	var params = map[string]string{
		"db": client.dbName,
	}

	ctx = WithParameters(ctx, params)

	batch, err := client.PrepareBatch(ctx, query)
	if err != nil {
		return fmt.Errorf("[clickhouse.Client.InsertTestReports] %w", err)
	}

	for _, ts := range suites {
		ids, statuses, names := convertTestCasesSummary(ts.TestCases)

		err = batch.Append(
			ts.ID,
			map[string]interface{}{
				"id":          ts.TestReport.ID,
				"pipeline_id": ts.TestReport.PipelineID,
			},
			ts.Name,
			ts.TotalTime,
			ts.TotalCount,
			ts.SuccessCount,
			ts.FailedCount,
			ts.SkippedCount,
			ts.ErrorCount,
			ids,
			statuses,
			names,
		)
		if err != nil {
			return fmt.Errorf("[clickhouse.Client.InsertTestSuites] %w", err)
		}
	}

	return batch.Send()
}

func convertTestCasesSummary(cases []*models.PipelineTestCase) (ids []int64, statuses []string, names []string) {
	for _, tc := range cases {
		ids = append(ids, tc.ID)
		statuses = append(statuses, tc.Status)
		names = append(names, tc.Name)
	}

	return
}

func InsertTestCases(ctx context.Context, cases []*models.PipelineTestCase, client *Client) error {
	const query string = `INSERT INTO {db: Identifier}.testcases`
	var params = map[string]string{
		"db": client.dbName,
	}

	ctx = WithParameters(ctx, params)

	batch, err := client.PrepareBatch(ctx, query)
	if err != nil {
		return fmt.Errorf("[clickhouse.Client.InsertTestReports] %w", err)
	}

	for _, tc := range cases {
		err = batch.Append(
			tc.ID,
			map[string]interface{}{
				"id": tc.TestSuite.ID,
			},
			map[string]interface{}{
				"id":          tc.TestReport.ID,
				"pipeline_id": tc.TestReport.PipelineID,
			},
			tc.Status,
			tc.Name,
			tc.Classname,
			tc.File,
			tc.ExecutionTime,
			tc.SystemOutput,
			tc.StackTrace,
			tc.AttachmentURL,
			map[string]interface{}{
				"count":       tc.RecentFailures.Count,
				"base_branch": tc.RecentFailures.BaseBranch,
			},
		)
		if err != nil {
			return fmt.Errorf("[clickhouse.Client.InsertTestCases] %w", err)
		}
	}

	return batch.Send()
}

func InsertJobMetrics(ctx context.Context, metrics []*models.JobMetric, client *Client) error {
	const query string = `INSERT INTO {db: Identifier}.{table: Identifier}`
	var params = map[string]string{
		"db":    client.dbName,
		"table": "job_metrics",
	}

	ctx = WithParameters(ctx, params)

	batch, err := client.PrepareBatch(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare batch: %w", err)
	}

	for _, m := range metrics {
		err = batch.Append(
			m.Name,
			m.Labels,
			m.Value,
			m.Timestamp,
			map[string]interface{}{
				"id":   m.Job.ID,
				"name": m.Job.Name,
			},
		)
		if err != nil {
			return fmt.Errorf("append batch:  %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("send batch: %w", err)
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

func InsertTraces(ctx context.Context, traces []models.Trace, client *Client) error {
	const query string = `INSERT INTO {db: Identifier}.traces`
	var params = map[string]string{
		"db": client.dbName,
	}

	ctx = WithParameters(ctx, params)

	batch, err := client.PrepareBatch(ctx, query)
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
	return InsertTraces(ctx, []models.Trace{trace}, client)
}
