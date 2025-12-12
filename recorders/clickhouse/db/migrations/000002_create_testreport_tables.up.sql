-- testreports
CREATE TABLE IF NOT EXISTS testreports (
    id String,
    pipeline_id Int64,
    total_time Float64,
    total_count Int64,
    success_count Int64,
    failed_count Int64,
    skipped_count Int64,
    error_count Int64,
)
ENGINE ReplacingMergeTree()
ORDER BY id
;


-- testsuites
CREATE TABLE IF NOT EXISTS testsuites (
    id String,
    testreport_id String,
    pipeline_id Int64,
    name String,
    total_time Float64,
    total_count Int64,
    success_count Int64,
    failed_count Int64,
    skipped_count Int64,
    error_count Int64,
)
ENGINE ReplacingMergeTree()
ORDER BY id
;

-- testcases
CREATE TABLE IF NOT EXISTS testcases (
    id String,
    testsuite_id String,
    testreport_id String,
    pipeline_id Int64,
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
