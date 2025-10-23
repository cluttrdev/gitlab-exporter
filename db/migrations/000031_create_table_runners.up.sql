-- runners (current state - deduplicated)
CREATE TABLE IF NOT EXISTS runners (
	`id` Int64,
	`short_sha` String,
	`description` String,

	`runner_type` String,
	`tag_list` Array(String),
	`status` String,

	`locked` Bool,
	`paused` Bool,

	`run_protected` Bool,
	`run_untagged` Bool,

	`created_at` Float64,
	`contacted_at` Float64,

	`created_by_id` Int64,
	`created_by_username` String,
	`created_by_name` String,

	`_fetched_at` Float64
)
ENGINE = ReplacingMergeTree(_fetched_at)
ORDER BY (id)
;

-- _runners_raw (event log - all records)
CREATE TABLE IF NOT EXISTS _runners_raw (
	`id` Int64,
	`short_sha` String,
	`description` String,

	`runner_type` String,
	`tag_list` Array(String),
	`status` String,

	`locked` Bool,
	`paused` Bool,

	`run_protected` Bool,
	`run_untagged` Bool,

	`created_at` Float64,
	`contacted_at` Float64,

	`created_by_id` Int64,
	`created_by_username` String,
	`created_by_name` String,

	`_fetched_at` Float64
)
ENGINE = MergeTree()
ORDER BY (id, _fetched_at)
;

-- runners_in
CREATE TABLE IF NOT EXISTS runners_in AS _runners_raw ENGINE = Null;

-- runners_mv
CREATE MATERIALIZED VIEW IF NOT EXISTS runners_mv TO runners AS
SELECT runners_in.* FROM runners_in LEFT OUTER JOIN runners ON runners_in.id = runners.id
WHERE runners_in._fetched_at > runners._fetched_at
;

-- _runners_raw_mv
CREATE MATERIALIZED VIEW IF NOT EXISTS _runners_raw_mv TO _runners_raw AS
SELECT * FROM runners_in
;
