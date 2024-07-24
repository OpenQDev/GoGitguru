-- name: GetCommit :one
SELECT * FROM commits WHERE commit_hash = $1;

-- name: GetCommits :many
SELECT * FROM commits;


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
INSERT INTO commits (
    commit_hash, 
    author, 
    author_email, 
    author_date, 
    committer_date, 
    message,
    files_changed, 
    repo_url
) 
SELECT
    unnest(sqlc.arg(commitHashes)::varchar[]),  
    unnest(sqlc.arg(authors)::varchar[]),  
    unnest(sqlc.arg(authorEmails)::varchar[]),  
    unnest(sqlc.arg(authorDates)::bigint[]),  
    unnest(sqlc.arg(committerDates)::bigint[]),  
    unnest(sqlc.arg(messages)::text[]),
    unnest(sqlc.arg(filesChanged)::int[]),  
    sqlc.arg(repoUrl)
ON CONFLICT (commit_hash, repo_url) DO NOTHING;

-- name: GetLatestUncheckedCommitPerAuthor :many
SELECT DISTINCT ON (c.author_email)
c.commit_hash,
c.author_email,
c.author_date,
c.repo_url,
g.email as github_user_email
FROM (
    SELECT * FROM commits WHERE has_checked_user = FALSE
) c
LEFT JOIN github_user_rest_id_author_emails g
ON c.author_email = g.email
ORDER BY c.author_email, c.author_date DESC;


-- name: SetAllCommitsToChecked :exec
UPDATE commits
SET has_checked_user = TRUE
WHERE commits.author_email = ANY($1::VARCHAR[]);


-- name: GetUserCommitsForRepos :many
WITH commits_cte AS (
    SELECT * FROM commits WHERE author_date BETWEEN $1 AND $2
)
SELECT * FROM commits_cte c
INNER JOIN github_user_rest_id_author_emails gure
ON c.author_email = gure.email
INNER JOIN github_users gu
ON gure.rest_id = gu.github_rest_id
WHERE gu.login = $3
AND c.repo_url = ANY($4::VARCHAR[])
ORDER BY c.author_date DESC;

-- name: GetLatestCommitterDate :one
SELECT (committer_date + 1)::bigint AS next_committer_date
FROM commits
WHERE repo_url = CAST($1 AS VARCHAR)
ORDER BY committer_date DESC
LIMIT 1;

-- name: GetFirstCommit :one
SELECT * FROM commits c
INNER JOIN github_user_rest_id_author_emails gure
ON c.author_email = gure.email
INNER JOIN github_users gu
ON gure.rest_id = gu.github_rest_id
WHERE c.repo_url = $1
AND gu.login ILIKE $2
ORDER BY c.author_date ASC
LIMIT 1;

-- name: GetFirstAndLastCommit :one
SELECT 
    MIN(author_date) as first_commit_date,
    MAX(author_date) as last_commit_date
FROM commits
WHERE author_email = $1;
