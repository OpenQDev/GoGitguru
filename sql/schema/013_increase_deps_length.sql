-- +goose Up 

ALTER TABLE dependencies
    ALTER COLUMN dependency_name TYPE VARCHAR(255);