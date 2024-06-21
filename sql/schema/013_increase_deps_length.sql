-- +goose Up 

ALTER TABLE dependencies
    ALTER COLUMN dependency_name TYPE VARCHAR(255);
    ALTER COLUMN dependency_file TYPE VARCHAR(255);