-- name: GetRepoURL :one
SELECT * FROM repo_urls WHERE url = $1;

-- name: GetRepoURLs :many
SELECT * FROM repo_urls;

-- name: InsertRepoURL :exec
INSERT INTO repo_urls (url, created_at, updated_at) 
VALUES ($1, NOW(), NOW())
RETURNING *;

-- name: UpdateStatus :exec
UPDATE repo_urls SET status = $1 WHERE url = $2 AND status != 'failed';

-- name: UpdateStatusAndUpdatedAt :exec
UPDATE repo_urls SET status = $1, updated_at = NOW() WHERE url = $2 AND status != 'failed';