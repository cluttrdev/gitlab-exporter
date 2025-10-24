-- issues
CREATE TABLE IF NOT EXISTS issues (
    `id` Int64,
    `iid` Int64,
    `project_id` Int64,

    `created_at` Float64,
    `updated_at` Float64,
    `closed_at` Float64,

    `title` String,
    `labels` Array(String),

    `type` String,
    `severity` String,
    `state` String,
)
ENGINE = ReplacingMergeTree(updated_at)
ORDER BY (project_id, id)
;

-- issues_in
CREATE TABLE IF NOT EXISTS issues_in AS issues ENGINE = Null;

-- issues_mv
CREATE MATERIALIZED VIEW IF NOT EXISTS issues_mv TO issues AS
SELECT issues_in.* FROM issues_in LEFT OUTER JOIN issues USING id
WHERE issues_in.updated_at > issues.updated_at
;
