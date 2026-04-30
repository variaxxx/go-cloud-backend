CREATE SCHEMA cloud;

CREATE TYPE cloud.file_status AS ENUM ('uploaded', 'processed', 'failed');

CREATE TABLE cloud.users (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  username VARCHAR(100) NOT NULL UNIQUE CHECK (char_length(username) BETWEEN 3 AND 100),
  password_hash TEXT NOT NULL
);

CREATE TABLE cloud.folders (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  name VARCHAR(255) NOT NULL,
  user_id BIGINT NOT NULL REFERENCES cloud.users(id) ON DELETE CASCADE,
  parent_id BIGINT REFERENCES cloud.folders(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX uq_folders_user_parent_name
ON cloud.folders (user_id, COALESCE(parent_id, 0), name);

CREATE TABLE cloud.files (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ,

  filename TEXT NOT NULL,
  mimetype VARCHAR(255),
  status cloud.file_status NOT NULL DEFAULT 'uploaded',
  storage_path TEXT NOT NULL,
  size_bytes BIGINT NOT NULL CHECK (size_bytes >= 0),

  user_id BIGINT NOT NULL REFERENCES cloud.users(id) ON DELETE CASCADE,
  folder_id BIGINT REFERENCES cloud.folders(id) ON DELETE CASCADE
);

CREATE INDEX idx_files_user_status ON cloud.files(user_id, status);
CREATE UNIQUE INDEX uq_files_user_folder_filename
ON cloud.files (user_id, COALESCE(folder_id, 0), filename)
WHERE deleted_at IS NULL;

CREATE TABLE cloud.refresh_tokens (
  id BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMPTZ NOT NULL,
  revoked_at TIMESTAMPTZ,

  token_hash TEXT NOT NULL UNIQUE,

  user_id BIGINT NOT NULL REFERENCES cloud.users(id) ON DELETE CASCADE,
  replaced_by_id BIGINT REFERENCES cloud.refresh_tokens(id) ON DELETE SET NULL
);

CREATE INDEX idx_refresh_tokens_user_id ON cloud.refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON cloud.refresh_tokens(expires_at);
