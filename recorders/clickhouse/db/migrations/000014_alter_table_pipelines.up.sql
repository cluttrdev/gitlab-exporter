-- pipelines
ALTER TABLE pipelines ALTER COLUMN IF EXISTS yaml_errors TYPE Bool;

ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS failure_reason String AFTER status;
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS warnings Bool AFTER tag;
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS child Bool;
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS upstream_pipeline_id Int64;
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS upstream_pipeline_iid Int64;
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS upstream_pipeline_project_id Int64;
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS merge_request_id Int64;
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS merge_request_iid Int64;
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS merge_request_project_id Int64;
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS user_id Int64;

ALTER TABLE pipelines DROP COLUMN IF EXISTS before_sha;
ALTER TABLE pipelines DROP COLUMN IF EXISTS tag;
ALTER TABLE pipelines DROP COLUMN IF EXISTS web_url;

-- pipelines_in
DROP TABLE IF EXISTS pipelines_in;
CREATE TABLE IF NOT EXISTS pipelines_in AS pipelines ENGINE = Null;
