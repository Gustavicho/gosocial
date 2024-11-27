CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id BIGSERIAL PRIMARY KEY,
  username VARCHAR(255) NOT NULL,
  email citext NOT NULL UNIQUE,
  password bytea NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now()
);
