-- testreports
ALTER TABLE testreports ADD COLUMN IF NOT EXISTS job_id Int64 AFTER id;

DROP TABLE IF EXISTS testreports_in;
CREATE TABLE IF NOT EXISTS testreports_in AS testreports ENGINE = Null;

-- testsuites
ALTER TABLE testsuites ADD COLUMN IF NOT EXISTS job_id Int64 AFTER testreport_id;

DROP TABLE IF EXISTS testsuites_in;
CREATE TABLE IF NOT EXISTS testsuites_in AS testsuites ENGINE = Null;

-- testcases
ALTER TABLE testcases ADD COLUMN IF NOT EXISTS job_id Int64 AFTER testreport_id;

DROP TABLE IF EXISTS testcases_in;
CREATE TABLE IF NOT EXISTS testcases_in AS testcases ENGINE = Null;
