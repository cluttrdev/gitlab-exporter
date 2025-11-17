-- projects
CREATE TABLE IF NOT EXISTS projects (
    id INTEGER PRIMARY KEY,
    namespace_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_projects_namespace ON projects(namespace_id);

-- pipelines
CREATE TABLE IF NOT EXISTS pipelines (
    id INTEGER PRIMARY KEY,
    iid INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_pipelines_project ON pipelines(project_id);

-- jobs
CREATE TABLE IF NOT EXISTS jobs (
    id INTEGER PRIMARY KEY,
    pipeline_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_jobs_pipeline ON jobs(project_id, pipeline_id);

-- sections
CREATE TABLE IF NOT EXISTS sections (
    id INTEGER PRIMARY KEY,
    job_id INTEGER NOT NULL,
    pipeline_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_sections_job ON sections(project_id, pipeline_id, job_id);

-- metrics
CREATE TABLE IF NOT EXISTS metrics (
    id TEXT PRIMARY KEY,
    iid INTEGER NOT NULL,
    job_id INTEGER NOT NULL,
    pipeline_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_metrics_job ON metrics(project_id, pipeline_id, job_id);

-- coverage_reports
CREATE TABLE IF NOT EXISTS coverage_reports (
    id TEXT PRIMARY KEY,
    job_id INTEGER NOT NULL,
    pipeline_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_coverage_reports_job ON coverage_reports(project_id, pipeline_id, job_id);

-- coverage_packages
CREATE TABLE IF NOT EXISTS coverage_packages (
    id TEXT PRIMARY KEY,
    report_id TEXT NOT NULL,
    job_id INTEGER NOT NULL,
    pipeline_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_coverage_packages_job ON coverage_packages(project_id, pipeline_id, job_id);
CREATE INDEX IF NOT EXISTS idx_coverage_packages_report ON coverage_packages(report_id);

-- coverage_classes
CREATE TABLE IF NOT EXISTS coverage_classes (
    id TEXT PRIMARY KEY,
    package_id TEXT NOT NULL,
    report_id TEXT NOT NULL,
    job_id INTEGER NOT NULL,
    pipeline_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_coverage_classes_job ON coverage_classes(project_id, pipeline_id, job_id);
CREATE INDEX IF NOT EXISTS idx_coverage_classes_package ON coverage_classes(report_id, package_id);

-- coverage_methods
CREATE TABLE IF NOT EXISTS coverage_methods (
    id TEXT PRIMARY KEY,
    class_id TEXT,
    package_id TEXT,
    report_id TEXT NOT NULL,
    job_id INTEGER NOT NULL,
    pipeline_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_coverage_methods_job ON coverage_methods(project_id, pipeline_id, job_id);
CREATE INDEX IF NOT EXISTS idx_coverage_methods_class ON coverage_methods(report_id, class_id);
CREATE INDEX IF NOT EXISTS idx_coverage_methods_package ON coverage_methods(report_id, package_id);

-- test_reports
CREATE TABLE IF NOT EXISTS test_reports (
    id TEXT PRIMARY KEY,
    job_id INTEGER NOT NULL,
    pipeline_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_test_reports_job ON test_reports(project_id, pipeline_id, job_id);

-- test_suites
CREATE TABLE IF NOT EXISTS test_suites (
    id TEXT PRIMARY KEY,
    test_report_id TEXT NOT NULL,
    job_id INTEGER NOT NULL,
    pipeline_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_test_suites_job ON test_suites(project_id, pipeline_id, job_id);
CREATE INDEX IF NOT EXISTS idx_test_suites_report ON test_suites(test_report_id);

-- test_cases
CREATE TABLE IF NOT EXISTS test_cases (
    id TEXT PRIMARY KEY,
    test_suite_id TEXT NOT NULL,
    test_report_id TEXT NOT NULL,
    job_id INTEGER NOT NULL,
    pipeline_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_test_cases_job ON test_cases(project_id, pipeline_id, job_id);
CREATE INDEX IF NOT EXISTS idx_test_cases_suite ON test_cases(test_report_id, test_suite_id);

-- deployments
CREATE TABLE IF NOT EXISTS deployments (
    id INTEGER PRIMARY KEY,
    iid INTEGER NOT NULL,
    job_id INTEGER NOT NULL,
    pipeline_id INTEGER NOT NULL,
    project_id INTEGER NOT NULL,
    environment_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_deployments_project ON deployments(project_id);
CREATE INDEX IF NOT EXISTS idx_deployments_environment ON deployments(project_id, environment_id);
CREATE INDEX IF NOT EXISTS idx_deployments_job ON deployments(project_id, pipeline_id, job_id);

-- issues
CREATE TABLE IF NOT EXISTS issues (
    id INTEGER PRIMARY KEY,
    iid INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_issues_project ON issues(project_id);

-- merge_requests
CREATE TABLE IF NOT EXISTS merge_requests (
    id INTEGER PRIMARY KEY,
    iid INTEGER NOT NULL,
    project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_merge_requests_project ON merge_requests(project_id);

-- merge_request_note_events
CREATE TABLE IF NOT EXISTS merge_request_note_events (
    id INTEGER PRIMARY KEY,
    merge_request_id INTEGER NOT NULL,
    merge_request_iid INTEGER NOT NULL,
    merge_request_project_id INTEGER NOT NULL,

    _data BLOB NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_merge_request_note_events_mr ON merge_request_note_events(merge_request_project_id, merge_request_id);

-- runners
CREATE TABLE IF NOT EXISTS runners (
    id INTEGER PRIMARY KEY,

    _data BLOB NOT NULL
);

-- traces (OpenTelemetry)
CREATE TABLE IF NOT EXISTS traces (
    Timestamp REAL, -- Unix timestamp
    TraceId TEXT,
    SpanId TEXT PRIMARY KEY,
    ParentSpanId TEXT,
    TraceState TEXT,
    SpanName TEXT,
    SpanKind TEXT,
    ServiceName TEXT,
    ResourceAttributes TEXT, -- JSON object
    ScopeName TEXT,
    ScopeVersion TEXT,
    SpanAttributes TEXT, -- JSON object
    Duration INTEGER,
    StatusCode TEXT,
    StatusMessage TEXT,
    Events TEXT, -- JSON array of {Timestamp, Name, Attributes}
    Links TEXT -- JSON array of {TraceId, SpanId, TraceState, Attributes}
);

CREATE INDEX IF NOT EXISTS idx_traces_traceid ON traces(TraceId);
CREATE INDEX IF NOT EXISTS idx_traces_service ON traces(ServiceName);
CREATE INDEX IF NOT EXISTS idx_traces_timestamp ON traces(Timestamp);

-- trace_view (Grafana)
CREATE VIEW IF NOT EXISTS trace_view AS
SELECT
    TraceId AS `traceID`,
    SpanId AS `spanID`,
    SpanName AS `operationName`,
    ParentSpanId AS `parentSpanID`,
    ServiceName AS `serviceName`,
    Duration / 1000000 AS `duration`,
    Timestamp AS `startTime`,
    SpanAttributes AS `tags`,
    ResourceAttributes AS `serviceTags`,
    Links AS `references`
FROM traces
;

