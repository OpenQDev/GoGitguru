-- +goose Up
CREATE TYPE repo_status AS ENUM ('pending', 'queued', 'syncing_repo', 'synced', 'storing_commits', 'failed', 'not_listed');

CREATE TABLE repo_urls (
    url VARCHAR(150) PRIMARY KEY,
    status repo_status NOT NULL DEFAULT 'pending'::repo_status,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL
);

-- +goose Down
DROP TABLE repo_urls;