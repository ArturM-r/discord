-- +migrate Down
DROP INDEX IF EXISTS idx_member_user_id;
DROP INDEX IF EXISTS idx_member_server_id;
DROP VIEW IF EXISTS members;
DROP TABLE IF EXISTS member;
