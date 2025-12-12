-- projects
ALTER TABLE projects DROP COLUMN IF EXISTS default_branch;

DROP TABLE IF EXISTS projects_in;
CREATE TABLE IF NOT EXISTS projects_in AS projects ENGINE = Null;
