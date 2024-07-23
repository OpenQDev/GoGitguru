-- name: InsertUser :exec

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
ON CONFLICT (github_rest_id) DO UPDATE
SET
    github_graphql_id = EXCLUDED.github_graphql_id,
    login = EXCLUDED.login,
    name = EXCLUDED.name,
    email = EXCLUDED.email,
    avatar_url = EXCLUDED.avatar_url,
    company = EXCLUDED.company,
    location = EXCLUDED.location,
    bio = EXCLUDED.bio,
    blog = EXCLUDED.blog,
    hireable = EXCLUDED.hireable,
    twitter_username = EXCLUDED.twitter_username,
    followers = EXCLUDED.followers,
    following = EXCLUDED.following,
    type = EXCLUDED.type,
    created_at = EXCLUDED.created_at,
    updated_at = EXCLUDED.updated_at;

-- name: CheckGithubUserIdExists :one
SELECT EXISTS(SELECT 1 FROM github_users WHERE github_rest_id = $1);
-- name: GetGithubUserByRestId :one
SELECT 1 FROM github_users WHERE github_rest_id = $1;


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
