-- +goose Up

CREATE TABLE github_user_rest_id_author_emails (
    rest_id INT DEFAULT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    PRIMARY KEY (rest_id, email)
);