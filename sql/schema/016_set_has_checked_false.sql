-- +goose Up

ALTER TABLE commits ALTER COLUMN has_checked_user SET DEFAULT FALSE;

-- +goose Down

ALTER TABLE commits ALTER COLUMN has_checked_user SET DEFAULT TRUE;