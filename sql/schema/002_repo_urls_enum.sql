-- +goose Up
CREATE TYPE repo_status AS ENUM ('pending', 'queued', 'syncing_database', 'syncing_repo', 'synced', 'storing_commits', 'failed', 'not_listed');

-- In order to add a repo_status typed default to status, we need to
-- 1) DROP the default
ALTER TABLE repo_urls ALTER COLUMN status DROP DEFAULT;

-- 2) TYPE the status column to repo_status
ALTER TABLE repo_urls ALTER COLUMN status TYPE repo_status USING status::repo_status;

-- 3) SET 'pending' as the new default
ALTER TABLE repo_urls ALTER COLUMN status SET DEFAULT 'pending';

-- +goose Down
-- MUST drop all dependents of repo_status enum before removing it as a requirement on repo_urls::status
ALTER TABLE repo_urls ALTER COLUMN status DROP DEFAULT;
ALTER TABLE repo_urls ALTER COLUMN status TYPE VARCHAR(30) USING status::VARCHAR;
DROP TYPE repo_status;
ALTER TABLE repo_urls ALTER COLUMN status SET DEFAULT 'pending';
