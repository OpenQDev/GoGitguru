-- +goose Up

CREATE INDEX idx_author_email ON commits (author_email);
CREATE INDEX idx_email ON github_user_rest_id_author_emails (email);

-- +goose Down

DROP INDEX idx_author_email ON commits;
DROP INDEX idx_email ON github_user_rest_id_author_emails;