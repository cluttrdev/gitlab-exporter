-- testreports
ALTER TABLE testreports DROP COLUMN IF EXISTS project_id;

-- testreports_in
DROP TABLE IF EXISTS testreports_in;
CREATE TABLE IF NOT EXISTS testreports_in AS testreports ENGINE = Null;

-- testsuites
ALTER TABLE testsuites DROP COLUMN IF EXISTS project_id;

-- testsuites_in
DROP TABLE IF EXISTS testsuites_in;
CREATE TABLE IF NOT EXISTS testsuites_in AS testsuites ENGINE = Null;

-- testcases
ALTER TABLE testcases ADD COLUMN IF NOT EXISTS stack_trace String AFTER system_output;
ALTER TABLE testcases ADD COLUMN IF NOT EXISTS recent_failures Tuple(count Int64, base_branch String) AFTER attachment_url;

ALTER TABLE testcases DROP COLUMN IF EXISTS project_id;

-- testcases_in
DROP TABLE IF EXISTS testcases_in;
CREATE TABLE IF NOT EXISTS testcases_in AS testcases ENGINE = Null;
