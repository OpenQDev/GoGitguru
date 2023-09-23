-- name: CreateRepo :one
INSERT INTO repos (id, created_at, updated_at, organization, repository, url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetRepos :many
SELECT * FROM repos;