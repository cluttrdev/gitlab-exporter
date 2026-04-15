--
-- mergerequest_commits
--

-- storage table
DROP TABLE IF EXISTS mergerequest_commits;

-- insertion table
DROP TABLE IF EXISTS mergerequest_commits_in;

-- deduplication view
DROP VIEW IF EXISTS mergerequest_commits_m;

--
-- mergerequests
--

-- storage table
ALTER TABLE mergerequests DROP COLUMN IF EXISTS `commits_id`;

-- insertion table
DROP TABLE IF EXISTS mergerequests_in;
CREATE TABLE IF NOT EXISTS mergerequests_in AS mergerequests ENGINE = Null;
