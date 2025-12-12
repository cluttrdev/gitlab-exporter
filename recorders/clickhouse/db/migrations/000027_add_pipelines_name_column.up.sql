-- pipelines
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS name String AFTER project_id;

-- pipelines_in
DROP TABLE IF EXISTS pipelines_in;
CREATE TABLE IF NOT EXISTS pipelines_in AS pipelines ENGINE = Null;
