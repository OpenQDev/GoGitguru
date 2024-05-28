-- +goose Up
ALTER TABLE dependencies_to_users
ADD COLUMN updated_at DATE;

ALTER TABLE repos_to_dependencies
ADD COLUMN commit_hash TEXT;