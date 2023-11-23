-- +goose Up

CREATE TABLE dependencies_to_github_repos (
    internal_id SERIAL PRIMARY KEY,
	github__repo_internal_id VARCHAR(150) REFERENCES github_repos(internal_id),
    dependency_name VARCHAR(150) REFERENCES dependencies(dependency_name),
	first_commit_date DATE,
	date_added DATE,
	date_removed DATE,
	UNIQUE(dependency_name)
);