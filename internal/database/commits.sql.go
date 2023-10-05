// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: commits.sql

package database

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const bulkInsertCommits = `-- name: BulkInsertCommits :exec
INSERT INTO commits (commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, files_changed, repo_url) VALUES (  
  unnest($1::varchar[]),  
  unnest($2::varchar[]),  
  unnest($3::varchar[]),  
  unnest($4::bigint[]),  
  unnest($5::bigint[]),  
  unnest($6::text[]),  
  unnest($7::int[]),  
  unnest($8::int[]),  
  unnest($9::int[]),  
  unnest($10::varchar[])  
)
`

type BulkInsertCommitsParams struct {
	Column1  []string `json:"column_1"`
	Column2  []string `json:"column_2"`
	Column3  []string `json:"column_3"`
	Column4  []int64  `json:"column_4"`
	Column5  []int64  `json:"column_5"`
	Column6  []string `json:"column_6"`
	Column7  []int32  `json:"column_7"`
	Column8  []int32  `json:"column_8"`
	Column9  []int32  `json:"column_9"`
	Column10 []string `json:"column_10"`
}

func (q *Queries) BulkInsertCommits(ctx context.Context, arg BulkInsertCommitsParams) error {
	_, err := q.exec(ctx, q.bulkInsertCommitsStmt, bulkInsertCommits,
		pq.Array(arg.Column1),
		pq.Array(arg.Column2),
		pq.Array(arg.Column3),
		pq.Array(arg.Column4),
		pq.Array(arg.Column5),
		pq.Array(arg.Column6),
		pq.Array(arg.Column7),
		pq.Array(arg.Column8),
		pq.Array(arg.Column9),
		pq.Array(arg.Column10),
	)
	return err
}

const getCommit = `-- name: GetCommit :one
SELECT commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, lines_changed, files_changed, repo_url FROM commits WHERE commit_hash = $1
`

func (q *Queries) GetCommit(ctx context.Context, commitHash string) (Commit, error) {
	row := q.queryRow(ctx, q.getCommitStmt, getCommit, commitHash)
	var i Commit
	err := row.Scan(
		&i.CommitHash,
		&i.Author,
		&i.AuthorEmail,
		&i.AuthorDate,
		&i.CommitterDate,
		&i.Message,
		&i.Insertions,
		&i.Deletions,
		&i.LinesChanged,
		&i.FilesChanged,
		&i.RepoUrl,
	)
	return i, err
}

const getCommits = `-- name: GetCommits :many
SELECT commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, lines_changed, files_changed, repo_url FROM commits
`

func (q *Queries) GetCommits(ctx context.Context) ([]Commit, error) {
	rows, err := q.query(ctx, q.getCommitsStmt, getCommits)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Commit
	for rows.Next() {
		var i Commit
		if err := rows.Scan(
			&i.CommitHash,
			&i.Author,
			&i.AuthorEmail,
			&i.AuthorDate,
			&i.CommitterDate,
			&i.Message,
			&i.Insertions,
			&i.Deletions,
			&i.LinesChanged,
			&i.FilesChanged,
			&i.RepoUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCommitsWithAuthorInfo = `-- name: GetCommitsWithAuthorInfo :many
SELECT commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, lines_changed, files_changed, repo_url, rest_id, gure.email, internal_id, github_rest_id, github_graphql_id, login, name, gu.email, avatar_url, company, location, bio, blog, hireable, twitter_username, followers, following, type, created_at, updated_at
FROM (
    SELECT commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, lines_changed, files_changed, repo_url
    FROM commits
    WHERE repo_url = $1
    AND author_date BETWEEN $2 AND $3
) c
INNER JOIN github_user_rest_id_author_emails gure
ON c.author_email = gure.email
INNER JOIN github_users gu
ON gure.rest_id = gu.github_rest_id
ORDER BY author_date DESC
`

type GetCommitsWithAuthorInfoParams struct {
	RepoUrl      sql.NullString `json:"repo_url"`
	AuthorDate   sql.NullInt64  `json:"author_date"`
	AuthorDate_2 sql.NullInt64  `json:"author_date_2"`
}

type GetCommitsWithAuthorInfoRow struct {
	CommitHash      string         `json:"commit_hash"`
	Author          sql.NullString `json:"author"`
	AuthorEmail     sql.NullString `json:"author_email"`
	AuthorDate      sql.NullInt64  `json:"author_date"`
	CommitterDate   sql.NullInt64  `json:"committer_date"`
	Message         sql.NullString `json:"message"`
	Insertions      sql.NullInt32  `json:"insertions"`
	Deletions       sql.NullInt32  `json:"deletions"`
	LinesChanged    sql.NullInt32  `json:"lines_changed"`
	FilesChanged    sql.NullInt32  `json:"files_changed"`
	RepoUrl         sql.NullString `json:"repo_url"`
	RestID          sql.NullInt32  `json:"rest_id"`
	Email           string         `json:"email"`
	InternalID      int32          `json:"internal_id"`
	GithubRestID    int32          `json:"github_rest_id"`
	GithubGraphqlID string         `json:"github_graphql_id"`
	Login           string         `json:"login"`
	Name            sql.NullString `json:"name"`
	Email_2         sql.NullString `json:"email_2"`
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

func (q *Queries) GetCommitsWithAuthorInfo(ctx context.Context, arg GetCommitsWithAuthorInfoParams) ([]GetCommitsWithAuthorInfoRow, error) {
	rows, err := q.query(ctx, q.getCommitsWithAuthorInfoStmt, getCommitsWithAuthorInfo, arg.RepoUrl, arg.AuthorDate, arg.AuthorDate_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCommitsWithAuthorInfoRow
	for rows.Next() {
		var i GetCommitsWithAuthorInfoRow
		if err := rows.Scan(
			&i.CommitHash,
			&i.Author,
			&i.AuthorEmail,
			&i.AuthorDate,
			&i.CommitterDate,
			&i.Message,
			&i.Insertions,
			&i.Deletions,
			&i.LinesChanged,
			&i.FilesChanged,
			&i.RepoUrl,
			&i.RestID,
			&i.Email,
			&i.InternalID,
			&i.GithubRestID,
			&i.GithubGraphqlID,
			&i.Login,
			&i.Name,
			&i.Email_2,
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
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLatestUncheckedCommitPerAuthor = `-- name: GetLatestUncheckedCommitPerAuthor :many

WITH LatestUncheckedCommitPerAuthor AS (
    SELECT DISTINCT ON (author_email)
    commit_hash,
    author_email,
    repo_url
    FROM commits
    WHERE author_email NOT IN (
        SELECT email FROM github_user_rest_id_author_emails
    )
    ORDER BY author_email, author_date DESC
)
SELECT commit_hash, author_email, repo_url
FROM LatestUncheckedCommitPerAuthor
ORDER BY repo_url DESC
`

type GetLatestUncheckedCommitPerAuthorRow struct {
	CommitHash  string         `json:"commit_hash"`
	AuthorEmail sql.NullString `json:"author_email"`
	RepoUrl     sql.NullString `json:"repo_url"`
}

func (q *Queries) GetLatestUncheckedCommitPerAuthor(ctx context.Context) ([]GetLatestUncheckedCommitPerAuthorRow, error) {
	rows, err := q.query(ctx, q.getLatestUncheckedCommitPerAuthorStmt, getLatestUncheckedCommitPerAuthor)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLatestUncheckedCommitPerAuthorRow
	for rows.Next() {
		var i GetLatestUncheckedCommitPerAuthorRow
		if err := rows.Scan(&i.CommitHash, &i.AuthorEmail, &i.RepoUrl); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertCommit = `-- name: InsertCommit :one
INSERT INTO commits (commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, files_changed, repo_url) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, lines_changed, files_changed, repo_url
`

type InsertCommitParams struct {
	CommitHash    string         `json:"commit_hash"`
	Author        sql.NullString `json:"author"`
	AuthorEmail   sql.NullString `json:"author_email"`
	AuthorDate    sql.NullInt64  `json:"author_date"`
	CommitterDate sql.NullInt64  `json:"committer_date"`
	Message       sql.NullString `json:"message"`
	Insertions    sql.NullInt32  `json:"insertions"`
	Deletions     sql.NullInt32  `json:"deletions"`
	FilesChanged  sql.NullInt32  `json:"files_changed"`
	RepoUrl       sql.NullString `json:"repo_url"`
}

func (q *Queries) InsertCommit(ctx context.Context, arg InsertCommitParams) (Commit, error) {
	row := q.queryRow(ctx, q.insertCommitStmt, insertCommit,
		arg.CommitHash,
		arg.Author,
		arg.AuthorEmail,
		arg.AuthorDate,
		arg.CommitterDate,
		arg.Message,
		arg.Insertions,
		arg.Deletions,
		arg.FilesChanged,
		arg.RepoUrl,
	)
	var i Commit
	err := row.Scan(
		&i.CommitHash,
		&i.Author,
		&i.AuthorEmail,
		&i.AuthorDate,
		&i.CommitterDate,
		&i.Message,
		&i.Insertions,
		&i.Deletions,
		&i.LinesChanged,
		&i.FilesChanged,
		&i.RepoUrl,
	)
	return i, err
}
