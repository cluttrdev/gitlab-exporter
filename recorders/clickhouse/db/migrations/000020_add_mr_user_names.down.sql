-- mergerequests
ALTER TABLE mergerequests DROP COLUMN IF EXISTS author_username;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS author_name;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS assignees_username;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS assignees_name;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS reviewers_username;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS reviewers_name;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS approvers_username;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS approvers_name;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS merge_user_username;
ALTER TABLE mergerequests DROP COLUMN IF EXISTS merge_user_name;

-- mergerequests_in
DROP TABLE IF EXISTS mergerequests_in;
CREATE TABLE IF NOT EXISTS mergerequests_in AS mergerequests ENGINE = Null;

-- mergerequest_noteevents
ALTER TABLE mergerequest_noteevents DROP COLUMN IF EXISTS author_username;
ALTER TABLE mergerequest_noteevents DROP COLUMN IF EXISTS author_name;
ALTER TABLE mergerequest_noteevents DROP COLUMN IF EXISTS resolver_username;
ALTER TABLE mergerequest_noteevents DROP COLUMN IF EXISTS resolver_name;

-- mergerequest_noteevents_in
DROP TABLE IF EXISTS mergerequest_noteevents_in;
CREATE TABLE IF NOT EXISTS mergerequest_noteevents_in AS mergerequest_noteevents ENGINE = Null;
