-- coverage_reports
CREATE TABLE IF NOT EXISTS coverage_reports (
    id String,
    job_id Int64,
    pipeline_id Int64,
    project_id Int64,

    line_rate Float32,
    lines_covered Int32,
    lines_valid Int32,

    branch_rate Float32,
    branches_covered Int32,
    branches_valid Int32,

    complexity Float32,

    version String,
    timestamp Float64,

    source_paths Array(String)
)
ENGINE = ReplacingMergeTree()
PARTITION BY toStartOfMonth(toDateTime(timestamp))
ORDER BY (project_id, pipeline_id, job_id, id)
;


-- coverage_reports_in
CREATE TABLE IF NOT EXISTS coverage_reports_in AS coverage_reports ENGINE = Null;

-- coverage_reports_mv
CREATE MATERIALIZED VIEW IF NOT EXISTS coverage_reports_mv TO coverage_reports
AS
SELECT * FROM coverage_reports_in
WHERE id NOT IN (
    SELECT id FROM coverage_reports
    WHERE job_id IN (
        SELECT DISTINCT job_id FROM coverage_reports_in
    )
)
;

-- coverage_packages
CREATE TABLE IF NOT EXISTS coverage_packages (
    id String,
    report_id String,
    job_id Int64,
    pipeline_id Int64,
    project_id Int64,

    name String,

    line_rate Float32,
    branch_rate Float32,
    complexity Float32,
)
ENGINE = ReplacingMergeTree()
ORDER BY (project_id, pipeline_id, job_id, id)
;


-- coverage_packages_in
CREATE TABLE IF NOT EXISTS coverage_packages_in AS coverage_packages ENGINE = Null;

-- coverage_packages_mv
CREATE MATERIALIZED VIEW IF NOT EXISTS coverage_packages_mv TO coverage_packages
AS
SELECT * FROM coverage_packages_in
WHERE id NOT IN (
    SELECT id FROM coverage_packages
    WHERE job_id IN (
        SELECT DISTINCT job_id FROM coverage_packages_in
    )
)
;

-- coverage_classes
CREATE TABLE IF NOT EXISTS coverage_classes (
    id String,
    package_id String,
    report_id String,
    job_id Int64,
    pipeline_id Int64,
    project_id Int64,

    package_name String,
    name String,
    filename String,

    line_rate Float32,
    branch_rate Float32,
    complexity Float32,
)
ENGINE = ReplacingMergeTree()
ORDER BY (project_id, pipeline_id, job_id, id)
;


-- coverage_classes_in
CREATE TABLE IF NOT EXISTS coverage_classes_in AS coverage_classes ENGINE = Null;

-- coverage_classes_mv
CREATE MATERIALIZED VIEW IF NOT EXISTS coverage_classes_mv TO coverage_classes
AS
SELECT * FROM coverage_classes_in
WHERE id NOT IN (
    SELECT id FROM coverage_classes
    WHERE job_id IN (
        SELECT DISTINCT job_id FROM coverage_classes_in
    )
)
;

-- coverage_methods
CREATE TABLE IF NOT EXISTS coverage_methods (
    id String,
    class_id String,
    package_id String,
    report_id String,
    job_id Int64,
    pipeline_id Int64,
    project_id Int64,

    package_name String,
    class_name String,
    name String,
    signature String,

    line_rate Float32,
    branch_rate Float32,
    complexity Float32,
)
ENGINE = ReplacingMergeTree()
ORDER BY (project_id, pipeline_id, job_id, id)
;


-- coverage_methods_in
CREATE TABLE IF NOT EXISTS coverage_methods_in AS coverage_methods ENGINE = Null;

-- coverage_methods_mv
CREATE MATERIALIZED VIEW IF NOT EXISTS coverage_methods_mv TO coverage_methods
AS
SELECT * FROM coverage_methods_in
WHERE id NOT IN (
    SELECT id FROM coverage_methods
    WHERE job_id IN (
        SELECT DISTINCT job_id FROM coverage_methods_in
    )
)
;
