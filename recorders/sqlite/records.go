package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"go.cluttr.dev/gitlab-exporter/protobuf/servicepb"
	"google.golang.org/protobuf/proto"
)

// numFields returns the number of exported fields in a struct
func numFields(s any) int {
	fields := reflect.VisibleFields(reflect.TypeOf(s))
	count := 0
	for _, field := range fields {
		if field.IsExported() {
			count += 1
		}
	}
	return count
}

// structToSlice converts a struct to a slice of its exported field values
func structToSlice(s any) []any {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	fields := reflect.VisibleFields(reflect.TypeOf(s))
	result := make([]any, 0, len(fields))

	for _, field := range fields {
		if field.IsExported() {
			fv := v.FieldByIndex(field.Index)
			result = append(result, fv.Interface())
		}
	}
	return result
}

// record is a generic function to record protobuf messages into a specified table.
// It takes a conversion function to transform protobuf messages into the appropriate struct type.
func record[P proto.Message, T any](ctx context.Context, db *sql.DB, table string, msgs []P, convert func(p P) (T, error)) (int32, error) {
	nrows := len(msgs)
	var t T
	ncols := numFields(t)
	stmt, err := prepareBatchInsert(ctx, db, table, nrows, ncols)
	if err != nil {
		return 0, fmt.Errorf("prepare statement: %w", err)
	}

	vals := make([]any, 0, nrows*ncols)
	for _, msg := range msgs {
		val, err := convert(msg)
		if err != nil {
			return 0, fmt.Errorf("convert message: %w", err)
		}
		vals = append(vals, structToSlice(val)...)
	}
	if len(vals) != nrows*ncols {
		return 0, fmt.Errorf("invalid number of values: got %d, expected %d", len(vals), nrows*ncols)
	}

	err = withTransaction(ctx, db, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.StmtContext(ctx, stmt).ExecContext(ctx, vals...)
		if err != nil {
			return fmt.Errorf("exec: %w", err)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}

	return int32(nrows), nil
}

func (r *Recorder) RecordCoverageClasses(ctx context.Context, req *servicepb.RecordCoverageClassesRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "coverage_classes", req.Data, ConvertCoverageClass)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordCoverageMethods(ctx context.Context, req *servicepb.RecordCoverageMethodsRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "coverage_methods", req.Data, ConvertCoverageMethod)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordCoveragePackages(ctx context.Context, req *servicepb.RecordCoveragePackagesRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "coverage_packages", req.Data, ConvertCoveragePackage)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordCoverageReports(ctx context.Context, req *servicepb.RecordCoverageReportsRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "coverage_reports", req.Data, ConvertCoverageReport)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordDeployments(ctx context.Context, req *servicepb.RecordDeploymentsRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "deployments", req.Data, ConvertDeployment)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordIssues(ctx context.Context, req *servicepb.RecordIssuesRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "issues", req.Data, ConvertIssue)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordJobs(ctx context.Context, req *servicepb.RecordJobsRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "jobs", req.Data, ConvertJob)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordMergeRequestNoteEvents(ctx context.Context, req *servicepb.RecordMergeRequestNoteEventsRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "merge_request_note_events", req.Data, ConvertMergeRequestNoteEvent)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordMergeRequests(ctx context.Context, req *servicepb.RecordMergeRequestsRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "merge_requests", req.Data, ConvertMergeRequest)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordMetrics(ctx context.Context, req *servicepb.RecordMetricsRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "metrics", req.Data, ConvertMetric)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordPipelines(ctx context.Context, req *servicepb.RecordPipelinesRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "pipelines", req.Data, ConvertPipeline)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordProjects(ctx context.Context, req *servicepb.RecordProjectsRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "projects", req.Data, ConvertProject)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordRunners(ctx context.Context, req *servicepb.RecordRunnersRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "runners", req.Data, ConvertRunner)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordSections(ctx context.Context, req *servicepb.RecordSectionsRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "sections", req.Data, ConvertSection)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordTestCases(ctx context.Context, req *servicepb.RecordTestCasesRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "test_cases", req.Data, ConvertTestCase)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordTestReports(ctx context.Context, req *servicepb.RecordTestReportsRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "test_reports", req.Data, ConvertTestReport)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordTestSuites(ctx context.Context, req *servicepb.RecordTestSuitesRequest) (*servicepb.RecordSummary, error) {
	n, err := record(ctx, r.db, "test_suites", req.Data, ConvertTestSuite)
	return &servicepb.RecordSummary{
		RecordedCount: n,
	}, err
}

func (r *Recorder) RecordTraces(ctx context.Context, req *servicepb.RecordTracesRequest) (*servicepb.RecordSummary, error) {
	var spans []TraceSpan
	for _, msg := range req.Data {
		traceSpans, err := ConvertTrace(msg)
		if err != nil {
			return &servicepb.RecordSummary{}, fmt.Errorf("convert trace: %w", err)
		}
		spans = append(spans, traceSpans...)
	}
	if len(spans) == 0 {
		return &servicepb.RecordSummary{}, nil
	}

	nrows := len(spans)
	ncols := numFields(spans[0])
	stmt, err := prepareBatchInsert(ctx, r.db, "traces", nrows, ncols)
	if err != nil {
		return &servicepb.RecordSummary{}, fmt.Errorf("prepare statement: %w", err)
	}

	vals := make([]any, 0, nrows*ncols)
	for _, val := range spans {
		vals = append(vals, structToSlice(val)...)
	}
	if len(vals) != nrows*ncols {
		return &servicepb.RecordSummary{}, fmt.Errorf("invalid number of values: got %d, expected %d", len(vals), nrows*ncols)
	}

	err = withTransaction(ctx, r.db, func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.StmtContext(ctx, stmt).ExecContext(ctx, vals...)
		if err != nil {
			return fmt.Errorf("exec: %w", err)
		}
		return nil
	})
	if err != nil {
		return &servicepb.RecordSummary{}, err
	}

	return &servicepb.RecordSummary{
		RecordedCount: int32(nrows),
	}, nil
}
