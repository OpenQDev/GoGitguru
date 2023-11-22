-- name: InsertRestIdToEmail :one

INSERT INTO github_user_rest_id_author_emails (
    rest_id,
    email
) VALUES (
    $1, $2
)
RETURNING *;

-- name: CheckGithubUserRestIdAuthorEmailExists :one
SELECT EXISTS(SELECT 1 FROM github_user_rest_id_author_emails WHERE rest_id = $1 AND email = $2);
