-- projects


-- pipelines

-- jobs
ALTER TABLE jobs_mv MODIFY QUERY
    SELECT * FROM jobs_in WHERE id NOT IN ( SELECT id FROM jobs )
;

-- sections
ALTER TABLE sections_mv MODIFY QUERY
    SELECT * FROM sections_in WHERE id NOT IN ( SELECT id FROM sections )
;

-- bridges
ALTER TABLE bridges_mv MODIFY QUERY
    SELECT * FROM bridges_in WHERE id NOT IN ( SELECT id FROM bridges )
;

-- testreports
ALTER TABLE testreports_mv MODIFY QUERY
    SELECT * FROM testreports_in WHERE id NOT IN ( SELECT id FROM testreports )
;

-- testsuites
ALTER TABLE testsuites_mv MODIFY QUERY
    SELECT * FROM testsuites_in WHERE id NOT IN ( SELECT id FROM testsuites )
;

-- testcases
ALTER TABLE testcases_mv MODIFY QUERY
    SELECT * FROM testcases_in WHERE id NOT IN ( SELECT id FROM testcases )
;

-- mergerequests

-- metrics
ALTER TABLE metrics_mv MODIFY QUERY
    SELECT * FROM metrics_in WHERE id NOT IN ( SELECT id FROM metrics )
;

-- traces
ALTER TABLE traces_mv MODIFY QUERY
    SELECT * FROM traces_in WHERE SpanId NOT IN ( SELECT SpanId FROM traces )
;


