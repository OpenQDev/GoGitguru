-- +goose Up

CREATE TABLE users_to_repo_urls (
    user_id INT NOT NULL,
    url VARCHAR(255) NOT NULL,
    first_commit_date BIGINT DEFAULT NULL,
    last_commit_date BIGINT DEFAULT NULL,
    updated_at  BIGINT DEFAULT NULL,
    PRIMARY KEY (user_id, url),
    FOREIGN KEY (user_id) REFERENCES github_users(internal_id) ON DELETE CASCADE,
    FOREIGN KEY (url) REFERENCES repo_urls(url),
    UNIQUE (user_id, url)
);