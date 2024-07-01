-- +goose Up

CREATE TABLE users_to_dependencies (
    dependency_id INT NOT NULL,
    user_id INT NOT NULL,
    first_use_date BIGINT DEFAULT NULL,
    last_use_date BIGINT DEFAULT NULL,
    updated_at  BIGINT DEFAULT NULL,
    PRIMARY KEY (dependency_id, user_id),
    FOREIGN KEY (dependency_id) REFERENCES dependencies (internal_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES github_users (internal_id) ON DELETE CASCADE,
    UNIQUE (dependency_id, user_id)
);
