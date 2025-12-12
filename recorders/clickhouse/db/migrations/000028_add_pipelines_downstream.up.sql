-- pipelines
ALTER TABLE pipelines
ADD COLUMN IF NOT EXISTS downstream_pipelines Array(Tuple(id Int64, iid Int64, project_id Int64)) AFTER upstream_pipeline_project_id
;

-- pipelines_in
DROP TABLE IF EXISTS pipelines_in;
CREATE TABLE IF NOT EXISTS pipelines_in AS pipelines ENGINE = Null;
