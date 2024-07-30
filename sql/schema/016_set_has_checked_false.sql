-- +goose Up

ALTER TABLE commits ALTER COLUMN has_checked_user SET DEFAULT FALSE;