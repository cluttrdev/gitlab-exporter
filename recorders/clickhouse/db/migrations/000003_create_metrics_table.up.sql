-- metrics
CREATE TABLE IF NOT EXISTS metrics (
    name String,
    labels Map(String, String),
    value Float64,
    timestamp Int64,
    job_id Int64,
    job_name String
)
ENGINE MergeTree()
ORDER BY (job_id, name, timestamp)
;
