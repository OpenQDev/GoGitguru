-- +goose Up
CREATE TABLE If NOT EXISTS repo_urls (
    url VARCHAR(150) PRIMARY KEY,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NULL
);

-- +goose Down
DROP TABLE repo_urls;