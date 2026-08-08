-- +migrate Down
-- drop view before table to avoid dependency errors
DROP VIEW IF EXISTS channels;
DROP TABLE IF EXISTS channel;
