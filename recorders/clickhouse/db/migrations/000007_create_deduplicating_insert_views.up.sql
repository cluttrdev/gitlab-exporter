-- projects
CREATE TABLE IF NOT EXISTS projects_in AS projects ENGINE = Null;

CREATE MATERIALIZED VIEW IF NOT EXISTS projects_mv
TO projects
AS SELECT projects_in.* FROM projects_in LEFT OUTER JOIN projects ON projects_in.id = projects.id
WHERE projects_in.last_activity_at > projects.last_activity_at
;

-- pipelines
CREATE TABLE IF NOT EXISTS pipelines_in AS pipelines ENGINE = Null;

CREATE MATERIALIZED VIEW IF NOT EXISTS pipelines_mv
TO pipelines
AS SELECT pipelines_in.* FROM pipelines_in LEFT OUTER JOIN pipelines ON pipelines_in.id = pipelines.id
WHERE pipelines_in.updated_at > pipelines.updated_at
;

-- jobs
CREATE TABLE IF NOT EXISTS jobs_in AS jobs ENGINE = Null;

CREATE MATERIALIZED VIEW IF NOT EXISTS jobs_mv
TO jobs
AS SELECT jobs_in.* FROM jobs_in
WHERE id NOT IN ( SELECT id FROM jobs )
;

-- sections
CREATE TABLE IF NOT EXISTS sections_in AS sections ENGINE = Null;

CREATE MATERIALIZED VIEW IF NOT EXISTS sections_mv
TO sections
AS SELECT sections_in.* FROM sections_in
WHERE id NOT IN ( SELECT id FROM sections )
;

-- bridges
CREATE TABLE IF NOT EXISTS bridges_in AS bridges ENGINE = Null;

CREATE MATERIALIZED VIEW IF NOT EXISTS bridges_mv
TO bridges
AS SELECT bridges_in.* FROM bridges_in
WHERE id NOT IN ( SELECT id FROM bridges )
;

-- testreports
CREATE TABLE IF NOT EXISTS testreports_in AS testreports ENGINE = Null;

CREATE MATERIALIZED VIEW IF NOT EXISTS testreports_mv
TO testreports
AS SELECT testreports_in.* FROM testreports_in
WHERE id NOT IN ( SELECT id FROM testreports )
;

-- testsuites
CREATE TABLE IF NOT EXISTS testsuites_in AS testsuites ENGINE = Null;

CREATE MATERIALIZED VIEW IF NOT EXISTS testsuites_mv
TO testsuites
AS SELECT testsuites_in.* FROM testsuites_in
WHERE id NOT IN ( SELECT id FROM testsuites )
;

-- testcases
CREATE TABLE IF NOT EXISTS testcases_in AS testcases ENGINE = Null;

CREATE MATERIALIZED VIEW IF NOT EXISTS testcases_mv
TO testcases
AS SELECT testcases_in.* FROM testcases_in
WHERE id NOT IN ( SELECT id FROM testcases )
;

-- mergerequests
CREATE TABLE IF NOT EXISTS mergerequests_in AS mergerequests ENGINE = Null;

CREATE MATERIALIZED VIEW IF NOT EXISTS mergerequests_mv
TO mergerequests
AS SELECT mergerequests_in.* FROM mergerequests_in LEFT OUTER JOIN mergerequests ON mergerequests_in.id = mergerequests.id
WHERE mergerequests_in.updated_at > mergerequests.updated_at
;

-- metrics
CREATE TABLE IF NOT EXISTS metrics_in AS metrics ENGINE = Null;

CREATE MATERIALIZED VIEW IF NOT EXISTS metrics_mv
TO metrics
AS SELECT metrics_in.* FROM metrics_in
WHERE id NOT IN ( SELECT id FROM metrics )
;

-- traces
CREATE TABLE IF NOT EXISTS traces_in AS traces ENGINE = Null;

CREATE MATERIALIZED VIEW IF NOT EXISTS traces_mv
TO traces
AS SELECT traces_in.* FROM traces_in
WHERE SpanId NOT IN ( SELECT SpanId FROM traces )
;
