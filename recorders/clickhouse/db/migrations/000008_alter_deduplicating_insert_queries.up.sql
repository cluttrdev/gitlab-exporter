-- projects

-- pipelines

-- jobs
ALTER TABLE jobs_mv MODIFY QUERY
    SELECT * FROM jobs_in WHERE id NOT IN (
        SELECT id FROM jobs WHERE pipeline.id IN (
            SELECT DISTINCT tupleElement(pipeline, 'id') FROM jobs_in
        )
    )
;

-- sections
ALTER TABLE sections_mv MODIFY QUERY
    SELECT * FROM sections_in WHERE id NOT IN (
        SELECT id FROM sections WHERE job.id IN (
            SELECT DISTINCT tupleElement(job, 'id') FROM sections_in
        )
    )
;

-- bridges
ALTER TABLE bridges_mv MODIFY QUERY
    SELECT * FROM bridges_in WHERE id NOT IN (
        SELECT id FROM bridges WHERE pipeline.id IN (
            SELECT DISTINCT tupleElement(pipeline, 'id') FROM bridges_in
        )
    )
;

-- testreports
ALTER TABLE testreports_mv MODIFY QUERY
    SELECT * FROM testreports_in WHERE id NOT IN (
        SELECT id FROM testreports WHERE pipeline_id IN (
            SELECT DISTINCT pipeline_id FROM testreports_in
        )
    )
;

-- testsuites
ALTER TABLE testsuites_mv MODIFY QUERY
    SELECT * FROM testsuites_in WHERE id NOT IN (
        SELECT id FROM testsuites WHERE testreport_id IN (
            SELECT DISTINCT testreport_id FROM testsuites_in
        )
    )
;

-- testcases
ALTER TABLE testcases_mv MODIFY QUERY
    SELECT * FROM testcases_in WHERE id NOT IN (
        SELECT id FROM testcases WHERE testsuite_id IN (
            SELECT DISTINCT testsuite_id FROM testcases_in
        )
    )
;

-- mergerequests

-- metrics
ALTER TABLE metrics_mv MODIFY QUERY
    SELECT * FROM metrics_in WHERE id NOT IN (
        SELECT id FROM metrics WHERE job_id IN (
            SELECT DISTINCT job_id FROM metrics_in
        )
    )
;

-- traces
ALTER TABLE traces_mv MODIFY QUERY
    SELECT * FROM traces_in WHERE SpanId NOT IN (
        SELECT SpanId FROM traces WHERE TraceId IN (
            SELECT DISTINCT TraceId FROM traces_in
        )
    )
;

