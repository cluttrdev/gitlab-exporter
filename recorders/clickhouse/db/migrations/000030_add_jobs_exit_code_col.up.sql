-- jobs
ALTER TABLE jobs ADD COLUMN IF NOT EXISTS exit_code Int64 AFTER failure_reason;

-- jobs_in
DROP TABLE IF EXISTS jobs_in;
CREATE TABLE IF NOT EXISTS jobs_in AS jobs ENGINE = Null;
