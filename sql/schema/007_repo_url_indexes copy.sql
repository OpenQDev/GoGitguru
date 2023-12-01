-- +goose Up

CREATE INDEX idx_repo_url ON commits (repo_url);
CREATE INDEX idx_full_name ON github_repos (full_name);
CREATE INDEX idx_url ON repo_urls (url);
CREATE INDEX idx_rest_id ON github_user_rest_id_author_emails (rest_id);
CREATE INDEX idx_login ON github_users (login);
