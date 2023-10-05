-- name: InsertRestIdToEmail :one

INSERT INTO github_user_rest_id_author_emails (
    rest_id,
    email
) VALUES (
    $1, $2
)
RETURNING *;
