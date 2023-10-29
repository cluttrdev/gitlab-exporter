package models

import (
	"fmt"
	"time"

	"golang.org/x/exp/slices"
)

type Span struct {
	TraceID      string
	SpanID       string
	TraceState   string
	ParentSpanID string
	Name         string
	Kind         SpanKind
	StartTime    uint64
	EndTime      uint64
	Attributes   map[string]string
	Events       []SpanEvent
	Links        []SpanLink
	Status       SpanStatus
	Resource     *Resource
}

type SpanKind int32

const (
	SpanKindUnspecified SpanKind = 0
	SpanKindInternal    SpanKind = 1
	SpanKindServer      SpanKind = 2
	SpanKindClient      SpanKind = 3
	SpanKindProducer    SpanKind = 4
	SpanKindConsumer    SpanKind = 5
)

var (
	SpanKindNames = map[int32]string{
		0: "Unspecified",
		1: "Internal",
		2: "Server",
		3: "Client",
		4: "Producer",
		5: "Consumer",
	}

	SpanKindValues = map[string]int32{
		"Unspecified": 0,
		"Internal":    1,
		"Server":      2,
		"Client":      3,
		"Producer":    4,
		"Consumer":    5,
	}
)

func (sk *SpanKind) Name() string {
	return SpanKindNames[int32(*sk)]
}

type SpanEvent struct {
	Time       uint64
	Name       string
	Attributes map[string]string
}

type SpanLink struct {
	TraceID    string
	SpanID     string
	TraceState string
	Attributes map[string]string
}

type SpanStatus struct {
	Message string
	Code    StatusCode
}

type StatusCode int32

const (
	StatusCodeUnset StatusCode = 0
	StatusCodeOk    StatusCode = 1
	StatusCodeError StatusCode = 2
)

var (
	StatusCodeNames = map[int32]string{
		0: "Unset",
		1: "Ok",
		2: "Error",
	}

	StatusCodeValues = map[string]int32{
		"Unset": 0,
		"Ok":    1,
		"Error": 2,
	}
)

func (sc *StatusCode) Name() string {
	return StatusCodeNames[int32(*sc)]
}

type Resource struct {
	Attributes map[string]string
}

func unixNano(t *time.Time) uint64 {
	if t == nil {
		return 0
	}
	return uint64(t.UnixNano())
}

func spanStatusCode(gitlabStatus string) (status StatusCode) {
	status = StatusCodeUnset
	if gitlabStatus == "success" {
		status = StatusCodeOk
	} else if gitlabStatus == "failed" {
		status = StatusCodeError
	}
	return
}

func NewPipelineSpan(traceID string, parentID string, pipeline *Pipeline) *Span {
	return &Span{
		TraceID:      traceID,
		SpanID:       fmt.Sprintf("%d", pipeline.ID),
		TraceState:   "",
		ParentSpanID: parentID,
		Name:         pipeline.Ref,
		Kind:         SpanKindInternal,
		StartTime:    unixNano(pipeline.StartedAt),
		EndTime:      unixNano(pipeline.FinishedAt),
		Attributes:   pipelineSpanAttributes(pipeline),
		Events:       []SpanEvent{},
		Links:        []SpanLink{},
		Status: SpanStatus{
			Message: pipeline.Status,
			Code:    spanStatusCode(pipeline.Status),
		},
		Resource: &Resource{
			Attributes: map[string]string{
				"service.name": "gitlab_ci.pipeline",
			},
		},
	}
}

func NewJobSpan(traceID string, job *Job) *Span {
	return &Span{
		TraceID:      traceID,
		SpanID:       fmt.Sprintf("%d", job.ID),
		TraceState:   "",
		ParentSpanID: fmt.Sprintf("%d", job.Pipeline.ID),
		Name:         job.Name,
		Kind:         SpanKindInternal,
		StartTime:    unixNano(job.StartedAt),
		EndTime:      unixNano(job.FinishedAt),
		Attributes:   jobSpanAttributes(job),
		Events:       []SpanEvent{},
		Links:        []SpanLink{},
		Status: SpanStatus{
			Message: job.Status,
			Code:    spanStatusCode(job.Status),
		},
		Resource: &Resource{
			Attributes: map[string]string{
				"service.name": "gitlab_ci.job",
			},
		},
	}
}

func NewBridgeSpan(traceID string, bridge *Bridge) *Span {
	return &Span{
		TraceID:      traceID,
		SpanID:       fmt.Sprintf("%d", bridge.ID),
		TraceState:   "",
		ParentSpanID: fmt.Sprintf("%d", bridge.Pipeline.ID),
		Name:         bridge.Name,
		Kind:         SpanKindInternal,
		StartTime:    unixNano(bridge.StartedAt),
		EndTime:      unixNano(bridge.FinishedAt),
		Attributes:   bridgeSpanAttributes(bridge),
		Events:       []SpanEvent{},
		Links:        []SpanLink{},
		Status: SpanStatus{
			Message: bridge.Status,
			Code:    spanStatusCode(bridge.Status),
		},
		Resource: &Resource{
			Attributes: map[string]string{
				"service.name": "gitlab_ci.bridge",
			},
		},
	}
}

func NewSectionSpan(traceID string, section *Section) *Span {
	return &Span{
		TraceID:      traceID,
		SpanID:       fmt.Sprintf("%d", section.ID),
		TraceState:   "",
		ParentSpanID: fmt.Sprintf("%d", section.Job.ID),
		Name:         section.Name,
		Kind:         SpanKindInternal,
		StartTime:    unixNano(section.StartedAt),
		EndTime:      unixNano(section.FinishedAt),
		Attributes:   sectionSpanAttributes(section),
		Events:       []SpanEvent{},
		Links:        []SpanLink{},
		Status: SpanStatus{
			Message: "",
			Code:    StatusCodeUnset,
		},
		Resource: &Resource{
			Attributes: map[string]string{
				"service.name": "gitlab_ci.section",
			},
		},
	}
}

func NewPipelineHierarchyTrace(traceID string, parentID string, ph *PipelineHierarchy) []*Span {
	if traceID == "" {
		traceID = fmt.Sprintf("%d", ph.Pipeline.ID)
	}

	var trace = []*Span{}

	trace = append(trace, NewPipelineSpan(traceID, parentID, ph.Pipeline))

	for _, job := range ph.Jobs {
		trace = append(trace, NewJobSpan(traceID, job))
	}

	for _, section := range ph.Sections {
		trace = append(trace, NewSectionSpan(traceID, section))
	}

	for _, bridge := range ph.Bridges {
		span := NewBridgeSpan(traceID, bridge)
		trace = append(trace, span)

		index := slices.IndexFunc(ph.DownstreamPipelines, func(dph *PipelineHierarchy) bool {
			return dph.Pipeline.ID == bridge.DownstreamPipeline.ID
		})
		if index > -1 {
			dph := ph.DownstreamPipelines[index]
			trace = append(trace, NewPipelineHierarchyTrace(span.TraceID, span.SpanID, dph)...)
		}
	}

	return trace
}

func pipelineSpanAttributes(pipeline *Pipeline) map[string]string {
	attr := map[string]string{}

	if pipeline != nil {
		attr["ci.pipeline.status"] = pipeline.Status
		attr["ci.pipeline.web_url"] = pipeline.WebURL
	}

	return attr
}

func jobSpanAttributes(job *Job) map[string]string {
	attr := map[string]string{}

	if job != nil {
		attr["ci.job.status"] = job.Status
		attr["ci.job.web_url"] = job.WebURL
	}

	return attr
}

func sectionSpanAttributes(section *Section) map[string]string {
	attr := map[string]string{}

	return attr
}

func bridgeSpanAttributes(bridge *Bridge) map[string]string {
	attr := map[string]string{}

	if bridge != nil {
		attr["ci.bridge.status"] = bridge.Status
		attr["ci.bridge.web_url"] = bridge.WebURL

		if bridge.DownstreamPipeline != nil {
			attr["ci.bridge.downstream_pipeline.status"] = bridge.DownstreamPipeline.Status
			attr["ci.bridge.downstream_pipeline.web_url"] = bridge.DownstreamPipeline.WebURL
		}
	}

	return attr
}
