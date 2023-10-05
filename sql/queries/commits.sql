-- name: GetCommit :one
SELECT * FROM commits WHERE commit_hash = $1;

-- name: GetCommits :many
SELECT * FROM commits;

-- name: InsertCommit :one
INSERT INTO commits (commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, files_changed, repo_url) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetCommitsWithAuthorInfo :many
SELECT *
FROM (
    SELECT *
    FROM commits
    WHERE repo_url = $1
    AND author_date BETWEEN $2 AND $3
) c
INNER JOIN github_user_rest_id_author_emails gure
ON c.author_email = gure.email
INNER JOIN github_users gu
ON gure.rest_id = gu.github_rest_id
ORDER BY author_date DESC;