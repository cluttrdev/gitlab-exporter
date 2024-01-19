package client

import (
	"context"
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	otlp_commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	otlp_tracepb "go.opentelemetry.io/proto/otlp/trace/v1"

	pb "github.com/cluttrdev/gitlab-exporter/grpc/exporterpb"
	"github.com/cluttrdev/gitlab-exporter/pkg/models"
)

type ClientConfig struct {
	Address string
	Options []grpc.DialOption
}

type Client struct {
	conn   grpc.ClientConnInterface
	client pb.GitLabExporterClient
}

func NewCLient(cfg ClientConfig) (*Client, error) {
	conn, err := grpc.Dial(cfg.Address, cfg.Options...)
	if err != nil {
		return nil, err
	}

	client := pb.NewGitLabExporterClient(conn)

	return &Client{
		conn:   conn,
		client: client,
	}, nil
}

func (c *Client) RecordPipelines(ctx context.Context, ps []*models.Pipeline) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordPipelines(ctx)
	if err != nil {
		return err
	}

	for _, p := range ps {
		err := stream.Send(&pb.Pipeline{
			Id:             p.ID,
			Iid:            p.IID,
			ProjectId:      p.ProjectID,
			Status:         p.Status,
			Source:         p.Source,
			Ref:            p.Ref,
			Sha:            p.SHA,
			BeforeSha:      p.BeforeSHA,
			Tag:            p.Tag,
			YamlErrors:     p.YamlErrors,
			CreatedAt:      convertTime(p.CreatedAt),
			UpdatedAt:      convertTime(p.UpdatedAt),
			StartedAt:      convertTime(p.StartedAt),
			FinishedAt:     convertTime(p.FinishedAt),
			CommittedAt:    convertTime(p.CommittedAt),
			Duration:       convertDuration(p.Duration),
			QueuedDuration: convertDuration(p.QueuedDuration),
			Coverage:       p.Coverage,
			WebUrl:         p.WebURL,
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordJobs(ctx context.Context, jobs []*models.Job) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordJobs(ctx)
	if err != nil {
		return err
	}

	for _, job := range jobs {
		err := stream.Send(&pb.Job{
			Pipeline: &pb.PipelineReference{
				Id:        job.Pipeline.ID,
				ProjectId: job.Pipeline.ProjectID,
				Ref:       job.Pipeline.Ref,
				Status:    job.Pipeline.Status,
			},
			Id:     job.ID,
			Name:   job.Name,
			Ref:    job.Ref,
			Stage:  job.Stage,
			Status: job.Status,

			CreatedAt:      convertTime(job.CreatedAt),
			StartedAt:      convertTime(job.StartedAt),
			FinishedAt:     convertTime(job.FinishedAt),
			ErasedAt:       convertTime(job.ErasedAt),
			Duration:       convertDuration(job.Duration),
			QueuedDuration: convertDuration(job.QueuedDuration),

			Coverage: job.Coverage,

			Tag:           job.Tag,
			AllowFailure:  job.AllowFailure,
			FailureReason: job.FailureReason,
			WebUrl:        job.WebURL,
			TagList:       job.TagList,
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordSections(ctx context.Context, sections []*models.Section) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordSections(ctx)
	if err != nil {
		return err
	}

	for _, section := range sections {
		err := stream.Send(&pb.Section{
			Job: &pb.JobReference{
				Id:     section.Job.ID,
				Name:   section.Job.Name,
				Status: section.Job.Status,
			},
			Pipeline: &pb.PipelineReference{
				Id:        section.Pipeline.ID,
				ProjectId: section.Pipeline.ProjectID,
				Ref:       section.Pipeline.Ref,
				Status:    section.Pipeline.Status,
			},
			Id:   section.ID,
			Name: section.Name,

			StartedAt:  convertTime(section.StartedAt),
			FinishedAt: convertTime(section.FinishedAt),
			Duration:   convertDuration(section.Duration),
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordBridges(ctx context.Context, bridges []*models.Bridge) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordBridges(ctx)
	if err != nil {
		return err
	}

	for _, bridge := range bridges {
		err := stream.Send(&pb.Bridge{
			Pipeline: &pb.PipelineInfo{
				Id:        bridge.Pipeline.ID,
				Iid:       bridge.Pipeline.IID,
				ProjectId: bridge.Pipeline.ProjectID,
				Status:    bridge.Pipeline.Status,
				Source:    bridge.Pipeline.Source,
				Ref:       bridge.Pipeline.Ref,
				Sha:       bridge.Pipeline.SHA,
				WebUrl:    bridge.Pipeline.WebURL,
				CreatedAt: convertTime(bridge.Pipeline.CreatedAt),
				UpdatedAt: convertTime(bridge.Pipeline.UpdatedAt),
			},
			Id:     bridge.ID,
			Name:   bridge.Name,
			Ref:    bridge.Ref,
			Stage:  bridge.Stage,
			Status: bridge.Status,

			DownstreamPipeline: &pb.PipelineInfo{
				Id:        bridge.DownstreamPipeline.ID,
				Iid:       bridge.DownstreamPipeline.IID,
				ProjectId: bridge.DownstreamPipeline.ProjectID,
				Status:    bridge.DownstreamPipeline.Status,
				Source:    bridge.DownstreamPipeline.Source,
				Ref:       bridge.DownstreamPipeline.Ref,
				Sha:       bridge.DownstreamPipeline.SHA,
				WebUrl:    bridge.DownstreamPipeline.WebURL,
				CreatedAt: convertTime(bridge.DownstreamPipeline.CreatedAt),
				UpdatedAt: convertTime(bridge.DownstreamPipeline.UpdatedAt),
			},

			CreatedAt:      convertTime(bridge.CreatedAt),
			StartedAt:      convertTime(bridge.StartedAt),
			FinishedAt:     convertTime(bridge.FinishedAt),
			ErasedAt:       convertTime(bridge.ErasedAt),
			Duration:       convertDuration(bridge.Duration),
			QueuedDuration: convertDuration(bridge.QueuedDuration),

			Coverage: bridge.Coverage,

			Tag:           bridge.Tag,
			AllowFailure:  bridge.AllowFailure,
			FailureReason: bridge.FailureReason,
			WebUrl:        bridge.WebURL,
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordTestReports(ctx context.Context, reports []*models.PipelineTestReport) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordTestReports(ctx)
	if err != nil {
		return err
	}

	for _, tr := range reports {
		err := stream.Send(&pb.TestReport{
			Id:           tr.ID,
			PipelineId:   tr.PipelineID,
			TotalTime:    tr.TotalTime,
			TotalCount:   tr.TotalCount,
			SuccessCount: tr.SkippedCount,
			FailedCount:  tr.FailedCount,
			SkippedCount: tr.SkippedCount,
			ErrorCount:   tr.ErrorCount,
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordTestSuites(ctx context.Context, suites []*models.PipelineTestSuite) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordTestSuites(ctx)
	if err != nil {
		return err
	}

	for _, ts := range suites {
		err := stream.Send(&pb.TestSuite{
			Id:           ts.ID,
			TestreportId: ts.TestReport.ID,
			PipelineId:   ts.TestReport.PipelineID,
			Name:         ts.Name,
			TotalTime:    ts.TotalTime,
			TotalCount:   ts.TotalCount,
			SuccessCount: ts.SkippedCount,
			FailedCount:  ts.FailedCount,
			SkippedCount: ts.SkippedCount,
			ErrorCount:   ts.ErrorCount,
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordTestCases(ctx context.Context, cases []*models.PipelineTestCase) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordTestCases(ctx)
	if err != nil {
		return err
	}

	for _, tc := range cases {
		err := stream.Send(&pb.TestCase{
			Id:            tc.ID,
			TestsuiteId:   tc.TestSuite.ID,
			TestreportId:  tc.TestReport.ID,
			PipelineId:    tc.TestReport.PipelineID,
			Status:        tc.Status,
			Name:          tc.Name,
			Classname:     tc.Classname,
			File:          tc.File,
			ExecutionTime: tc.ExecutionTime,
			SystemOutput:  tc.SystemOutput,
			StackTrace:    tc.StackTrace,
			AttachmentUrl: tc.AttachmentURL,
			RecentFailures: &pb.TestCaseRecentFailures{
				Count:      tc.RecentFailures.Count,
				BaseBranch: tc.RecentFailures.BaseBranch,
			},
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordLogEmbeddedMetrics(ctx context.Context, metrics []*models.JobMetric) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordLogEmbeddedMetrics(ctx)
	if err != nil {
		return err
	}

	for _, m := range metrics {
		labels := make([]*pb.LogEmbeddedMetric_Label, 0, len(m.Labels))
		for name, value := range m.Labels {
			labels = append(labels, &pb.LogEmbeddedMetric_Label{
				Name:  name,
				Value: value,
			})
		}
		err := stream.Send(&pb.LogEmbeddedMetric{
			Name:      m.Name,
			Labels:    labels,
			Value:     m.Value,
			Timestamp: convertUnixMilli(m.Timestamp),
			Job: &pb.LogEmbeddedMetric_JobReference{
				Id:   m.Job.ID,
				Name: m.Job.Name,
			},
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RecordTraces(ctx context.Context, traces []*models.Trace) error {
	if c == nil {
		return errors.New("nil client")
	}
	stream, err := c.client.RecordTraces(ctx)
	if err != nil {
		return err
	}

	for _, t := range traces {
		spans := make([]*otlp_tracepb.Span, 0, len(*t))
		for _, s := range *t {
			spans = append(spans, &otlp_tracepb.Span{
				TraceId:                []byte(s.TraceID),
				SpanId:                 []byte(s.SpanID),
				TraceState:             s.TraceState,
				ParentSpanId:           []byte(s.ParentSpanID),
				Name:                   s.Name,
				Kind:                   otlp_tracepb.Span_SpanKind(s.Kind),
				StartTimeUnixNano:      s.StartTime,
				EndTimeUnixNano:        s.EndTime,
				Attributes:             convertSpanAttributes(s.Attributes),
				DroppedAttributesCount: 0,
				Events:                 convertSpanEvents(s.Events),
				DroppedEventsCount:     0,
				Links:                  convertSpanLinks(s.Links),
				DroppedLinksCount:      0,
				Status:                 convertSpanStatus(s.Status),
			})
		}
		err := stream.Send(&pb.Trace{
			Spans: spans,
		})
		if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

func convertUnixMilli(ts int64) *timestamppb.Timestamp {
	const msPerSecond int64 = 1_000
	const nsPerMilli int64 = 1_000
	return &timestamppb.Timestamp{
		Seconds: ts / msPerSecond,
		Nanos:   int32((ts % msPerSecond) * nsPerMilli),
	}
}

func convertTime(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func convertDuration(d float64) *durationpb.Duration {
	return durationpb.New(time.Duration(d * float64(time.Second)))
}

func convertSpanAttributes(attrs map[string]string) []*otlp_commonpb.KeyValue {
	list := make([]*otlp_commonpb.KeyValue, 0, len(attrs))

	for key, value := range attrs {
		list = append(list, &otlp_commonpb.KeyValue{
			Key: key,
			Value: &otlp_commonpb.AnyValue{
				Value: &otlp_commonpb.AnyValue_StringValue{
					StringValue: value,
				},
			},
		})
	}

	return list
}

func convertSpanEvents(events []models.SpanEvent) []*otlp_tracepb.Span_Event {
	list := make([]*otlp_tracepb.Span_Event, 0, len(events))

	for _, e := range events {
		list = append(list, &otlp_tracepb.Span_Event{
			TimeUnixNano:           e.Time,
			Name:                   e.Name,
			Attributes:             convertSpanAttributes(e.Attributes),
			DroppedAttributesCount: 0,
		})
	}

	return list
}

func convertSpanLinks(links []models.SpanLink) []*otlp_tracepb.Span_Link {
	list := make([]*otlp_tracepb.Span_Link, 0, len(links))

	for _, l := range links {
		list = append(list, &otlp_tracepb.Span_Link{
			TraceId:                []byte(l.TraceID),
			SpanId:                 []byte(l.SpanID),
			TraceState:             l.TraceState,
			Attributes:             convertSpanAttributes(l.Attributes),
			DroppedAttributesCount: 0,
		})
	}

	return list
}

func convertSpanStatus(s models.SpanStatus) *otlp_tracepb.Status {
	var status otlp_tracepb.Status

	status.Message = s.Message

	switch s.Code {
	case models.StatusCodeUnset:
		status.Code = otlp_tracepb.Status_STATUS_CODE_UNSET
	case models.StatusCodeOk:
		status.Code = otlp_tracepb.Status_STATUS_CODE_OK
	case models.StatusCodeError:
		status.Code = otlp_tracepb.Status_STATUS_CODE_ERROR
	}

	return &status
}
