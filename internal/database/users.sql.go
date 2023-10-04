// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
)

const getGithubUser = `-- name: GetGithubUser :one

SELECT internal_id, github_rest_id, github_graphql_id, login, name, email, avatar_url, company, location, bio, blog, hireable, twitter_username, followers, following, type, created_at, updated_at FROM github_users WHERE login = $1
`

func (q *Queries) GetGithubUser(ctx context.Context, login string) (GithubUser, error) {
	row := q.queryRow(ctx, q.getGithubUserStmt, getGithubUser, login)
	var i GithubUser
	err := row.Scan(
		&i.InternalID,
		&i.GithubRestID,
		&i.GithubGraphqlID,
		&i.Login,
		&i.Name,
		&i.Email,
		&i.AvatarUrl,
		&i.Company,
		&i.Location,
		&i.Bio,
		&i.Blog,
		&i.Hireable,
		&i.TwitterUsername,
		&i.Followers,
		&i.Following,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :one

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
RETURNING internal_id, github_rest_id, github_graphql_id, login, name, email, avatar_url, company, location, bio, blog, hireable, twitter_username, followers, following, type, created_at, updated_at
`

type InsertUserParams struct {
	GithubRestID    int32          `json:"github_rest_id"`
	GithubGraphqlID string         `json:"github_graphql_id"`
	Login           string         `json:"login"`
	Name            sql.NullString `json:"name"`
	Email           sql.NullString `json:"email"`
	AvatarUrl       sql.NullString `json:"avatar_url"`
	Company         sql.NullString `json:"company"`
	Location        sql.NullString `json:"location"`
	Bio             sql.NullString `json:"bio"`
	Blog            sql.NullString `json:"blog"`
	Hireable        sql.NullBool   `json:"hireable"`
	TwitterUsername sql.NullString `json:"twitter_username"`
	Followers       sql.NullInt32  `json:"followers"`
	Following       sql.NullInt32  `json:"following"`
	Type            string         `json:"type"`
	CreatedAt       sql.NullTime   `json:"created_at"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (GithubUser, error) {
	row := q.queryRow(ctx, q.insertUserStmt, insertUser,
		arg.GithubRestID,
		arg.GithubGraphqlID,
		arg.Login,
		arg.Name,
		arg.Email,
		arg.AvatarUrl,
		arg.Company,
		arg.Location,
		arg.Bio,
		arg.Blog,
		arg.Hireable,
		arg.TwitterUsername,
		arg.Followers,
		arg.Following,
		arg.Type,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i GithubUser
	err := row.Scan(
		&i.InternalID,
		&i.GithubRestID,
		&i.GithubGraphqlID,
		&i.Login,
		&i.Name,
		&i.Email,
		&i.AvatarUrl,
		&i.Company,
		&i.Location,
		&i.Bio,
		&i.Blog,
		&i.Hireable,
		&i.TwitterUsername,
		&i.Followers,
		&i.Following,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}