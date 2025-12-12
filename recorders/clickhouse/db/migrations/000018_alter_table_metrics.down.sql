-- metrics
ALTER TABLE metrics DROP COLUMN IF EXISTS pipeline_id;
ALTER TABLE metrics DROP COLUMN IF EXISTS project_id;

-- metrics_in
DROP TABLE IF EXISTS metrics_in;
CREATE TABLE IF NOT EXISTS metrics_in AS metrics ENGINE = Null;
