-- +goose Up

ALTER TABLE repos_to_dependencies
    ALTER COLUMN url TYPE VARCHAR(255);
ALTER TABLE github_repos
    ALTER COLUMN url TYPE VARCHAR(255);
ALTER TABLE repo_urls
    ALTER COLUMN url TYPE VARCHAR(255);