-- name: InsertUser :one

INSERT INTO github_users (
    github_rest_id,
    github_graphql_id,
    login,
    name,
    email,
    avatar_url,
    company,
    location,
    bio,
    blog,
    hireable,
    twitter_username,
    followers,
    following,
    type,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
)
RETURNING internal_id;

-- name: CheckGithubUserId :one
SELECT internal_id FROM github_users WHERE login = $1
LIMIT 1;

-- name: CheckGithubUserExists :one
SELECT EXISTS(SELECT 1 FROM github_users WHERE login = $1);

-- name: GetGithubUserByCommitEmail :many
SELECT gu.internal_id, array_agg(DISTINCT gure.email)::text[] AS emails FROM github_users gu
INNER JOIN github_user_rest_id_author_emails gure
ON gu.github_rest_id = gure.rest_id
WHERE gure.email = ANY(sqlc.arg(user_emails)::TEXT[])
GROUP BY gu.internal_id;

-- name: GetGithubUser :one
SELECT * FROM github_users WHERE login = $1;

-- name: GetGroupOfEmails :one
SELECT github_rest_id FROM github_users WHERE github_rest_id = ANY($1::INT[]);
