-- pipelines
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS ref_path String AFTER ref;

DROP TABLE IF EXISTS pipelines_in;
CREATE TABLE IF NOT EXISTS pipelines_in AS pipelines ENGINE = Null;

-- jobs
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS ref_path String AFTER ref;

DROP TABLE IF EXISTS jobs_in;
CREATE TABLE IF NOT EXISTS jobs_in AS jobs ENGINE = Null;
