-- metrics
ALTER TABLE metrics ADD COLUMN IF NOT EXISTS pipeline_id Int64 AFTER job_id;
ALTER TABLE metrics ADD COLUMN IF NOT EXISTS project_id Int64 AFTER pipeline_id;

-- metrics_in
DROP TABLE IF EXISTS metrics_in;
CREATE TABLE IF NOT EXISTS metrics_in AS metrics ENGINE = Null;
