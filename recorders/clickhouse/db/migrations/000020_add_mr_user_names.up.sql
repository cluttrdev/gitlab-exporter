-- mergerequests
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS author_username String AFTER author_id;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS author_name String AFTER author_username;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS assignees_username Array(String) AFTER assignees_id;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS assignees_name Array(String) AFTER assignees_username;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS reviewers_username Array(String) AFTER reviewers_id;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS reviewers_name Array(String) AFTER reviewers_username;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS approvers_username Array(String) AFTER approvers_id;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS approvers_name Array(String) AFTER approvers_username;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS merge_user_username String AFTER merge_user_id;
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS merge_user_name String AFTER merge_user_username;

-- mergerequests_in
DROP TABLE IF EXISTS mergerequests_in;
CREATE TABLE IF NOT EXISTS mergerequests_in AS mergerequests ENGINE = Null;

-- mergerequest_noteevents
ALTER TABLE mergerequest_noteevents ADD COLUMN IF NOT EXISTS author_username String AFTER author_id;
ALTER TABLE mergerequest_noteevents ADD COLUMN IF NOT EXISTS author_name String AFTER author_username;
ALTER TABLE mergerequest_noteevents ADD COLUMN IF NOT EXISTS resolver_username String AFTER resolver_id;
ALTER TABLE mergerequest_noteevents ADD COLUMN IF NOT EXISTS resolver_name String AFTER resolver_username;

-- mergerequest_noteevents_in
DROP TABLE IF EXISTS mergerequest_noteevents_in;
CREATE TABLE IF NOT EXISTS mergerequest_noteevents_in AS mergerequest_noteevents ENGINE = Null;
