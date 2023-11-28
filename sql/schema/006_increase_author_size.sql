-- +goose Up

ALTER TABLE commits
ALTER COLUMN author
TYPE VARCHAR(255);