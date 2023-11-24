-- +goose Up

CREATE TABLE dependencies_to_repos (
    internal_id SERIAL PRIMARY KEY,
	github_repo_internal_id SERIAL REFERENCES github_repos(internal_id),
    dependency_name VARCHAR(150) REFERENCES dependencies(dependency_name),
	first_commit_date INT,
	date_added INT,
	date_removed INT,
	UNIQUE(dependency_name, github_repo_internal_id)
);