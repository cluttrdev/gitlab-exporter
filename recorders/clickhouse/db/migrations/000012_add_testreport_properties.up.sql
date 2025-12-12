-- testsuites
ALTER TABLE testsuites
ADD COLUMN IF NOT EXISTS properties Array(Tuple(name String, value String))
;

ALTER TABLE testsuites_in
ADD COLUMN IF NOT EXISTS properties Array(Tuple(name String, value String))
;

-- testcases
ALTER TABLE testcases
ADD COLUMN IF NOT EXISTS properties Array(Tuple(name String, value String))
;

ALTER TABLE testcases_in
ADD COLUMN IF NOT EXISTS properties Array(Tuple(name String, value String))
;
