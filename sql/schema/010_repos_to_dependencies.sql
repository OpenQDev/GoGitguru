-- +goose Up

CREATE TABLE repos_to_dependencies (
    url VARCHAR(150),
    dependency_id INT NOT NULL,
    first_use_data BIGINT DEFAULT NULL,
    last_use_data BIGINT DEFAULT NULL,
    updated_at  BIGINT DEFAULT NULL,
    status VARCHAR(10) CHECK (status IN ('queued', 'synced')) DEFAULT 'queued',
    PRIMARY KEY (url, dependency_id),
    FOREIGN KEY (url) REFERENCES repo_urls(url),
    FOREIGN KEY (dependency_id) REFERENCES dependencies (internal_id) ON DELETE CASCADE,
    UNIQUE (url, dependency_id)

);

