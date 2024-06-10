// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: commits.sql

package database

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const bulkInsertCommits = `-- name: BulkInsertCommits :exec
INSERT INTO commits (
    commit_hash, 
    author, 
    author_email, 
    author_date, 
    committer_date, 
    message,
    files_changed, 
    repo_url
) 
SELECT
    unnest($1::varchar[]),  
    unnest($2::varchar[]),  
    unnest($3::varchar[]),  
    unnest($4::bigint[]),  
    unnest($5::bigint[]),  
    unnest($6::text[]),
    unnest($7::int[]),  
    $8
ON CONFLICT (commit_hash, repo_url) DO NOTHING
`

type BulkInsertCommitsParams struct {
	Commithashes   []string       `json:"commithashes"`
	Authors        []string       `json:"authors"`
	Authoremails   []string       `json:"authoremails"`
	Authordates    []int64        `json:"authordates"`
	Committerdates []int64        `json:"committerdates"`
	Messages       []string       `json:"messages"`
	Fileschanged   []int32        `json:"fileschanged"`
	Repourl        sql.NullString `json:"repourl"`
}

func (q *Queries) BulkInsertCommits(ctx context.Context, arg BulkInsertCommitsParams) error {
	_, err := q.exec(ctx, q.bulkInsertCommitsStmt, bulkInsertCommits,
		pq.Array(arg.Commithashes),
		pq.Array(arg.Authors),
		pq.Array(arg.Authoremails),
		pq.Array(arg.Authordates),
		pq.Array(arg.Committerdates),
		pq.Array(arg.Messages),
		pq.Array(arg.Fileschanged),
		arg.Repourl,
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
	RestID          int32          `json:"rest_id"`
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

const getFirstCommit = `-- name: GetFirstCommit :one
SELECT commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, lines_changed, files_changed, repo_url, rest_id, gure.email, internal_id, github_rest_id, github_graphql_id, login, name, gu.email, avatar_url, company, location, bio, blog, hireable, twitter_username, followers, following, type, created_at, updated_at FROM commits c
INNER JOIN github_user_rest_id_author_emails gure
ON c.author_email = gure.email
INNER JOIN github_users gu
ON gure.rest_id = gu.github_rest_id
WHERE c.repo_url = $1
AND gu.login ILIKE $2
ORDER BY c.author_date ASC
LIMIT 1
`

type GetFirstCommitParams struct {
	RepoUrl sql.NullString `json:"repo_url"`
	Login   string         `json:"login"`
}

type GetFirstCommitRow struct {
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
	RestID          int32          `json:"rest_id"`
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

func (q *Queries) GetFirstCommit(ctx context.Context, arg GetFirstCommitParams) (GetFirstCommitRow, error) {
	row := q.queryRow(ctx, q.getFirstCommitStmt, getFirstCommit, arg.RepoUrl, arg.Login)
	var i GetFirstCommitRow
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
	)
	return i, err
}

const getLatestCommitterDate = `-- name: GetLatestCommitterDate :one
SELECT (committer_date + 1)::bigint AS next_committer_date
FROM commits
WHERE repo_url = CAST($1 AS VARCHAR)
ORDER BY committer_date DESC
LIMIT 1
`

func (q *Queries) GetLatestCommitterDate(ctx context.Context, dollar_1 string) (int64, error) {
	row := q.queryRow(ctx, q.getLatestCommitterDateStmt, getLatestCommitterDate, dollar_1)
	var next_committer_date int64
	err := row.Scan(&next_committer_date)
	return next_committer_date, err
}

const getLatestUncheckedCommitPerAuthor = `-- name: GetLatestUncheckedCommitPerAuthor :many
SELECT DISTINCT ON (c.author_email)
c.commit_hash,
c.author_email,
c.repo_url
FROM commits c
LEFT JOIN github_user_rest_id_author_emails g
ON c.author_email = g.email
WHERE g.email IS NULL
ORDER BY c.author_email, c.author_date DESC
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

const getUserCommitsForRepos = `-- name: GetUserCommitsForRepos :many
WITH commits_cte AS (
    SELECT commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, lines_changed, files_changed, repo_url FROM commits WHERE author_date BETWEEN $1 AND $2
)
SELECT commit_hash, author, author_email, author_date, committer_date, message, insertions, deletions, lines_changed, files_changed, repo_url, rest_id, gure.email, internal_id, github_rest_id, github_graphql_id, login, name, gu.email, avatar_url, company, location, bio, blog, hireable, twitter_username, followers, following, type, created_at, updated_at FROM commits_cte c
INNER JOIN github_user_rest_id_author_emails gure
ON c.author_email = gure.email
INNER JOIN github_users gu
ON gure.rest_id = gu.github_rest_id
WHERE gu.login = $3
AND c.repo_url = ANY($4::VARCHAR[])
ORDER BY c.author_date DESC
`

type GetUserCommitsForReposParams struct {
	AuthorDate   sql.NullInt64 `json:"author_date"`
	AuthorDate_2 sql.NullInt64 `json:"author_date_2"`
	Login        string        `json:"login"`
	Column4      []string      `json:"column_4"`
}

type GetUserCommitsForReposRow struct {
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
	RestID          int32          `json:"rest_id"`
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

func (q *Queries) GetUserCommitsForRepos(ctx context.Context, arg GetUserCommitsForReposParams) ([]GetUserCommitsForReposRow, error) {
	rows, err := q.query(ctx, q.getUserCommitsForReposStmt, getUserCommitsForRepos,
		arg.AuthorDate,
		arg.AuthorDate_2,
		arg.Login,
		pq.Array(arg.Column4),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserCommitsForReposRow
	for rows.Next() {
		var i GetUserCommitsForReposRow
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
