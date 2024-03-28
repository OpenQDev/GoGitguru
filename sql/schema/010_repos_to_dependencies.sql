-- +goose Up

CREATE TABLE repos_to_dependencies (
    url VARCHAR(150),
    dependency_id VARCHAR(255) UNIQUE NOT NULL,
    first_use_data TIMESTAMP DEFAULT NULL,
    last_use_data TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (url, dependency_id),
    FOREIGN KEY (url) REFERENCES repo_urls(url),
    FOREIGN KEY (dependency_id) REFERENCES dependencies (internal_id) ON DELETE CASCADE,
    UNIQUE (url, dependency_id)
)

