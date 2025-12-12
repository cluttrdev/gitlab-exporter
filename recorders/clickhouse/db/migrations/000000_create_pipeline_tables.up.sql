-- pipelines
CREATE TABLE IF NOT EXISTS pipelines (
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

-- jobs
CREATE TABLE IF NOT EXISTS jobs (
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

-- sections
CREATE TABLE IF NOT EXISTS sections (
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

-- bridges
CREATE TABLE IF NOT EXISTS bridges (
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

