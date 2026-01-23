--
-- PREPARE
--

-- Create old table.
CREATE TABLE _testcases_old
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
)
ENGINE = ReplacingMergeTree()
ORDER BY id
;

-- Create old insertion table.
CREATE TABLE _testcases_in_old AS _testcases_old
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
    WHERE pipeline_id IN (
        SELECT DISTINCT pipeline_id
        FROM testcases_in
    )
)
;

--
-- MIGRATE
--

-- Copy data from new table to old table, dropping the new report_created_at column.
-- Note: This may take a while depending on the size of the testcases table.
INSERT INTO _testcases_old (*)
SELECT * EXCEPT(report_created_at)
FROM `testcases`
SETTINGS max_execution_time=0
;

-- Exchange tables.
EXCHANGE TABLES testcases AND _testcases_old
;
EXCHANGE TABLES testcases_in AND _testcases_in_old
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
    WHERE testsuite_id IN (
        SELECT DISTINCT testsuite_id
        FROM testcases_in
    )
)
;

-- Copy data from temporary table to old table, dropping the new report_created_at column.
-- Note: This is to not lose any data that was inserted during the migration.
INSERT INTO testcases (*)
SELECT * EXCEPT(report_created_at)
FROM `_testcases_tmp`
SETTINGS max_execution_time=0
;

--
-- CLEANUP
--

-- Drop temporary table.
DROP TABLE _testcases_tmp
;

-- Drop new (exchanged) tables.
DROP TABLE _testcases_in_old
;
DROP TABLE _testcases_old
;
