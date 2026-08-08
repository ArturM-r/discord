-- +migrate Up

-- create singular table 'channel' because some repo files INSERT/DELETE from `channel` (singular)
CREATE TABLE IF NOT EXISTS channel (
  id UUID PRIMARY KEY,
  server_id UUID NOT NULL REFERENCES servers(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- provide a view named 'channels' so code that SELECTs/JOINs against `channels` continues to work
CREATE OR REPLACE VIEW channels AS SELECT * FROM channel;
