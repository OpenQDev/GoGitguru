-- +goose Up
ALTER TABLE repos_to_dependencies ALTER COLUMN url TYPE VARCHAR(150);
