-- jobs
ALTER TABLE jobs MODIFY COLUMN IF EXISTS id Int64 FIRST;
ALTER TABLE jobs MODIFY COLUMN IF EXISTS coverage Float64 AFTER queued_duration;
ALTER TABLE jobs MODIFY COLUMN IF EXISTS allow_failure Bool AFTER tag_list;

ALTER TABLE jobs ADD COLUMN IF NOT EXISTS pipeline_id Int64 AFTER id;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS project_id Int64 AFTER pipeline_id;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS queued_at Float64 AFTER created_at;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS manual Bool AFTER allow_failure;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS retried Bool AFTER manual;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS retryable Bool AFTER retried;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS kind String;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS downstream_pipeline_id Int64;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS downstream_pipeline_iid Int64;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS downstream_pipeline_project_id Int64;
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS runner_id String;

ALTER TABLE jobs DROP COLUMN IF EXISTS tag;
ALTER TABLE jobs DROP COLUMN IF EXISTS web_url;

ALTER TABLE jobs UPDATE pipeline_id = tupleElement(pipeline, 'id') WHERE pipeline_id = 0;
ALTER TABLE jobs UPDATE project_id = tupleElement(pipeline, 'project_id') WHERE project_id = 0;
/* ALTER ... UPDATE queries (mutations) run async, so dropping column would fail
ALTER TABLE jobs DROP COLUMN IF EXISTS pipeline
*/
ALTER TABLE jobs COMMENT COLUMN IF EXISTS pipeline 'deprecated';

-- jobs_in
DROP TABLE IF EXISTS jobs_in;
CREATE TABLE IF NOT EXISTS jobs_in AS jobs ENGINE = Null;
