-- +goose Up

ALTER TABLE commits ADD COLUMN has_checked_user BOOLEAN DEFAULT TRUE;

-- +goose Down

ALTER TABLE commits DROP COLUMN has_checked_user;