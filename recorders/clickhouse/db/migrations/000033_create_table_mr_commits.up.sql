--
-- mergerequest_commits
--

-- storage table
CREATE TABLE IF NOT EXISTS mergerequest_commits (
    `id` String,
    `mergerequest_id` Int64,
    `mergerequest_iid` Int64,
    `project_id` Int64,

    `sha` String,

    `title` String,
    `message` String,
    `trailers` Array(Tuple(key String, value String)),

    `author_id` Int64,
    `author_username` String,

    `authored_date` DateTime,
    `committed_date` DateTime,

    `author_name` String,
    `author_email` String,
    `committer_name` String,
    `committer_email` String
)
ENGINE = ReplacingMergeTree()
PRIMARY KEY (project_id, mergerequest_iid)
ORDER BY (project_id, mergerequest_iid, authored_date, id)
;

-- insertion table
CREATE TABLE IF NOT EXISTS mergerequest_commits_in AS mergerequest_commits
ENGINE = Null
;

-- deduplication view
CREATE MATERIALIZED VIEW IF NOT EXISTS mergerequest_commits_mv
TO mergerequest_commits AS
SELECT mergerequest_commits_in.* FROM mergerequest_commits_in
WHERE id NOT IN (
    SELECT id FROM mergerequest_commits
    WHERE project_id IN (SELECT DISTINCT project_id FROM mergerequest_commits_in)
      AND mergerequest_iid IN (SELECT DISTINCT mergerequest_iid FROM mergerequest_commits_in)
)
;

--
-- mergerequests
--

-- storage table
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS `commit_shas` Array(String);

-- insertion table
DROP TABLE IF EXISTS mergerequests_in;
CREATE TABLE IF NOT EXISTS mergerequests_in AS mergerequests ENGINE = Null;
