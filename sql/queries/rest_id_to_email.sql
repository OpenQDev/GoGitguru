-- name: InsertRestIdToEmail :exec
INSERT INTO github_user_rest_id_author_emails (
    rest_id,
    email
) VALUES (
    $1, $2
)
ON CONFLICT DO NOTHING;

-- name: CheckGithubUserRestIdAuthorEmailExists :one
SELECT EXISTS(SELECT 1 FROM github_user_rest_id_author_emails WHERE email = $1);
