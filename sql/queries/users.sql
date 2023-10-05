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
RETURNING *;

-- name: GetGithubUser :one

SELECT * FROM github_users WHERE login = $1;

-- name: GetGroupOfEmails :one
SELECT github_rest_id FROM github_users WHERE github_rest_id = ANY($1::INT[]);
