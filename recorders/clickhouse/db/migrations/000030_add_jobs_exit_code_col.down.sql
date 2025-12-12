-- jobs
ALTER TABLE jobs DROP COLUMN IF EXISTS exit_code;

-- jobs_in
DROP TABLE IF EXISTS jobs_in;
CREATE TABLE IF NOT EXISTS jobs_in AS jobs ENGINE = Null;
