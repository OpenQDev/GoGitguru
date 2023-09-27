-- name: GetCommit :one
SELECT * FROM commits WHERE commit_hash = $1;

-- name: GetCommits :many
SELECT * FROM commits;

-- name: InsertCommit :exec
INSERT INTO commits (commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, files_changed, repo_url) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;
