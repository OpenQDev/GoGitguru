// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: repo_urls.sql

package database

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const deleteRepoURL = `-- name: DeleteRepoURL :exec
DELETE FROM repo_urls WHERE url = $1
`

func (q *Queries) DeleteRepoURL(ctx context.Context, url string) error {
	_, err := q.exec(ctx, q.deleteRepoURLStmt, deleteRepoURL, url)
	return err
}

const getAndUpdateRepoURL = `-- name: GetAndUpdateRepoURL :one
BEGIN
`

type GetAndUpdateRepoURLRow struct {
}

func (q *Queries) GetAndUpdateRepoURL(ctx context.Context) (GetAndUpdateRepoURLRow, error) {
	row := q.queryRow(ctx, q.getAndUpdateRepoURLStmt, getAndUpdateRepoURL)
	var i GetAndUpdateRepoURLRow
	err := row.Scan()
	return i, err
}

const getRepoURL = `-- name: GetRepoURL :one
SELECT url, status, created_at, updated_at FROM repo_urls WHERE url = $1
`

func (q *Queries) GetRepoURL(ctx context.Context, url string) (RepoUrl, error) {
	row := q.queryRow(ctx, q.getRepoURLStmt, getRepoURL, url)
	var i RepoUrl
	err := row.Scan(
		&i.Url,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getRepoURLs = `-- name: GetRepoURLs :many
SELECT url, status, created_at, updated_at FROM repo_urls
`

func (q *Queries) GetRepoURLs(ctx context.Context) ([]RepoUrl, error) {
	rows, err := q.query(ctx, q.getRepoURLsStmt, getRepoURLs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RepoUrl
	for rows.Next() {
		var i RepoUrl
		if err := rows.Scan(
			&i.Url,
			&i.Status,
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

const getReposStatus = `-- name: GetReposStatus :many
SELECT 
    r.url,
    r.status,
    r.updated_at,
    COUNT(DISTINCT c.author_email) FILTER (WHERE g.email IS NULL) AS pending_authors

FROM
    repo_urls r
LEFT JOIN commits c ON c.repo_url = r.url
LEFT JOIN github_user_rest_id_author_emails g ON c.author_email = g.email

WHERE
    r.url = ANY($1::text[])
GROUP BY
    r.url, r.status, r.updated_at
ORDER BY r.status, r.updated_at DESC
`

type GetReposStatusRow struct {
	Url            string       `json:"url"`
	Status         RepoStatus   `json:"status"`
	UpdatedAt      sql.NullTime `json:"updated_at"`
	PendingAuthors int64        `json:"pending_authors"`
}

func (q *Queries) GetReposStatus(ctx context.Context, dollar_1 []string) ([]GetReposStatusRow, error) {
	rows, err := q.query(ctx, q.getReposStatusStmt, getReposStatus, pq.Array(dollar_1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetReposStatusRow
	for rows.Next() {
		var i GetReposStatusRow
		if err := rows.Scan(
			&i.Url,
			&i.Status,
			&i.UpdatedAt,
			&i.PendingAuthors,
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

const updateStatusAndUpdatedAt = `-- name: UpdateStatusAndUpdatedAt :exec
UPDATE repo_urls SET status = $1, updated_at = NOW() WHERE url = $2
`

type UpdateStatusAndUpdatedAtParams struct {
	Status RepoStatus `json:"status"`
	Url    string     `json:"url"`
}

func (q *Queries) UpdateStatusAndUpdatedAt(ctx context.Context, arg UpdateStatusAndUpdatedAtParams) error {
	_, err := q.exec(ctx, q.updateStatusAndUpdatedAtStmt, updateStatusAndUpdatedAt, arg.Status, arg.Url)
	return err
}

const upsertRepoURL = `-- name: UpsertRepoURL :exec
INSERT INTO repo_urls (url, created_at, updated_at) 
VALUES ($1, NOW(), NOW())
ON CONFLICT (url) DO UPDATE SET updated_at = NOW(), status = 'pending'::repo_status
RETURNING url, status, created_at, updated_at
`

func (q *Queries) UpsertRepoURL(ctx context.Context, url string) error {
	_, err := q.exec(ctx, q.upsertRepoURLStmt, upsertRepoURL, url)
	return err
}
