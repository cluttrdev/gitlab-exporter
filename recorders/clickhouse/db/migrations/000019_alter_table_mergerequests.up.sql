-- mergerequests
ALTER TABLE mergerequests RENAME COLUMN IF EXISTS has_conflicts TO conflicts;
ALTER TABLE mergerequests RENAME COLUMN IF EXISTS diff_ref_base_sha TO base_sha;
ALTER TABLE mergerequests RENAME COLUMN IF EXISTS diff_ref_head_sha TO head_sha;
ALTER TABLE mergerequests RENAME COLUMN IF EXISTS diff_ref_start_sha TO start_sha;

ALTER TABLE mergerequests MODIFY COLUMN IF EXISTS labels Array(String) AFTER title;
ALTER TABLE mergerequests MODIFY COLUMN IF EXISTS merge_error String AFTER merge_status;
ALTER TABLE mergerequests MODIFY COLUMN IF EXISTS merge_commit_sha String AFTER start_sha;

ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS name String AFTER title;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS additions Int64 AFTER target_branch;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS changes Int64 AFTER additions;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS deletions Int64 AFTER changes;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS file_count Int64 AFTER deletions;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS commit_count Int64 AFTER file_count;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS rebase_commit_sha String AFTER merge_commit_sha;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS approvers_id Array(Int64) AFTER reviewers_id;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS approved Bool AFTER conflicts;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS mergeable Bool AFTER approved;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS milestone_iid Int64 AFTER milestone_id;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS milestone_project_id Int64 AFTER milestone_iid;

ALTER TABLE mergerequests DROP COLUMN IF EXISTS assignee_id;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS close_user_id;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS sha;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS squash_commit_sha;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS changes_count;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS user_notes_count;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS upvotes;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS downvotes;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS head_pipeline_id;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS web_url;

-- mergerequests_in
DROP TABLE IF EXISTS mergerequests_in;
CREATE TABLE IF NOT EXISTS mergerequests_in AS mergerequests ENGINE = Null;

-- mergerequest_noteevents
ALTER TABLE mergerequest_noteevents RENAME COLUMN IF EXISTS project_id TO mergerequest_project_id;

ALTER TABLE mergerequest_noteevents ADD COLUMN IF NOT EXISTS resolved_at Float64 AFTER updated_at;

ALTER TABLE mergerequest_noteevents DROP COLUMN IF EXISTS confidential;

-- mergerequest_noteevents_in
DROP TABLE IF EXISTS mergerequest_noteevents_in;
CREATE TABLE IF NOT EXISTS mergerequest_noteevents_in AS mergerequest_noteevents ENGINE = Null;
