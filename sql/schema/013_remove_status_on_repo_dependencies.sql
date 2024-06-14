-- +goose Up
ALTER TABLE repos_to_dependencies DROP COLUMN status;