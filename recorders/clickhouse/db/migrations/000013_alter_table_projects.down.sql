-- projects
ALTER TABLE projects ADD COLUMN IF NOT EXISTS owner_id Int64 AFTER namespace_id;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS creator_id Int64 AFTER owner_id;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS default_branch Int64 AFTER topics;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS empty_repo Bool AFTER default_branch;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS web_url Int64 AFTER open_issues_count;

ALTER TABLE projects RENAME COLUMN IF EXISTS full_name TO name_with_namespace;
ALTER TABLE projects RENAME COLUMN IF EXISTS full_path TO path_with_namespace;

ALTER TABLE projects DROP COLUMN IF EXISTS container_registry_size;

-- projects_in
DROP TABLE IF EXISTS projects_in;
CREATE TABLE IF NOT EXISTS projects_in AS projects ENGINE = Null;
