-- name: GetRepoURL :one
SELECT * FROM repo_urls WHERE url = $1;

-- name: GetRepoURLs :many
SELECT * FROM repo_urls;

-- name: InsertRepoURL :exec
INSERT INTO repo_urls (url, created_at, updated_at) 
VALUES ($1, NOW(), NOW())
RETURNING *;
