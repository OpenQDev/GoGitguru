-- +goose Up
ALTER TABLE users_to_dependencies
ADD COLUMN resync_all BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE users_to_dependencies
DROP COLUMN resync_all;