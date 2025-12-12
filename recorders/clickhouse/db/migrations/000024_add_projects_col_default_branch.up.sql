-- projects
ALTER TABLE projects ADD COLUMN IF NOT EXISTS default_branch String;

DROP TABLE IF EXISTS projects_in;
CREATE TABLE IF NOT EXISTS projects_in AS projects ENGINE = Null;
