-- mergerequest_noteevents

-- storage table
CREATE TABLE IF NOT EXISTS mergerequest_noteevents (
    id Int64,
    mergerequest_id Int64,
    mergerequest_iid Int64,
    project_id Int64,
    created_at Float64,
    updated_at Float64,
    type String,
    system Bool,
    author_id Int64,
    resolvable Bool,
    resolved Bool,
    resolver_id Int64,
    confidential Bool,
    internal Bool
)
ENGINE = ReplacingMergeTree(updated_at)
ORDER BY id
;

-- insertion table
CREATE TABLE IF NOT EXISTS mergerequest_noteevents_in AS mergerequest_noteevents ENGINE = Null;

-- deduplication view
CREATE MATERIALIZED VIEW IF NOT EXISTS mergerequest_noteevents_mv
TO mergerequest_noteevents
AS SELECT mergerequest_noteevents_in.* FROM mergerequest_noteevents_in LEFT OUTER JOIN mergerequest_noteevents USING (id)
    WHERE mergerequest_noteevents_in.updated_at > mergerequest_noteevents.updated_at
;
