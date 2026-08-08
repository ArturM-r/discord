-- +migrate Up

-- create singular table 'member' because some repo files INSERT/DELETE from `member`
CREATE TABLE IF NOT EXISTS member (
  id UUID PRIMARY KEY,
  server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  role TEXT NOT NULL,
  CONSTRAINT member_server_user_unique UNIQUE (server_id, user_id)
);

-- view to satisfy queries that reference `members`
CREATE OR REPLACE VIEW members AS SELECT * FROM member;

CREATE INDEX IF NOT EXISTS idx_member_server_id ON member(server_id);
CREATE INDEX IF NOT EXISTS idx_member_user_id ON member(user_id);
