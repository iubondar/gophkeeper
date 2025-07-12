-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  login TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  salt TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS records (
  id UUID PRIMARY KEY,
  user_id UUID REFERENCES users(id),
  label TEXT NOT NULL UNIQUE,
  type TEXT NOT NULL,
  metadata TEXT,
  encrypted_data BYTEA,
  file_key TEXT,
  version INTEGER NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ,
  updated_at TIMESTAMPTZ
);

CREATE INDEX idx_records_user_label ON records(user_id, label);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP INDEX IF EXISTS idx_records_user_label;

DROP TABLE IF EXISTS records;

DROP TABLE IF EXISTS users;