-- pipelines
ALTER TABLE pipelines ALTER COLUMN IF EXISTS yaml_errors TYPE String;

ALTER TABLE pipelines DROP COLUMN IF EXISTS failure_reason;
ALTER TABLE pipelines DROP COLUMN IF EXISTS warnings;
ALTER TABLE pipelines DROP COLUMN IF EXISTS child;
ALTER TABLE pipelines DROP COLUMN IF EXISTS upstream_pipeline_id;
ALTER TABLE pipelines DROP COLUMN IF EXISTS upstream_pipeline_iid;
ALTER TABLE pipelines DROP COLUMN IF EXISTS upstream_pipeline_project_id;
ALTER TABLE pipelines DROP COLUMN IF EXISTS merge_request_id;
ALTER TABLE pipelines DROP COLUMN IF EXISTS merge_request_iid;
ALTER TABLE pipelines DROP COLUMN IF EXISTS merge_request_project_id;
ALTER TABLE pipelines DROP COLUMN IF EXISTS user_id;

ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS before_sha String AFTER sha;
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS tag Bool AFTER before_sha;
ALTER TABLE pipelines ADD COLUMN IF NOT EXISTS web_url String;

-- pipelines_in
DROP TABLE IF EXISTS pipelines_in;
CREATE TABLE IF NOT EXISTS pipelines_in AS pipelines ENGINE = Null;
