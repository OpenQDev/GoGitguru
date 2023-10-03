// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: commits.sql

package database

import (
	"context"
	"database/sql"
)

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