-- sections
ALTER TABLE sections ADD COLUMN IF NOT EXISTS job Tuple(id Int64, name String, status String) AFTER name;
ALTER TABLE sections ADD COLUMN IF NOT EXISTS pipeline Tuple(id Int64, project_id Int64, ref String, sha String, status String) AFTER job;

ALTER TABLE sections UPDATE job = tuple(job_id, '', '') WHERE job.id = 0;
ALTER TABLE sections UPDATE pipeline = tuple(pipeline_id, job_id, '', '', '') WHERE pipeline.id = 0;
/* ALTER ... UPDATE queries (mutations) run async, so dropping columns would fail
ALTER TABLE sections DROP COLUMN IF EXISTS job_id
ALTER TABLE sections DROP COLUMN IF EXISTS pipeline_id
ALTER TABLE sections DROP COLUMN IF EXISTS project_id
*/
ALTER TABLE sections MODIFY COLUMN IF EXISTS job REMOVE COMMENT;
ALTER TABLE sections MODIFY COLUMN IF EXISTS pipeline REMOVE COMMENT;

-- sections_in
DROP TABLE IF EXISTS sections_in;
CREATE TABLE IF NOT EXISTS sections_in AS sections ENGINE = Null;
