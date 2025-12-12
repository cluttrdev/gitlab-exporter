-- mergerequests
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS assignee_id Int64 AFTER author_id;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS close_user_id Int64 AFTER merge_user_id;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS sha String AFTER close_user_id;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS squash_commit_sha String AFTER sha;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS changes_count Int64 AFTER squash_commit_sha;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS user_notes_count Int64 AFTER changes_count;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS upvotes Int64 AFTER user_notes_count;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS downvotes Int64 AFTER upvotes;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS head_pipeline_id Int64 AFTER downvotes;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS web_url String;

ALTER TABLE mergerequests MODIFY COLUMN IF EXISTS labels Array(String) AFTER close_user_id;
ALTER TABLE mergerequests MODIFY COLUMN IF EXISTS merge_error String AFTER conflicts;
ALTER TABLE mergerequests MODIFY COLUMN IF EXISTS merge_commit_sha String AFTER sha;

ALTER TABLE mergerequests RENAME COLUMN IF EXISTS conflicts TO has_conflicts;
ALTER TABLE mergerequests RENAME COLUMN IF EXISTS base_sha TO diff_ref_base_sha;
ALTER TABLE mergerequests RENAME COLUMN IF EXISTS head_sha TO diff_ref_head_sha;
ALTER TABLE mergerequests RENAME COLUMN IF EXISTS start_sha TO diff_ref_start_sha;

ALTER TABLE mergerequests DROP COLUMN IF EXISTS name;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS additions;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS changes;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS deletions;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS file_count;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS commit_count;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS rebase_commit_sha;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS approvers_id;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS approved;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS mergeable;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS milestone_iid;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS milestone_project_id;

-- mergerequests_in
DROP TABLE IF EXISTS mergerequests_in;
CREATE TABLE IF NOT EXISTS mergerequests_in AS mergerequests ENGINE = Null;

-- mergerequest_notevents
ALTER TABLE mergerequest_noteevents ADD COLUMN IF NOT EXISTS confidential Bool AFTER resolver_id;

ALTER TABLE mergerequest_noteevents DROP COLUMN IF EXISTS resolved_at;

ALTER TABLE mergerequest_noteevents RENAME COLUMN IF EXISTS mergerequest_project_id TO project_id;

-- mergerequest_noteevents_in
DROP TABLE IF EXISTS mergerequest_noteevents_in;
CREATE TABLE IF NOT EXISTS mergerequest_noteevents_in AS mergerequest_noteevents ENGINE = Null;
