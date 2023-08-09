package clickhouseclient

import (
	"context"
	"fmt"
)

const (
	dbName string = "gitlab_ci"
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
    id FixedString(16),
    pipeline_id Int64,
    total_time Float64,
    total_count Int64,
    success_count Int64,
    failed_count Int64,
    skipped_count Int64,
    error_count Int64,
    test_suites Nested(
        id FixedString(16),
        name string,
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
    id FixedString(16),
    testreport Tuple(
        id FixedString(16),
        pipeline_id Int64
    ),
    name string,
    total_time Float64,
    total_count Int64,
    success_count Int64,
    failed_count Int64,
    skipped_count Int64,
    error_count Int64,
    test_cases Nested(
        id FixedString(16),
        status string,
        name string
    )
)
ENGINE ReplacingMergeTree()
ORDER BY id
;
    `

	createTestCasesTableSQL = `
CREATE TABLE IF NOT EXISTS %s.%s (
    id FixedString(16),
    testsuite Tuple(
        id FixedString(16),
    ),
    testreport Tuple(
        id FixedString(16),
        pipeline_id Int64
    ),
    status string,
    name string,
    classname string,
    file string,
    execution_time Float64,
    system_output string,
    stack_trace string,
    attachment_url string,
    recent_failures Tuple(
        count Int64,
        base_branch string
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

func createTables(ctx context.Context, client *Client) error {
	if err := client.Conn.Exec(ctx, renderCreatePipelinesTableSQL()); err != nil {
		return fmt.Errorf("exec create pipelines table: %w", err)
	}
	if err := client.Conn.Exec(ctx, renderCreateJobsTableSQL()); err != nil {
		return fmt.Errorf("exec create jobs table: %w", err)
	}
	if err := client.Conn.Exec(ctx, renderCreateSectionsTableSQL()); err != nil {
		return fmt.Errorf("exec create sections table: %w", err)
	}
	if err := client.Conn.Exec(ctx, renderCreateBridgesTableSQL()); err != nil {
		return fmt.Errorf("exec create bridges table: %w", err)
	}
	if err := client.Conn.Exec(ctx, renderCreateTestReportsTableSQL()); err != nil {
		return fmt.Errorf("exec create testreports table: %w", err)
	}
	if err := client.Conn.Exec(ctx, renderCreateTestSuitesTableSQL()); err != nil {
		return fmt.Errorf("exec create testsuites table: %w", err)
	}
	if err := client.Conn.Exec(ctx, renderCreateTestCasesTableSQL()); err != nil {
		return fmt.Errorf("exec create testcases table: %w", err)
	}

	if err := client.Conn.Exec(ctx, renderCreateTracesTableSQL()); err != nil {
		return fmt.Errorf("exec create traces table: %w", err)
	}
	if err := client.Conn.Exec(ctx, renderCreateTraceIdTsTableSQL()); err != nil {
		return fmt.Errorf("exec create traceIdTs table: %w", err)
	}
	if err := client.Conn.Exec(ctx, renderCreateTraceIdTsMaterializedViewSQL()); err != nil {
		return fmt.Errorf("exec create traceIdTs view: %w", err)
	}
	if err := client.Conn.Exec(ctx, renderTraceViewSQL()); err != nil {
		return fmt.Errorf("exec create trace view: %w", err)
	}

	return nil
}

func renderCreatePipelinesTableSQL() string {
	tableName := "pipelines"
	return fmt.Sprintf(createPipelinesTableSQL, dbName, tableName)
}

func renderCreateJobsTableSQL() string {
	tableName := "jobs"
	return fmt.Sprintf(createJobsTableSQL, dbName, tableName)
}

func renderCreateBridgesTableSQL() string {
	tableName := "bridges"
	return fmt.Sprintf(createBridgesTableSQL, dbName, tableName)
}

func renderCreateSectionsTableSQL() string {
	tableName := "sections"
	return fmt.Sprintf(createSectionsTableSQL, dbName, tableName)
}

func renderCreateTestReportsTableSQL() string {
	tableName := "testreports"
	return fmt.Sprintf(createTestReportsTableSQL, dbName, tableName)
}

func renderCreateTestSuitesTableSQL() string {
	tableName := "testsuites"
	return fmt.Sprintf(createTestSuitesTableSQL, dbName, tableName)
}

func renderCreateTestCasesTableSQL() string {
	tableName := "testcases"
	return fmt.Sprintf(createTestCasesTableSQL, dbName, tableName)
}

func renderCreateTracesTableSQL() string {
	tableName := "traces"
	return fmt.Sprintf(createTracesTableSQL, dbName, tableName)
}

func renderCreateTraceIdTsTableSQL() string {
	tableName := "traces"
	return fmt.Sprintf(createTraceIdTsTableSQL, dbName, tableName)
}

func renderCreateTraceIdTsMaterializedViewSQL() string {
	tableName := "traces"
	return fmt.Sprintf(
		createTraceIdTsMaterializedViewSQL,
		dbName, tableName,
		dbName, tableName,
		dbName, tableName,
	)
}

func renderTraceViewSQL() string {
	viewName := "trace_view"
	tableName := "traces"
	return fmt.Sprintf(
		createTraceViewSQL,
		dbName, viewName,
		dbName, tableName,
	)
}
