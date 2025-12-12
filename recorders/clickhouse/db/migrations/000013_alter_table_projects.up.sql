-- projects
ALTER TABLE projects DROP COLUMN IF EXISTS owner_id;
ALTER TABLE projects DROP COLUMN IF EXISTS creator_id;
ALTER TABLE projects DROP COLUMN IF EXISTS default_branch;
ALTER TABLE projects DROP COLUMN IF EXISTS empty_repo;
ALTER TABLE projects DROP COLUMN IF EXISTS web_url;

ALTER TABLE projects RENAME COLUMN IF EXISTS name_with_namespace TO full_name;
ALTER TABLE projects RENAME COLUMN IF EXISTS path_with_namespace TO full_path;

ALTER TABLE projects ADD COLUMN IF NOT EXISTS container_registry_size Int64 AFTER job_artifacts_size;

-- projects_in
DROP TABLE IF EXISTS projects_in;
CREATE TABLE IF NOT EXISTS projects_in AS projects ENGINE = Null;
