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
ON CONFLICT (login) DO UPDATE SET
    github_rest_id = EXCLUDED.github_rest_id,
    github_graphql_id = EXCLUDED.github_graphql_id,
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
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: CheckGithubUserExists :one
SELECT EXISTS(SELECT 1 FROM github_users WHERE login = $1);

-- name: GetGithubUser :one
SELECT * FROM github_users WHERE login = $1;

-- name: GetGroupOfEmails :one
SELECT github_rest_id FROM github_users WHERE github_rest_id = ANY($1::INT[]);
