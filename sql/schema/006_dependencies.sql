-- +goose Up

CREATE TABLE dependencies (
    internal_id SERIAL PRIMARY KEY,
    dependency_name VARCHAR(150) NOT NULL,
	dependency_files TEXT[],
	UNIQUE(dependency_name)
);