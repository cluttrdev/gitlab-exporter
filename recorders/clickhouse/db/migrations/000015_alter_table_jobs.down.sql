-- jobs
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS tag Bool AFTER failure_reason;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS web_url String AFTER tag;

ALTER TABLE jobs UPDATE pipeline = tuple(pipeline_id, project_id, '', '', '') WHERE pipeline.id = 0;
/* ALTER ... UPDATE queries (mutations) run async, so dropping columns would fail
ALTER TABLE jobs DROP COLUMN IF EXISTS pipeline_id
ALTER TABLE jobs DROP COLUMN IF EXISTS project_id
*/
ALTER TABLE jobs MODIFY COLUMN IF EXISTS pipeline REMOVE COMMENT;

ALTER TABLE jobs DROP COLUMN IF EXISTS queued_at;
ALTER TABLE jobs DROP COLUMN IF EXISTS manual;
ALTER TABLE jobs DROP COLUMN IF EXISTS retried;
ALTER TABLE jobs DROP COLUMN IF EXISTS retryable;
ALTER TABLE jobs DROP COLUMN IF EXISTS kind;
ALTER TABLE jobs DROP COLUMN IF EXISTS downstream_pipeline_id;
ALTER TABLE jobs DROP COLUMN IF EXISTS downstream_pipeline_iid;
ALTER TABLE jobs DROP COLUMN IF EXISTS downstream_pipeline_project_id;
ALTER TABLE jobs DROP COLUMN IF EXISTS runner_id;

ALTER TABLE jobs MODIFY COLUMN IF EXISTS id Int64 AFTER tag_list;
ALTER TABLE jobs MODIFY COLUMN IF EXISTS coverage Float64 FIRST;
ALTER TABLE jobs MODIFY COLUMN IF EXISTS allow_failure Bool AFTER coverage;

-- jobs_in
DROP TABLE IF EXISTS jobs_in;
CREATE TABLE IF NOT EXISTS jobs_in AS jobs ENGINE = Null;

