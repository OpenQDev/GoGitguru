-- +goose Up
CREATE TABLE repo_urls_v2 AS TABLE repo_urls WITH NO DATA;

-- +goose Down
DROP TABLE repo_urls_v2;
