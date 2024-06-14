-- +goose Up
ALTER TABLE user_to_dependencies ADD COLUMN updated_at BIGINT;
