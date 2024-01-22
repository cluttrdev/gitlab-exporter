package models

import (
	"fmt"

	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/cluttrdev/gitlab-exporter/grpc/exporterpb"
)

func NewPipelineHierarchyTrace(traceID []byte, parentID []byte, ph *PipelineHierarchy) *pb.Trace {
	var data tracepb.TracesData

	if len(traceID) == 0 {
		traceID = []byte(fmt.Sprint(ph.Pipeline.Id))
	}

	data.ResourceSpans = append(data.ResourceSpans, NewPipelineSpan(traceID, parentID, ph.Pipeline))

	for _, job := range ph.Jobs {
		data.ResourceSpans = append(data.ResourceSpans, NewJobSpan(traceID, job))
	}

	for _, section := range ph.Sections {
		data.ResourceSpans = append(data.ResourceSpans, NewSectionSpan(traceID, section))
	}

	for _, bridge := range ph.Bridges {
		data.ResourceSpans = append(data.ResourceSpans, NewBridgeSpan(traceID, bridge))

		index := slices.IndexFunc(ph.DownstreamPipelines, func(dph *PipelineHierarchy) bool {
			return dph.Pipeline.Id == bridge.DownstreamPipeline.Id
		})
		if index > -1 {
			dtrace := NewPipelineHierarchyTrace(traceID, []byte(fmt.Sprint(bridge.Id)), ph.DownstreamPipelines[index])
			data.ResourceSpans = append(data.ResourceSpans, dtrace.Data.ResourceSpans...)
		}
	}

	return &pb.Trace{
		Data: &data,
	}
}

func NewPipelineSpan(traceID []byte, parentID []byte, pipeline *pb.Pipeline) *tracepb.ResourceSpans {
	return &tracepb.ResourceSpans{
		Resource: &resourcepb.Resource{
			Attributes: convertAttributes(map[string]string{
				"service.name": "gitlab_ci.pipeline",
			}),
			DroppedAttributesCount: 0,
		},
		ScopeSpans: []*tracepb.ScopeSpans{
			&tracepb.ScopeSpans{
				Scope: &commonpb.InstrumentationScope{
					Name:                   "",
					Version:                "",
					Attributes:             []*commonpb.KeyValue{},
					DroppedAttributesCount: 0,
				},
				Spans: []*tracepb.Span{
					&tracepb.Span{
						TraceId:           traceID,
						SpanId:            []byte(fmt.Sprint(pipeline.Id)),
						TraceState:        "",
						ParentSpanId:      parentID,
						Name:              pipeline.Ref,
						Kind:              tracepb.Span_SPAN_KIND_INTERNAL,
						StartTimeUnixNano: timestampUnixNano(pipeline.StartedAt),
						EndTimeUnixNano:   timestampUnixNano(pipeline.FinishedAt),
						Attributes: convertAttributes(map[string]string{
							"ci.pipeline.status":  pipeline.Status,
							"ci.pipeline.web_url": pipeline.WebUrl,
						}),
						DroppedAttributesCount: 0,
						Events:                 []*tracepb.Span_Event{},
						DroppedEventsCount:     0,
						Links:                  []*tracepb.Span_Link{},
						DroppedLinksCount:      0,
						Status:                 convertStatus(pipeline.Status),
					},
				},
				SchemaUrl: "",
			},
		},
		SchemaUrl: "",
	}
}

func NewJobSpan(traceID []byte, job *pb.Job) *tracepb.ResourceSpans {
	return &tracepb.ResourceSpans{
		Resource: &resourcepb.Resource{
			Attributes: convertAttributes(map[string]string{
				"service.name": "gitlab_ci.job",
			}),
			DroppedAttributesCount: 0,
		},
		ScopeSpans: []*tracepb.ScopeSpans{
			&tracepb.ScopeSpans{
				Scope: &commonpb.InstrumentationScope{
					Name:                   "",
					Version:                "",
					Attributes:             []*commonpb.KeyValue{},
					DroppedAttributesCount: 0,
				},
				Spans: []*tracepb.Span{
					&tracepb.Span{
						TraceId:           traceID,
						SpanId:            []byte(fmt.Sprint(job.Id)),
						TraceState:        "",
						ParentSpanId:      []byte(fmt.Sprint(job.Pipeline.Id)),
						Name:              job.Name,
						Kind:              tracepb.Span_SPAN_KIND_INTERNAL,
						StartTimeUnixNano: timestampUnixNano(job.StartedAt),
						EndTimeUnixNano:   timestampUnixNano(job.FinishedAt),
						Attributes: convertAttributes(map[string]string{
							"ci.job.status":  job.Status,
							"ci.job.web_url": job.WebUrl,
						}),
						DroppedAttributesCount: 0,
						Events:                 []*tracepb.Span_Event{},
						DroppedEventsCount:     0,
						Links:                  []*tracepb.Span_Link{},
						DroppedLinksCount:      0,
						Status:                 convertStatus(job.Status),
					},
				},
				SchemaUrl: "",
			},
		},
		SchemaUrl: "",
	}
}

func NewBridgeSpan(traceID []byte, bridge *pb.Bridge) *tracepb.ResourceSpans {
	attrs := map[string]string{
		"ci.job.status":  bridge.Status,
		"ci.job.web_url": bridge.WebUrl,
	}
	if bridge.DownstreamPipeline != nil {
		attrs["ci.bridge.downstream_pipeline.status"] = bridge.DownstreamPipeline.Status
		attrs["ci.bridge.downstream_pipeline.web_url"] = bridge.DownstreamPipeline.WebUrl
	}

	return &tracepb.ResourceSpans{
		Resource: &resourcepb.Resource{
			Attributes: convertAttributes(map[string]string{
				"service.name": "gitlab_ci.bridge",
			}),
			DroppedAttributesCount: 0,
		},
		ScopeSpans: []*tracepb.ScopeSpans{
			&tracepb.ScopeSpans{
				Scope: &commonpb.InstrumentationScope{
					Name:                   "",
					Version:                "",
					Attributes:             []*commonpb.KeyValue{},
					DroppedAttributesCount: 0,
				},
				Spans: []*tracepb.Span{
					&tracepb.Span{
						TraceId:                traceID,
						SpanId:                 []byte(fmt.Sprint(bridge.Id)),
						TraceState:             "",
						ParentSpanId:           []byte(fmt.Sprint(bridge.Pipeline.Id)),
						Name:                   bridge.Name,
						Kind:                   tracepb.Span_SPAN_KIND_INTERNAL,
						StartTimeUnixNano:      timestampUnixNano(bridge.StartedAt),
						EndTimeUnixNano:        timestampUnixNano(bridge.FinishedAt),
						Attributes:             convertAttributes(attrs),
						DroppedAttributesCount: 0,
						Events:                 []*tracepb.Span_Event{},
						DroppedEventsCount:     0,
						Links:                  []*tracepb.Span_Link{},
						DroppedLinksCount:      0,
						Status:                 convertStatus(bridge.Status),
					},
				},
				SchemaUrl: "",
			},
		},
		SchemaUrl: "",
	}
}

func NewSectionSpan(traceID []byte, section *pb.Section) *tracepb.ResourceSpans {
	return &tracepb.ResourceSpans{
		Resource: &resourcepb.Resource{
			Attributes: convertAttributes(map[string]string{
				"service.name": "gitlab_ci.section",
			}),
			DroppedAttributesCount: 0,
		},
		ScopeSpans: []*tracepb.ScopeSpans{
			&tracepb.ScopeSpans{
				Scope: &commonpb.InstrumentationScope{
					Name:                   "",
					Version:                "",
					Attributes:             []*commonpb.KeyValue{},
					DroppedAttributesCount: 0,
				},
				Spans: []*tracepb.Span{
					&tracepb.Span{
						TraceId:                traceID,
						SpanId:                 []byte(fmt.Sprint(section.Id)),
						TraceState:             "",
						ParentSpanId:           []byte(fmt.Sprint(section.Job.Id)),
						Name:                   section.Name,
						Kind:                   tracepb.Span_SPAN_KIND_INTERNAL,
						StartTimeUnixNano:      timestampUnixNano(section.StartedAt),
						EndTimeUnixNano:        timestampUnixNano(section.FinishedAt),
						Attributes:             convertAttributes(map[string]string{}),
						DroppedAttributesCount: 0,
						Events:                 []*tracepb.Span_Event{},
						DroppedEventsCount:     0,
						Links:                  []*tracepb.Span_Link{},
						DroppedLinksCount:      0,
						Status: &tracepb.Status{
							Message: "",
							Code:    tracepb.Status_STATUS_CODE_UNSET,
						},
					},
				},
				SchemaUrl: "",
			},
		},
		SchemaUrl: "",
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
	switch status {
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

func timestampUnixNano(ts *timestamppb.Timestamp) uint64 {
	if ts == nil {
		return 0
	}
	const nsPerSecond uint64 = 1_000_000_000
	return uint64(ts.Seconds)*nsPerSecond + uint64(ts.Nanos)
}
