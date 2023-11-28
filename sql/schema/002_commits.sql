-- +goose Up

CREATE TABLE commits (
    commit_hash VARCHAR(40) NOT NULL,
    author VARCHAR(255),
    author_email VARCHAR(255),
    author_date BIGINT,
    committer_date BIGINT,
    message TEXT,
    insertions INT,
    deletions INT,
    lines_changed INT GENERATED ALWAYS AS (insertions + deletions) STORED,
    files_changed INT,
    repo_url VARCHAR(255) REFERENCES repo_urls(url),
    UNIQUE(commit_hash, repo_url)
);