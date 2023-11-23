package clickhouse

import (
	"context"
	"fmt"
)

const (
	defaultDBName string = "gitlab_ci"
)

const (
	createPipelinesTableSQL = `
CREATE TABLE IF NOT EXISTS %s.%s (
    id Int64,
    iid Int64,
    project_id Int64,
    status String,
    source String,
    ref String,
    sha String,
    before_sha String,
    tag Bool,
    yaml_errors String,
    created_at Float64,
    updated_at Float64,
    started_at Float64,
    finished_at Float64,
    committed_at Float64,
    duration Float64,
    queued_duration Float64,
    coverage Float64,
    web_url String
)
ENGINE ReplacingMergeTree(updated_at)
ORDER BY id
;
    `

	createJobsTableSQL = `
CREATE TABLE IF NOT EXISTS %s.%s (
    coverage Float64,
    allow_failure Bool,
    created_at Float64,
    started_at Float64,
    finished_at Float64,
    erased_at Float64,
    duration Float64,
    queued_duration Float64,
    tag_list Array(String),
    id Int64,
    name String,
    pipeline Tuple(
        id Int64,
        project_id Int64,
        ref String,
        sha String,
        status String
    ),
    ref String,
    stage String,
    status String,
    failure_reason String,
    tag Bool,
    web_url String
)
ENGINE ReplacingMergeTree()
ORDER BY id
;
    `

	createBridgesTableSQL = `
CREATE TABLE IF NOT EXISTS %s.%s (
    coverage Float64,
    allow_failure Bool,
    created_at Float64,
    started_at Float64,
    finished_at Float64,
    erased_at Float64,
    duration Float64,
    queued_duration Float64,
    id Int64,
    name String,
    pipeline Tuple(
        id Int64,
        iid Int64,
        project_id Int64,
        status String,
        source String,
        ref String,
        sha String,
        web_url String,
        created_at Float64,
        updated_at Float64
    ),
    ref String,
    stage String,
    status String,
    failure_reason String,
    tag Bool,
    web_url String,
    downstream_pipeline Tuple(
        id Int64,
        iid Int64,
        project_id Int64,
        status String,
        source String,
        ref String,
        sha String,
        web_url String,
        created_at Float64,
        updated_at Float64
    )
)
ENGINE ReplacingMergeTree()
ORDER BY id
;
    `

	createSectionsTableSQL = `
CREATE TABLE IF NOT EXISTS %s.%s (
    id Int64,
    name String,
    job Tuple(
        id Int64,
        name String,
        status String
    ),
    pipeline Tuple(
        id Int64,
        project_id Int64,
        ref String,
        sha String,
        status String
    ),
    started_at Float64,
    finished_at Float64,
    duration Float64
)
ENGINE ReplacingMergeTree()
ORDER BY id
;
    `

	createTestReportsTableSQL = `
CREATE TABLE IF NOT EXISTS %s.%s (
    id Int64,
    pipeline_id Int64,
    total_time Float64,
    total_count Int64,
    success_count Int64,
    failed_count Int64,
    skipped_count Int64,
    error_count Int64,
    test_suites Nested(
        id Int64,
        name String,
        total_time Float64,
        total_count Int64
    )
)
ENGINE ReplacingMergeTree()
ORDER BY id
;
    `

	createTestSuitesTableSQL = `
CREATE TABLE IF NOT EXISTS %s.%s (
    id Int64,
    testreport Tuple(
        id Int64,
        pipeline_id Int64
    ),
    name String,
    total_time Float64,
    total_count Int64,
    success_count Int64,
    failed_count Int64,
    skipped_count Int64,
    error_count Int64,
    test_cases Nested(
        id Int64,
        status String,
        name String
    )
)
ENGINE ReplacingMergeTree()
ORDER BY id
;
    `

	createTestCasesTableSQL = `
CREATE TABLE IF NOT EXISTS %s.%s (
    id Int64,
    testsuite Tuple(
        id Int64
    ),
    testreport Tuple(
        id Int64,
        pipeline_id Int64
    ),
    status String,
    name String,
    classname String,
    file String,
    execution_time Float64,
    system_output String,
    stack_trace String,
    attachment_url String,
    recent_failures Tuple(
        count Int64,
        base_branch String
    )
)
ENGINE ReplacingMergeTree()
ORDER BY id
;
    `
)

const (
	// OpenTelemetry Traces
	// schemas taken from https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/exporter/clickhouseexporter/exporter_traces.go

	createTracesTableSQL = `
CREATE TABLE IF NOT EXISTS %s.%s (
     Timestamp DateTime64(9) CODEC(Delta, ZSTD(1)),
     TraceId String CODEC(ZSTD(1)),
     SpanId String CODEC(ZSTD(1)),
     ParentSpanId String CODEC(ZSTD(1)),
     TraceState String CODEC(ZSTD(1)),
     SpanName LowCardinality(String) CODEC(ZSTD(1)),
     SpanKind LowCardinality(String) CODEC(ZSTD(1)),
     ServiceName LowCardinality(String) CODEC(ZSTD(1)),
     ResourceAttributes Map(LowCardinality(String), String) CODEC(ZSTD(1)),
     ScopeName String CODEC(ZSTD(1)),
     ScopeVersion String CODEC(ZSTD(1)),
     SpanAttributes Map(LowCardinality(String), String) CODEC(ZSTD(1)),
     Duration Int64 CODEC(ZSTD(1)),
     StatusCode LowCardinality(String) CODEC(ZSTD(1)),
     StatusMessage String CODEC(ZSTD(1)),
     Events Nested (
         Timestamp DateTime64(9),
         Name LowCardinality(String),
         Attributes Map(LowCardinality(String), String)
     ) CODEC(ZSTD(1)),
     Links Nested (
         TraceId String,
         SpanId String,
         TraceState String,
         Attributes Map(LowCardinality(String), String)
     ) CODEC(ZSTD(1)),
     INDEX idx_trace_id TraceId TYPE bloom_filter(0.001) GRANULARITY 1,
     INDEX idx_res_attr_key mapKeys(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_res_attr_value mapValues(ResourceAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_span_attr_key mapKeys(SpanAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_span_attr_value mapValues(SpanAttributes) TYPE bloom_filter(0.01) GRANULARITY 1,
     INDEX idx_duration Duration TYPE minmax GRANULARITY 1
) ENGINE MergeTree()
PARTITION BY toDate(Timestamp)
ORDER BY (ServiceName, SpanName, toUnixTimestamp(Timestamp), TraceId)
SETTINGS index_granularity=8192, ttl_only_drop_parts = 1
;
    `

	createTraceIdTsTableSQL = `
CREATE TABLE IF NOT EXISTS %s.%s_trace_id_ts (
     TraceId String CODEC(ZSTD(1)),
     Start DateTime64(9) CODEC(Delta, ZSTD(1)),
     End DateTime64(9) CODEC(Delta, ZSTD(1)),
     INDEX idx_trace_id TraceId TYPE bloom_filter(0.01) GRANULARITY 1
) ENGINE MergeTree()
ORDER BY (TraceId, toUnixTimestamp(Start))
SETTINGS index_granularity=8192
;
    `

	createTraceIdTsMaterializedViewSQL = `
CREATE MATERIALIZED VIEW IF NOT EXISTS %s.%s_trace_id_ts_mv
TO %s.%s_trace_id_ts
AS SELECT
    TraceId,
    min(Timestamp) as Start,
    max(Timestamp) as End
FROM %s.%s
WHERE TraceId != ''
GROUP BY TraceId
;
    `

	createTraceViewSQL = `
CREATE VIEW IF NOT EXISTS %s.%s AS
SELECT
    TraceId AS traceID,
    SpanId AS spanID,
    SpanName AS operationName,
    ParentSpanId AS parentSpanID,
    ServiceName AS serviceName,
    Duration / 1000000 AS duration,
    Timestamp AS startTime,
    arrayMap(key -> map('key', key, 'value', SpanAttributes[key]), mapKeys(SpanAttributes)) AS tags,
    arrayMap(key -> map('key', key, 'value', ResourceAttributes[key]), mapKeys(ResourceAttributes)) AS serviceTags
FROM %s.%s
WHERE TraceId = {trace_id:String}
;
    `
)

func createTables(ctx context.Context, db string, client *Client) error {
	if err := client.Exec(ctx, renderCreatePipelinesTableSQL(db)); err != nil {
		return fmt.Errorf("exec create pipelines table: %w", err)
	}
	if err := client.Exec(ctx, renderCreateJobsTableSQL(db)); err != nil {
		return fmt.Errorf("exec create jobs table: %w", err)
	}
	if err := client.Exec(ctx, renderCreateSectionsTableSQL(db)); err != nil {
		return fmt.Errorf("exec create sections table: %w", err)
	}
	if err := client.Exec(ctx, renderCreateBridgesTableSQL(db)); err != nil {
		return fmt.Errorf("exec create bridges table: %w", err)
	}
	if err := client.Exec(ctx, renderCreateTestReportsTableSQL(db)); err != nil {
		return fmt.Errorf("exec create testreports table: %w", err)
	}
	if err := client.Exec(ctx, renderCreateTestSuitesTableSQL(db)); err != nil {
		return fmt.Errorf("exec create testsuites table: %w", err)
	}
	if err := client.Exec(ctx, renderCreateTestCasesTableSQL(db)); err != nil {
		return fmt.Errorf("exec create testcases table: %w", err)
	}

	if err := client.Exec(ctx, renderCreateTracesTableSQL(db)); err != nil {
		return fmt.Errorf("exec create traces table: %w", err)
	}
	if err := client.Exec(ctx, renderCreateTraceIdTsTableSQL(db)); err != nil {
		return fmt.Errorf("exec create traceIdTs table: %w", err)
	}
	if err := client.Exec(ctx, renderCreateTraceIdTsMaterializedViewSQL(db)); err != nil {
		return fmt.Errorf("exec create traceIdTs view: %w", err)
	}
	if err := client.Exec(ctx, renderTraceViewSQL(db)); err != nil {
		return fmt.Errorf("exec create trace view: %w", err)
	}

	return nil
}

func renderCreatePipelinesTableSQL(db string) string {
	const tableName string = "pipelines"
	return fmt.Sprintf(createPipelinesTableSQL, db, tableName)
}

func renderCreateJobsTableSQL(db string) string {
	const tableName string = "jobs"
	return fmt.Sprintf(createJobsTableSQL, db, tableName)
}

func renderCreateBridgesTableSQL(db string) string {
	const tableName string = "bridges"
	return fmt.Sprintf(createBridgesTableSQL, db, tableName)
}

func renderCreateSectionsTableSQL(db string) string {
	const tableName string = "sections"
	return fmt.Sprintf(createSectionsTableSQL, db, tableName)
}

func renderCreateTestReportsTableSQL(db string) string {
	const tableName string = "testreports"
	return fmt.Sprintf(createTestReportsTableSQL, db, tableName)
}

func renderCreateTestSuitesTableSQL(db string) string {
	const tableName string = "testsuites"
	return fmt.Sprintf(createTestSuitesTableSQL, db, tableName)
}

func renderCreateTestCasesTableSQL(db string) string {
	const tableName string = "testcases"
	return fmt.Sprintf(createTestCasesTableSQL, db, tableName)
}

func renderCreateTracesTableSQL(db string) string {
	const tableName string = "traces"
	return fmt.Sprintf(createTracesTableSQL, db, tableName)
}

func renderCreateTraceIdTsTableSQL(db string) string {
	const tableName string = "traces"
	return fmt.Sprintf(createTraceIdTsTableSQL, db, tableName)
}

func renderCreateTraceIdTsMaterializedViewSQL(db string) string {
	const tableName string = "traces"
	return fmt.Sprintf(
		createTraceIdTsMaterializedViewSQL,
		db, tableName,
		db, tableName,
		db, tableName,
	)
}

func renderTraceViewSQL(db string) string {
	viewName := "trace_view"
	const tableName string = "traces"
	return fmt.Sprintf(
		createTraceViewSQL,
		db, viewName,
		db, tableName,
	)
}
