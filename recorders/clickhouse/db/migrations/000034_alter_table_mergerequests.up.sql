--
-- mergerequests
--

-- storage table
ALTER TABLE mergerequests ADD COLUMN IF NOT EXISTS `description` String AFTER `title`;

-- insertion table
DROP TABLE IF EXISTS mergerequests_in;
CREATE TABLE IF NOT EXISTS mergerequests_in AS mergerequests ENGINE = Null;
