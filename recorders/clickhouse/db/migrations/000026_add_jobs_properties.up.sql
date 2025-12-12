-- jobs
ALTER TABLE jobs
ADD COLUMN IF NOT EXISTS properties Array(Tuple(name String, value String)) AFTER tag_list
;

-- jobs_in
DROP TABLE IF EXISTS jobs_in;
CREATE TABLE IF NOT EXISTS jobs_in AS jobs ENGINE = Null;
