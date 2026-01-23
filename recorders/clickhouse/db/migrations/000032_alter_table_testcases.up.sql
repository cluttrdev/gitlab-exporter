--
-- PREPARE
--

-- Create new table.
CREATE TABLE _testcases_new
(
    `id` String,
    `testsuite_id` String,
    `testreport_id` String,
    `job_id` Int64,
    `pipeline_id` Int64,
    `project_id` Int64,
    `status` LowCardinality(String),
    `name` String,
    `classname` String,
    `file` String,
    `execution_time` Float64,
    `system_output` String,
    `attachment_url` String,
    `properties` Array(Tuple(
        name String,
        value String)),

    `report_created_at` UInt32 DEFAULT toUInt32(0),
)
ENGINE = ReplacingMergeTree(report_created_at)
PARTITION BY toYYYYMMDD(toDateTime(report_created_at))
ORDER BY (project_id, pipeline_id, job_id, report_created_at, id)
PRIMARY KEY (project_id, pipeline_id, job_id)
;

-- Create new insertion table.
CREATE TABLE _testcases_in_new AS _testcases_new
ENGINE = `Null`
;

-- Create temporary table to receive data inserted during migration.
CREATE TABLE _testcases_tmp AS testcases
;

-- Recreate materialized view to insert into temporary table.
DROP TABLE IF EXISTS testcases_mv
;
CREATE MATERIALIZED VIEW testcases_mv TO _testcases_tmp
AS SELECT *
FROM testcases_in
WHERE id NOT IN (
    SELECT id
    FROM testcases
    WHERE testsuite_id IN (
        SELECT DISTINCT testsuite_id
        FROM testcases_in
    )
)
;

--
-- MIGRATE
--

-- Copy data from old table to new table, adding values for the new report_created_at column.
-- Note: This may take a while depending on the size of the testcases table.
INSERT INTO _testcases_new (*)
SELECT
  testcases.*,
  CASE
    WHEN jps.job_finished_at > 0 THEN toDateTime(toUInt32(job_finished_at))
    WHEN jps.pipeline_finished_at > 0 THEN toDateTime(toUInt32(pipeline_finished_at))
    ELSE toDateTime(toUInt32(0))
  END AS report_created_at
FROM testcases AS testcases
  LEFT OUTER JOIN (
    SELECT
      jobs.id AS job_id,
      jobs.finished_at AS job_finished_at,
      pipelines.finished_at AS pipeline_finished_at
    FROM jobs AS jobs
      INNER JOIN pipelines AS pipelines ON jobs.pipeline_id = pipelines.id
  ) AS jps USING job_id
SETTINGS max_execution_time=0, max_partitions_per_insert_block=0
;

-- Exchange tables.
EXCHANGE TABLES testcases AND _testcases_new
;
EXCHANGE TABLES testcases_in AND _testcases_in_new
;

-- Recreate materialized view to insert into new table.
DROP TABLE IF EXISTS testcases_mv
;
CREATE MATERIALIZED VIEW testcases_mv TO testcases
AS SELECT *
FROM testcases_in
WHERE id NOT IN (
    SELECT id
    FROM testcases
    WHERE pipeline_id IN (
        SELECT DISTINCT pipeline_id
        FROM testcases_in
    )
)
;

-- Copy data from temporary table to new table, adding values for the new report_created_at column.
-- Note: This is to not lose any data that was inserted during the migration.
INSERT INTO testcases (*)
SELECT
  testcases_tmp.*,
  CASE
    WHEN jps.job_finished_at > 0 THEN toDateTime(toUInt32(job_finished_at))
    WHEN jps.pipeline_finished_at > 0 THEN toDateTime(toUInt32(pipeline_finished_at))
    ELSE toDateTime(toUInt32(0))
  END AS report_created_at
FROM _testcases_tmp AS testcases_tmp
  LEFT OUTER JOIN (
    SELECT
      jobs.id AS job_id,
      jobs.finished_at AS job_finished_at,
      pipelines.finished_at AS pipeline_finished_at
    FROM jobs AS jobs
      INNER JOIN pipelines AS pipelines ON jobs.pipeline_id = pipelines.id
  ) AS jps USING job_id
SETTINGS max_execution_time=0, max_partitions_per_insert_block=0
;

--
-- CLEANUP
--

-- Drop temporary table.
DROP TABLE _testcases_tmp
;

-- Drop old (exchanged) tables.
DROP TABLE _testcases_in_new
;
DROP TABLE _testcases_new
;
