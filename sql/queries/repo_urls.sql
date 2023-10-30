-- name: GetRepoURL :one
SELECT * FROM repo_urls WHERE url = $1;

-- name: GetRepoURLs :many
SELECT * FROM repo_urls;

-- name: InsertRepoURL :exec
INSERT INTO repo_urls (url, created_at, updated_at) 
VALUES ($1, NOW(), NOW())
RETURNING *;

-- name: UpdateStatusAndUpdatedAt :exec
UPDATE repo_urls SET status = $1, updated_at = NOW() WHERE url = $2;

-- name: GetReposStatus :many
SELECT
    url,
    status,
    (SELECT COUNT(DISTINCT(author_email))
     FROM commits
     WHERE author_email NOT IN (SELECT email FROM github_user_rest_id_author_emails)
     AND repo_url = url
    ) AS pending_authors
FROM repo_urls
WHERE url = ANY($1::text[])
ORDER BY status, updated_at DESC;

-- name: GetAndUpdateRepoURL :one
BEGIN;

WITH selected AS (
    SELECT url FROM repo_urls 
    WHERE status IN ('synced'::repo_status, 'pending'::repo_status, 'failed'::repo_status)
    AND (updated_at < NOW() - INTERVAL $1 OR updated_at IS NULL) 
    ORDER BY RANDOM() LIMIT 1
    FOR UPDATE
), updated AS (
    UPDATE repo_urls SET status = 'queued'::repo_status, updated_at = NOW() WHERE url IN (SELECT url FROM selected)
    RETURNING *
)

SELECT * FROM updated;

COMMIT;