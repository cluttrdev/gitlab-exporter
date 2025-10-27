package messages

import (
	"strconv"
	"strings"
	"time"

	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"

	"go.cluttr.dev/gitlab-exporter/internal/types"
)

func NewPipelineSpan(pipeline types.Pipeline) *tracepb.Span {
	var (
		traceId  = []byte(strconv.FormatInt(pipeline.Id, 10))
		spanId   = []byte(strconv.FormatInt(pipeline.Id, 10))
		parentId []byte

		links []spanLinkData
	)

	if pipeline.UpstreamPipeline != nil {
		links = append(links, spanLinkData{
			TraceId: []byte(strconv.FormatInt(pipeline.UpstreamPipeline.Id, 10)),
			SpanId:  []byte(strconv.FormatInt(pipeline.UpstreamPipeline.Id, 10)),
		})
	}

	name := pipeline.Name
	if name == "" {
		name = pipeline.Ref
	}

	return newSpan(
		spanData{
			TraceId:   traceId,
			SpanId:    spanId,
			ParentId:  parentId,
			Name:      name,
			StartTime: pipeline.StartedAt,
			EndTime:   pipeline.FinishedAt,
			Attributes: map[string]string{
				"ci.pipeline.status": pipeline.Status,
			},
			Links:  links,
			Status: pipeline.Status,
		},
	)
}

func NewJobSpan(job types.Job) *tracepb.Span {
	traceId := []byte(strconv.FormatInt(job.Pipeline.Id, 10))
	spanId := []byte(strconv.FormatInt(job.Id, 10))
	parentId := []byte(strconv.FormatInt(job.Pipeline.Id, 10))

	var links []spanLinkData
	if job.DownstreamPipeline != nil {
		links = append(links, spanLinkData{
			TraceId: []byte(strconv.FormatInt(job.DownstreamPipeline.Id, 10)),
			SpanId:  []byte(strconv.FormatInt(job.DownstreamPipeline.Id, 10)),
		})
	}

	return newSpan(
		spanData{
			TraceId:   traceId,
			SpanId:    spanId,
			ParentId:  parentId,
			Name:      job.Name,
			StartTime: job.StartedAt,
			EndTime:   job.FinishedAt,
			Attributes: map[string]string{
				"ci.job.status": job.Status,
			},
			Links:  links,
			Status: job.Status,
		},
	)
}

func NewSectionSpan(section types.Section) *tracepb.Span {
	traceId := []byte(strconv.FormatInt(section.Job.Pipeline.Id, 10))
	spanId := []byte(strconv.FormatInt(section.Id, 10))
	parentId := []byte(strconv.FormatInt(section.Job.Id, 10))

	return newSpan(
		spanData{
			TraceId:   traceId,
			SpanId:    spanId,
			ParentId:  parentId,
			Name:      section.Name,
			StartTime: section.StartedAt,
			EndTime:   section.FinishedAt,
			Status:    "",
		},
	)
}

type spanData struct {
	TraceId  []byte
	SpanId   []byte
	ParentId []byte

	Name      string
	StartTime *time.Time
	EndTime   *time.Time

	Attributes map[string]string
	Links      []spanLinkData
	Status     string
}

type spanLinkData struct {
	TraceId []byte
	SpanId  []byte
}

func newSpan(span spanData) *tracepb.Span {
	var startTime, endTime uint64
	if span.StartTime != nil {
		startTime = uint64(span.StartTime.UnixNano())
	}
	if span.EndTime != nil {
		endTime = uint64(span.EndTime.UnixNano())
	}

	s := &tracepb.Span{
		TraceId:                span.TraceId,
		SpanId:                 span.SpanId,
		TraceState:             "",
		ParentSpanId:           span.ParentId,
		Name:                   span.Name,
		Kind:                   tracepb.Span_SPAN_KIND_INTERNAL,
		StartTimeUnixNano:      startTime,
		EndTimeUnixNano:        endTime,
		Attributes:             convertAttributes(span.Attributes),
		DroppedAttributesCount: 0,
		Events:                 []*tracepb.Span_Event{},
		DroppedEventsCount:     0,
		Links:                  []*tracepb.Span_Link{},
		DroppedLinksCount:      0,
		Status:                 convertStatus(span.Status),
	}

	for _, l := range span.Links {
		s.Links = append(s.Links, &tracepb.Span_Link{
			TraceId: l.TraceId,
			SpanId:  l.SpanId,
		})
	}

	return s
}

func NewResourceSpan(attrs map[string]string, spans []*tracepb.Span) *tracepb.ResourceSpans {
	scopeSpans := &tracepb.ScopeSpans{
		Scope: &commonpb.InstrumentationScope{},
		Spans: spans,
	}

	return &tracepb.ResourceSpans{
		Resource: &resourcepb.Resource{
			Attributes: convertAttributes(attrs),
		},
		ScopeSpans: []*tracepb.ScopeSpans{
			scopeSpans,
		},
	}
}

func convertAttributes(attrs map[string]string) []*commonpb.KeyValue {
	list := make([]*commonpb.KeyValue, 0, len(attrs))

	for key, value := range attrs {
		list = append(list, &commonpb.KeyValue{
			Key: key,
			Value: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_StringValue{
					StringValue: value,
				},
			},
		})
	}

	return list
}

func convertStatus(status string) *tracepb.Status {
	var code tracepb.Status_StatusCode
	switch strings.ToLower(status) {
	case "success":
		code = tracepb.Status_STATUS_CODE_OK
	case "failed", "canceled":
		code = tracepb.Status_STATUS_CODE_ERROR
	default:
		code = tracepb.Status_STATUS_CODE_UNSET
	}

	return &tracepb.Status{
		Message: status,
		Code:    code,
	}
}
