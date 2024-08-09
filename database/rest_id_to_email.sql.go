// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: rest_id_to_email.sql

package database

import (
	"context"
)

const checkGithubUserRestIdAuthorEmailExists = `-- name: CheckGithubUserRestIdAuthorEmailExists :one
SELECT EXISTS(SELECT 1 FROM github_user_rest_id_author_emails WHERE email = $1)
`

func (q *Queries) CheckGithubUserRestIdAuthorEmailExists(ctx context.Context, email string) (bool, error) {
	row := q.queryRow(ctx, q.checkGithubUserRestIdAuthorEmailExistsStmt, checkGithubUserRestIdAuthorEmailExists, email)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const insertRestIdToEmail = `-- name: InsertRestIdToEmail :exec
INSERT INTO github_user_rest_id_author_emails (
    rest_id,
    email
) VALUES (
    $1, $2
)
ON CONFLICT DO NOTHING
`

type InsertRestIdToEmailParams struct {
	RestID int32  `json:"rest_id"`
	Email  string `json:"email"`
}

func (q *Queries) InsertRestIdToEmail(ctx context.Context, arg InsertRestIdToEmailParams) error {
	_, err := q.exec(ctx, q.insertRestIdToEmailStmt, insertRestIdToEmail, arg.RestID, arg.Email)
	return err
}
