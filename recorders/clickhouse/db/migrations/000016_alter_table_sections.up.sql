-- sections
ALTER TABLE sections ADD COLUMN IF NOT EXISTS job_id Int64 AFTER id;
ALTER TABLE sections ADD COLUMN IF NOT EXISTS pipeline_id Int64 AFTER job_id;
ALTER TABLE sections ADD COLUMN IF NOT EXISTS project_id Int64 AFTER pipeline_id;

ALTER TABLE sections UPDATE job_id = tupleElement(job, 'id') WHERE job_id = 0;
ALTER TABLE sections UPDATE pipeline_id = tupleElement(pipeline, 'id') WHERE pipeline_id = 0;
ALTER TABLE sections UPDATE project_id = tupleElement(pipeline, 'project_id') WHERE project_id = 0;
/* ALTER ... UPDATE queries (mutations) run async, so dropping columns would fail
ALTER TABLE sections DROP COLUMN IF EXISTS job
ALTER TABLE sections DROP COLUMN IF EXISTS pipeline
*/
ALTER TABLE sections COMMENT COLUMN IF EXISTS job 'deprecated';
ALTER TABLE sections COMMENT COLUMN IF EXISTS pipeline 'deprecated';

-- sections_in
DROP TABLE IF EXISTS sections_in;
CREATE TABLE IF NOT EXISTS sections_in AS sections ENGINE = Null;
