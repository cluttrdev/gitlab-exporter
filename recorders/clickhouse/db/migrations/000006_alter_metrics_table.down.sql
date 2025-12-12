CREATE TABLE metrics_old (
    name String,
    labels Map(String, String),
    value Float64,
    timestamp Int64,
    job_id Int64,
    job_name String
)
ENGINE = MergeTree()
ORDER BY (job_id, name, timestamp)
;

INSERT INTO metrics_old SELECT
    name,
    labels,
    value,
    timestamp,
    job_id,
    ''
FROM metrics
;

RENAME TABLE metrics TO metrics_new;

RENAME TABLE metrics_old TO metrics;

DROP TABLE IF EXISTS metrics_new;
