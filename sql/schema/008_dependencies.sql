-- +goose Up

CREATE TABLE dependencies (
    internal_id SERIAL PRIMARY KEY,
    dependency_name VARCHAR(255) NOT NULL,
    dependency_file VARCHAR(255) NOT NULL,
    UNIQUE (dependency_name, dependency_file)
);