-- testreports
ALTER TABLE testreports ADD COLUMN IF NOT EXISTS project_id Int64 AFTER pipeline_id;

-- testreports_in
DROP TABLE IF EXISTS testreports_in;
CREATE TABLE IF NOT EXISTS testreports_in AS testreports ENGINE = Null;

-- testsuites
ALTER TABLE testsuites ADD COLUMN IF NOT EXISTS project_id Int64 AFTER pipeline_id;

-- testsuites_in
DROP TABLE IF EXISTS testsuites_in;
CREATE TABLE IF NOT EXISTS testsuites_in AS testsuites ENGINE = Null;

-- testcases
ALTER TABLE testcases ADD COLUMN IF NOT EXISTS project_id Int64 AFTER pipeline_id;

ALTER TABLE testcases DROP COLUMN IF EXISTS stack_trace;
ALTER TABLE testcases DROP COLUMN IF EXISTS recent_failures;

-- testcases_in
DROP TABLE IF EXISTS testcases_in;
CREATE TABLE IF NOT EXISTS testcases_in AS testcases ENGINE = Null;
