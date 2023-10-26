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

-- name: BulkInsertCommits :exec
INSERT INTO commits (commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, files_changed, repo_url) VALUES (  
  unnest($1::varchar[]),  
  unnest($2::varchar[]),  
  unnest($3::varchar[]),  
  unnest($4::bigint[]),  
  unnest($5::bigint[]),  
  unnest($6::text[]),  
  unnest($7::int[]),  
  unnest($8::int[]),  
  unnest($9::int[]),  
  unnest($10::varchar[])  
);

-- name: GetLatestUncheckedCommitPerAuthor :many
WITH LatestUncheckedCommitPerAuthor AS (
    SELECT DISTINCT ON (author_email)
    commit_hash,
    author_email,
    repo_url
    FROM commits
    WHERE author_email NOT IN (
        SELECT email FROM github_user_rest_id_author_emails
    )
    ORDER BY author_email, author_date DESC
)
SELECT commit_hash, author_email, repo_url
FROM LatestUncheckedCommitPerAuthor
ORDER BY repo_url DESC;

-- name: MultiRowInsertCommits :exec
INSERT INTO commits (commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, files_changed, repo_url) VALUES (  
  unnest($1::varchar[]),  
  unnest($2::varchar[]),  
  unnest($3::varchar[]),  
  unnest($4::bigint[]),  
  unnest($5::bigint[]),  
  unnest($6::text[]),  
  unnest($7::int[]),  
  unnest($8::int[]),  
  unnest($9::int[]),  
  unnest($10::varchar[])  
);

-- name: GetUserCommitsForRepos :many
WITH commits AS (
    SELECT * FROM commits WHERE c.author_date BETWEEN $1 AND $2
)
SELECT * FROM commits c
INNER JOIN github_user_rest_id_author_emails gure
ON c.author_email = gure.email
INNER JOIN github_users gu
ON gure.rest_id = gu.github_rest_id
WHERE gu.login = $3
AND c.repo_url = ANY($4)
ORDER BY c.author_date DESC;

-- name: GetAllUserCommits :many
WITH commits_cte AS (
    SELECT * FROM commits WHERE author_date BETWEEN $1 AND $2
)
SELECT * FROM commits_cte c
INNER JOIN github_user_rest_id_author_emails gure
ON c.author_email = gure.email
INNER JOIN github_users gu
ON gure.rest_id = gu.github_rest_id
WHERE gu.login = $3
ORDER BY c.author_date DESC;

-- name: GetLatestCommitterDate :one
SELECT committer_date + 1 AS next_committer_date
FROM commits
WHERE repo_url = CAST($1 AS VARCHAR)
ORDER BY committer_date DESC
LIMIT 1;
