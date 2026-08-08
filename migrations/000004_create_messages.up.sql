-- +migrate Up

CREATE TABLE IF NOT EXISTS messages (
  id UUID PRIMARY KEY,
  channel_id UUID NOT NULL REFERENCES channel(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  username TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_messages_channel_id ON messages(channel_id);
