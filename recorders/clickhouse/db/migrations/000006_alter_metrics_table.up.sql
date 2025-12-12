CREATE TABLE metrics_new (
    id String,
    iid Int64,
    job_id Int64,
    name String,
    labels Map(String, String),
    value Float64,
    timestamp Int64,
)
ENGINE = ReplacingMergeTree()
ORDER BY (job_id, iid)
;

INSERT INTO metrics_new SELECT
    concatWithSeparator('-', job_id, hex(sipHash64(name, labels, value, timestamp))),
    sipHash64(name, labels, value, timestamp),
    job_id,
    name,
    labels,
    value,
    timestamp
FROM metrics
;

RENAME TABLE metrics TO metrics_old;

RENAME TABLE metrics_new TO metrics;

DROP TABLE IF EXISTS metrics_old;
