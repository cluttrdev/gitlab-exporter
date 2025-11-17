package sqlite

import (
	"encoding/json"
	"fmt"

	otlp_comonpb "go.opentelemetry.io/proto/otlp/common/v1"
	otlp_tracepb "go.opentelemetry.io/proto/otlp/trace/v1"

	"go.cluttr.dev/gitlab-exporter/protobuf/typespb"
)

func ConvertCoverageClass(msg *typespb.CoverageClass) (CoverageClass, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return CoverageClass{}, err
	}

	return CoverageClass{
		Id:        msg.GetId(),
		PackageId: msg.GetPackage().GetId(),
		ReportId:  msg.GetPackage().GetReport().GetId(),

		JobId:      int(msg.GetPackage().GetReport().GetJob().GetId()),
		PipelineId: int(msg.GetPackage().GetReport().GetJob().GetPipeline().GetId()),
		ProjectId:  int(msg.GetPackage().GetReport().GetJob().GetPipeline().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertCoverageMethod(msg *typespb.CoverageMethod) (CoverageMethod, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return CoverageMethod{}, err
	}

	return CoverageMethod{
		Id:        msg.GetId(),
		ClassId:   msg.GetClass().GetId(),
		PackageId: msg.GetClass().GetPackage().GetId(),
		ReportId:  msg.GetClass().GetPackage().GetReport().GetId(),

		JobId:      int(msg.GetClass().GetPackage().GetReport().GetJob().GetId()),
		PipelineId: int(msg.GetClass().GetPackage().GetReport().GetJob().GetPipeline().GetId()),
		ProjectId:  int(msg.GetClass().GetPackage().GetReport().GetJob().GetPipeline().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertCoveragePackage(msg *typespb.CoveragePackage) (CoveragePackage, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return CoveragePackage{}, err
	}

	return CoveragePackage{
		Id:       msg.GetId(),
		ReportId: msg.GetReport().GetId(),

		JobId:      int(msg.GetReport().GetJob().GetId()),
		PipelineId: int(msg.GetReport().GetJob().GetPipeline().GetId()),
		ProjectId:  int(msg.GetReport().GetJob().GetPipeline().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertCoverageReport(msg *typespb.CoverageReport) (CoverageReport, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return CoverageReport{}, err
	}

	return CoverageReport{
		Id:         msg.GetId(),
		JobId:      int(msg.GetJob().GetId()),
		PipelineId: int(msg.GetJob().GetPipeline().GetId()),
		ProjectId:  int(msg.GetJob().GetPipeline().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertDeployment(msg *typespb.Deployment) (Deployment, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return Deployment{}, err
	}

	return Deployment{
		Id:            int(msg.GetId()),
		Iid:           int(msg.GetIid()),
		EnvironmentId: int(msg.GetEnvironment().GetId()),

		JobId:      int(msg.GetJob().GetId()),
		PipelineId: int(msg.GetJob().GetPipeline().GetId()),
		ProjectId:  int(msg.GetJob().GetPipeline().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertIssue(msg *typespb.Issue) (Issue, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return Issue{}, err
	}

	return Issue{
		Id:        int(msg.GetId()),
		Iid:       int(msg.GetIid()),
		ProjectId: int(msg.GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertJob(msg *typespb.Job) (Job, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return Job{}, err
	}

	return Job{
		Id:         int(msg.GetId()),
		PipelineId: int(msg.GetPipeline().GetId()),
		ProjectId:  int(msg.GetPipeline().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertMergeRequest(msg *typespb.MergeRequest) (MergeRequest, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return MergeRequest{}, err
	}

	return MergeRequest{
		Id:        int(msg.GetId()),
		Iid:       int(msg.GetIid()),
		ProjectId: int(msg.GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertMergeRequestNoteEvent(msg *typespb.MergeRequestNoteEvent) (MergeRequestNoteEvent, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return MergeRequestNoteEvent{}, err
	}

	return MergeRequestNoteEvent{
		Id:                    int(msg.GetId()),
		MergeRequestId:        int(msg.GetMergeRequest().GetId()),
		MergeRequestIid:       int(msg.GetMergeRequest().GetIid()),
		MergeRequestProjectId: int(msg.GetMergeRequest().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertMetric(msg *typespb.Metric) (Metric, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return Metric{}, err
	}

	return Metric{
		Id:         string(msg.GetId()),
		Iid:        int(msg.GetIid()),
		JobId:      int(msg.GetJob().GetId()),
		PipelineId: int(msg.GetJob().GetPipeline().GetId()),
		ProjectId:  int(msg.GetJob().GetPipeline().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertPipeline(msg *typespb.Pipeline) (Pipeline, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return Pipeline{}, err
	}

	return Pipeline{
		Id:        int(msg.GetId()),
		Iid:       int(msg.GetIid()),
		ProjectId: int(msg.GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertProject(msg *typespb.Project) (Project, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return Project{}, err
	}

	return Project{
		Id:          int(msg.GetId()),
		NamespaceId: int(msg.GetNamespace().GetId()),

		Data: data,
	}, nil
}

func ConvertRunner(msg *typespb.Runner) (Runner, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return Runner{}, err
	}

	return Runner{
		Id: int(msg.GetId()),

		Data: data,
	}, nil
}

func ConvertSection(msg *typespb.Section) (Section, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return Section{}, err
	}

	return Section{
		Id:         int(msg.GetId()),
		JobId:      int(msg.GetJob().GetId()),
		PipelineId: int(msg.GetJob().GetPipeline().GetId()),
		ProjectId:  int(msg.GetJob().GetPipeline().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertTestCase(msg *typespb.TestCase) (TestCase, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return TestCase{}, err
	}

	return TestCase{
		Id:           msg.GetId(),
		TestSuiteId:  msg.GetTestSuite().GetId(),
		TestReportId: msg.GetTestSuite().GetTestReport().GetId(),

		JobId:      int(msg.GetTestSuite().GetTestReport().GetJob().GetId()),
		PipelineId: int(msg.GetTestSuite().GetTestReport().GetJob().GetPipeline().GetId()),
		ProjectId:  int(msg.GetTestSuite().GetTestReport().GetJob().GetPipeline().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertTestReport(msg *typespb.TestReport) (TestReport, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return TestReport{}, err
	}

	return TestReport{
		Id: msg.GetId(),

		JobId:      int(msg.GetJob().GetId()),
		PipelineId: int(msg.GetJob().GetPipeline().GetId()),
		ProjectId:  int(msg.GetJob().GetPipeline().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertTestSuite(msg *typespb.TestSuite) (TestSuite, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return TestSuite{}, err
	}

	return TestSuite{
		Id:           msg.GetId(),
		TestReportId: msg.GetTestReport().GetId(),

		JobId:      int(msg.GetTestReport().GetJob().GetId()),
		PipelineId: int(msg.GetTestReport().GetJob().GetPipeline().GetId()),
		ProjectId:  int(msg.GetTestReport().GetJob().GetPipeline().GetProject().GetId()),

		Data: data,
	}, nil
}

func ConvertTrace(msg *typespb.Trace) ([]TraceSpan, error) {
	var spans []TraceSpan

	for _, resourceSpans := range msg.Data.ResourceSpans {
		resourceAttrs := convertAttributes(resourceSpans.Resource.Attributes)
		serviceName := resourceAttrs["service.name"]
		resourceAttrsData, err := json.Marshal(resourceAttrs)
		if err != nil {
			return nil, fmt.Errorf("convert resource attributes: %w", err)
		}

		for _, scopeSpans := range resourceSpans.ScopeSpans {
			scopeName := scopeSpans.Scope.Name
			scopeVersion := scopeSpans.Scope.Version
			for _, span := range scopeSpans.Spans {
				spanAttrs := convertAttributes(span.Attributes)
				spanAttrsData, err := json.Marshal(spanAttrs)
				if err != nil {
					return nil, fmt.Errorf("convert span attributes: %w", err)
				}

				spanEvents := convertEvents(span.Events)
				spanEventsData, err := json.Marshal(spanEvents)
				if err != nil {
					return nil, fmt.Errorf("convert span events: %w", err)
				}

				spanLinks := convertLinks(span.Links)
				spanLinksData, err := json.Marshal(spanLinks)
				if err != nil {
					return nil, fmt.Errorf("convert span links: %w", err)
				}

				spans = append(spans, TraceSpan{
					Timestamp:          span.StartTimeUnixNano,
					TraceId:            span.TraceId,
					SpanId:             span.SpanId,
					ParentSpanId:       span.ParentSpanId,
					TraceState:         span.TraceState,
					SpanName:           span.Name,
					SpanKind:           span.Kind.String(),
					ServiceName:        serviceName,
					ResourceAttributes: resourceAttrsData,
					ScopeName:          scopeName,
					ScopeVersion:       scopeVersion,
					SpanAttributes:     spanAttrsData,
					Duration:           int64(span.EndTimeUnixNano) - int64(span.StartTimeUnixNano),
					StatusCode:         int32(span.GetStatus().GetCode()),
					StatusMessage:      span.GetStatus().GetMessage(),
					Events:             spanEventsData,
					Links:              spanLinksData,
				})
			}
		}
	}

	return spans, nil
}

func convertAttributes(list []*otlp_comonpb.KeyValue) map[string]string {
	attrs := make(map[string]string)

	for _, attr := range list {
		value, ok := attr.GetValue().Value.(*otlp_comonpb.AnyValue_StringValue)
		if ok {
			attrs[attr.Key] = value.StringValue
		}
	}

	return attrs
}

func convertEvents(events []*otlp_tracepb.Span_Event) []map[string]any {
	var eventMaps []map[string]any
	for _, event := range events {
		eventMaps = append(eventMaps, map[string]any{
			"Timestamp":  event.TimeUnixNano,
			"Name":       event.Name,
			"Attributes": convertAttributes(event.Attributes),
		})
	}
	return eventMaps
}

func convertLinks(links []*otlp_tracepb.Span_Link) []map[string]any {
	var linkMaps []map[string]any
	for _, link := range links {
		linkMaps = append(linkMaps, map[string]any{
			"TraceId":    string(link.TraceId),
			"SpanId":     string(link.SpanId),
			"TraceState": link.TraceState,
			"Attributes": convertAttributes(link.Attributes),
		})
	}
	return linkMaps
}
