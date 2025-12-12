-- testsuites
ALTER TABLE testsuites
DROP COLUMN IF EXISTS properties
;

-- testsuites_in
DROP TABLE IF EXISTS testsuites_in;
CREATE TABLE IF NOT EXISTS testsuites_in AS testsuites ENGINE = Null;

-- testcases
ALTER TABLE testcases
DROP COLUMN IF EXISTS properties
;

-- testcases_in
DROP TABLE IF EXISTS testcases_in;
CREATE TABLE IF NOT EXISTS testcases_in AS testcases ENGINE = Null;
