-- +goose Up

CREATE TABLE dependencies_to_users (
    dependency_id INT NOT NULL,
    user_id INT NOT NULL,
    dependency_id VARCHAR(255) NOT NULL,
    first_use_data TIMESTAMP DEFAULT NULL,
    last_use_data TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (dependency_id, user_id),
    FOREIGN KEY (dependency_id) REFERENCES dependencies (internal_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES github_users (internal_id) ON DELETE CASCADE,
    UNIQUE (dependency_id, user_id)
);
