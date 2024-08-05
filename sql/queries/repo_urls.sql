-- name: GetRepoURL :one
SELECT * FROM repo_urls WHERE url = $1;

-- name: DeleteRepoURL :exec
DELETE FROM repo_urls WHERE url = $1;

-- name: GetRepoURLs :many
SELECT * FROM repo_urls;

-- name: UpsertRepoURL :exec
INSERT INTO repo_urls (url, created_at, updated_at) 
VALUES ($1, NOW(), NOW())
ON CONFLICT (url) DO UPDATE SET updated_at = NOW(), status = 'pending'::repo_status
RETURNING *;

-- name: UpdateStatusAndUpdatedAt :exec
UPDATE repo_urls SET status = $1, updated_at = NOW() WHERE url = $2;

-- name: GetReposStatus :many
SELECT 
    r.url,
    r.status,
    r.updated_at,
    COUNT(DISTINCT c.author_email) FILTER (WHERE g.email IS NULL) AS pending_authors

FROM
    repo_urls r
LEFT JOIN commits c ON c.repo_url = r.url
LEFT JOIN github_user_rest_id_author_emails g ON c.author_email = g.email

WHERE
    r.url = ANY($1::text[])
GROUP BY
    r.url, r.status, r.updated_at
ORDER BY r.status, r.updated_at DESC;

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


