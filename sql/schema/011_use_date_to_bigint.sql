-- +goose Up
ALTER TABLE repos_to_dependencies ALTER COLUMN first_use_data TYPE BIGINT USING CAST(first_use_data AS TEXT)::BIGINT;
ALTER TABLE repos_to_dependencies ALTER COLUMN last_use_data TYPE BIGINT USING CAST(last_use_data AS TEXT)::BIGINT;
ALTER TABLE dependencies_to_users ALTER COLUMN first_use_data TYPE BIGINT USING  CAST(first_use_data AS TEXT)::BIGINT;
ALTER TABLE dependencies_to_users ALTER COLUMN last_use_data TYPE BIGINT USING  CAST(first_use_data AS TEXT)::BIGINT;