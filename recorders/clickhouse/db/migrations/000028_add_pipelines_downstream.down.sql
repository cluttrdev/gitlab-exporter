-- pipelines
ALTER TABLE pipelines DROP COLUMN IF EXISTS downstream_pipelines;

-- pipelines_in
DROP TABLE IF EXISTS pipelines_in;
CREATE TABLE IF NOT EXISTS pipelines_in AS pipelines ENGINE = Null;
