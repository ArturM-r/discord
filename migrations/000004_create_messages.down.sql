-- +migrate Down
DROP INDEX IF EXISTS idx_messages_channel_id;
DROP TABLE IF EXISTS messages;
