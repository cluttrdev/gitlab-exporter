-- pipelines
ALTER TABLE pipelines DROP COLUMN IF EXISTS name;

-- pipelines_in
DROP TABLE IF EXISTS pipelines_in;
CREATE TABLE IF NOT EXISTS pipelines_in AS pipelines ENGINE = Null;
