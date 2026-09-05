-- Bound the memory used by the deduplicating JOIN-pattern materialized views.
--
-- The original definitions performed a `LEFT OUTER JOIN` directly against the
-- full target table. ClickHouse builds the join hash table from the right-hand
-- side, so every insert scanned and materialised the whole target table in RAM,
-- eventually exceeding the server memory limit.
--
-- Each view is rewritten to join against a subquery that is constrained to the
-- keys present in the inserted block and collapsed with `max(<version>)`.
-- Because the target tables are ordered by the filtered key (prefix), the
-- `IN (...)` filter becomes a primary-key lookup that reads only a handful of
-- granules. The `max()` + `GROUP BY` also collapses unmerged ReplacingMergeTree
-- duplicates, preventing a single incoming row from fanning out to multiple
-- output rows.
--
-- Semantics are preserved: `LEFT OUTER JOIN` yields `<version> = 0` for
-- unmatched rows, so brand-new records still pass the `>` comparison.

-- projects
ALTER TABLE projects_mv MODIFY QUERY
SELECT projects_in.* FROM projects_in
LEFT OUTER JOIN (
    SELECT id, max(last_activity_at) AS last_activity_at
    FROM projects
    WHERE id IN (SELECT DISTINCT id FROM projects_in)
    GROUP BY id
) AS cur ON projects_in.id = cur.id
WHERE projects_in.last_activity_at > cur.last_activity_at
;

-- pipelines
ALTER TABLE pipelines_mv MODIFY QUERY
SELECT pipelines_in.* FROM pipelines_in
LEFT OUTER JOIN (
    SELECT id, max(updated_at) AS updated_at
    FROM pipelines
    WHERE id IN (SELECT DISTINCT id FROM pipelines_in)
    GROUP BY id
) AS cur ON pipelines_in.id = cur.id
WHERE pipelines_in.updated_at > cur.updated_at
;

-- mergerequests
ALTER TABLE mergerequests_mv MODIFY QUERY
SELECT mergerequests_in.* FROM mergerequests_in
LEFT OUTER JOIN (
    SELECT id, max(updated_at) AS updated_at
    FROM mergerequests
    WHERE id IN (SELECT DISTINCT id FROM mergerequests_in)
    GROUP BY id
) AS cur ON mergerequests_in.id = cur.id
WHERE mergerequests_in.updated_at > cur.updated_at
;

-- mergerequest_noteevents
ALTER TABLE mergerequest_noteevents_mv MODIFY QUERY
SELECT mergerequest_noteevents_in.* FROM mergerequest_noteevents_in
LEFT OUTER JOIN (
    SELECT id, max(updated_at) AS updated_at
    FROM mergerequest_noteevents
    WHERE id IN (SELECT DISTINCT id FROM mergerequest_noteevents_in)
    GROUP BY id
) AS cur ON mergerequest_noteevents_in.id = cur.id
WHERE mergerequest_noteevents_in.updated_at > cur.updated_at
;

-- issues (ORDER BY (project_id, id): filter both key columns)
ALTER TABLE issues_mv MODIFY QUERY
SELECT issues_in.* FROM issues_in
LEFT OUTER JOIN (
    SELECT id, max(updated_at) AS updated_at
    FROM issues
    WHERE project_id IN (SELECT DISTINCT project_id FROM issues_in)
      AND id IN (SELECT DISTINCT id FROM issues_in)
    GROUP BY id
) AS cur ON issues_in.id = cur.id
WHERE issues_in.updated_at > cur.updated_at
;

-- deployments (ORDER BY (project_id, id): filter both key columns)
ALTER TABLE deployments_mv MODIFY QUERY
SELECT deployments_in.* FROM deployments_in
LEFT OUTER JOIN (
    SELECT id, max(updated_at) AS updated_at
    FROM deployments
    WHERE project_id IN (SELECT DISTINCT project_id FROM deployments_in)
      AND id IN (SELECT DISTINCT id FROM deployments_in)
    GROUP BY id
) AS cur ON deployments_in.id = cur.id
WHERE deployments_in.updated_at > cur.updated_at
;

-- runners
ALTER TABLE runners_mv MODIFY QUERY
SELECT runners_in.* FROM runners_in
LEFT OUTER JOIN (
    SELECT id, max(_fetched_at) AS _fetched_at
    FROM runners
    WHERE id IN (SELECT DISTINCT id FROM runners_in)
    GROUP BY id
) AS cur ON runners_in.id = cur.id
WHERE runners_in._fetched_at > cur._fetched_at
;
