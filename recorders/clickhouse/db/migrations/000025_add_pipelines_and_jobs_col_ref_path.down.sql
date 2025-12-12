-- pipelines
ALTER TABLE pipelines DROP COLUMN IF EXISTS ref_path;

-- jobs
ALTER TABLE jobs DROP COLUMN IF EXISTS ref_path;
