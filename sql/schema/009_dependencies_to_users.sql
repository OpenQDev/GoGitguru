-- +goose Up

CREATE TABLE dependencies_to_users (
    dependency_id INT NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (dependency_id, user_id),
    FOREIGN KEY (dependency_id) REFERENCES dependencies (internal_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES github_users (internal_id) ON DELETE CASCADE
);
