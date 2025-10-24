-- deployments
CREATE TABLE IF NOT EXISTS deployments (
    id Int64,
    iid Int64,
    job_id Int64,
    pipeline_id Int64,
    project_id Int64,

    environment_id Int64,
    environment_name String,
    environment_tier String,

    triggerer_id Int64,
    triggerer_username String,
    triggerer_name String,

    created_at Float64,
    finished_at Float64,
    updated_at Float64,

    status String,
    ref String,
    sha String,
)
ENGINE ReplacingMergeTree(updated_at)
ORDER BY (project_id, id)
;

-- deployments_in
CREATE TABLE IF NOT EXISTS deployments_in AS deployments ENGINE = Null;

-- deployments_mv
CREATE MATERIALIZED VIEW IF NOT EXISTS deployments_mv TO deployments AS
    SELECT deployments_in.* FROM deployments_in LEFT OUTER JOIN deployments ON deployments_in.id = deployments.id
    WHERE deployments_in.updated_at > deployments.updated_at
;
