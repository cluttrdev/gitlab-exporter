-- testreports
ALTER TABLE testreports DROP COLUMN IF EXISTS job_id;

DROP TABLE IF EXISTS testreports_in;
CREATE TABLE IF NOT EXISTS testreports_in AS testreports ENGINE = Null;

-- testsuites
ALTER TABLE testsuites DROP COLUMN IF EXISTS job_id;

DROP TABLE IF EXISTS testsuites_in;
CREATE TABLE IF NOT EXISTS testsuites_in AS testsuites ENGINE = Null;

-- testcases
ALTER TABLE testcases DROP COLUMN IF EXISTS job_id;

DROP TABLE IF EXISTS testcases_in;
CREATE TABLE IF NOT EXISTS testcases_in AS testreports ENGINE = Null;
