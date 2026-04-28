DROP INDEX IF EXISTS cloud.idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS cloud.idx_refresh_tokens_user_id;
DROP TABLE IF EXISTS cloud.refresh_tokens;
DROP INDEX IF EXISTS cloud.idx_files_user_status;
DROP TABLE IF EXISTS cloud.files;
DROP TABLE IF EXISTS cloud.users;
DROP TYPE IF EXISTS cloud.file_status;
DROP SCHEMA IF EXISTS cloud;
