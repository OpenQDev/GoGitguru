-- name: GetGithubRepo :one
SELECT * FROM github_repos WHERE full_name = $1;

-- name: CheckGithubRepoExists :one
SELECT EXISTS(SELECT 1 FROM github_repos WHERE full_name = $1);

-- name: InsertGithubRepo :one

INSERT INTO github_repos (
    github_rest_id,
    github_graphql_id,
    url,
    name,
    full_name,
    private,
    owner_login,
    owner_avatar_url,
    description,
    homepage,
    fork,
    forks_count,
    archived,
    disabled,
    license,
    language,
    stargazers_count,
    watchers_count,
    open_issues_count,
    has_issues,
    has_discussions,
    has_projects,
    created_at,
    updated_at,
    pushed_at,
    visibility,
    size,
    default_branch
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28
)
ON CONFLICT (github_rest_id) DO NOTHING
RETURNING *;
