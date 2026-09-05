-- Restore the original unbounded JOIN-pattern materialized view definitions.

-- projects (from 000007)
ALTER TABLE projects_mv MODIFY QUERY
SELECT projects_in.* FROM projects_in LEFT OUTER JOIN projects ON projects_in.id = projects.id
WHERE projects_in.last_activity_at > projects.last_activity_at
;

-- pipelines (from 000007)
ALTER TABLE pipelines_mv MODIFY QUERY
SELECT pipelines_in.* FROM pipelines_in LEFT OUTER JOIN pipelines ON pipelines_in.id = pipelines.id
WHERE pipelines_in.updated_at > pipelines.updated_at
;

-- mergerequests (from 000007)
ALTER TABLE mergerequests_mv MODIFY QUERY
SELECT mergerequests_in.* FROM mergerequests_in LEFT OUTER JOIN mergerequests ON mergerequests_in.id = mergerequests.id
WHERE mergerequests_in.updated_at > mergerequests.updated_at
;

-- mergerequest_noteevents (from 000009)
ALTER TABLE mergerequest_noteevents_mv MODIFY QUERY
SELECT mergerequest_noteevents_in.* FROM mergerequest_noteevents_in LEFT OUTER JOIN mergerequest_noteevents USING (id)
WHERE mergerequest_noteevents_in.updated_at > mergerequest_noteevents.updated_at
;

-- issues (from 000029)
ALTER TABLE issues_mv MODIFY QUERY
SELECT issues_in.* FROM issues_in LEFT OUTER JOIN issues USING id
WHERE issues_in.updated_at > issues.updated_at
;

-- deployments (from 000021)
ALTER TABLE deployments_mv MODIFY QUERY
SELECT deployments_in.* FROM deployments_in LEFT OUTER JOIN deployments ON deployments_in.id = deployments.id
WHERE deployments_in.updated_at > deployments.updated_at
;

-- runners (from 000031)
ALTER TABLE runners_mv MODIFY QUERY
SELECT runners_in.* FROM runners_in LEFT OUTER JOIN runners ON runners_in.id = runners.id
WHERE runners_in._fetched_at > runners._fetched_at
;
