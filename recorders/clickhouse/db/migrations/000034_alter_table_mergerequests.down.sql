--
-- mergerequests
--

-- storage table
ALTER TABLE mergerequests DROP COLUMN IF EXISTS `description`;

-- insertion table
DROP TABLE IF EXISTS mergerequests_in;
CREATE TABLE IF NOT EXISTS mergerequests_in AS mergerequests ENGINE = Null;
