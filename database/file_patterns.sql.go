// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: file_patterns.sql

package database

import (
	"context"

	"github.com/lib/pq"
)

const bulkUpsertFilePatterns = `-- name: BulkUpsertFilePatterns :exec
INSERT INTO file_patterns (pattern, updated_at, creator) VALUES (unnest($1::text[]), $2::bigint, $3::text)
ON CONFLICT DO NOTHING
`

type BulkUpsertFilePatternsParams struct {
	Patterns  []string `json:"patterns"`
	UpdatedAt int64    `json:"updated_at"`
	Creator   string   `json:"creator"`
}

func (q *Queries) BulkUpsertFilePatterns(ctx context.Context, arg BulkUpsertFilePatternsParams) error {
	_, err := q.exec(ctx, q.bulkUpsertFilePatternsStmt, bulkUpsertFilePatterns, pq.Array(arg.Patterns), arg.UpdatedAt, arg.Creator)
	return err
}

const getAllFilePatterns = `-- name: GetAllFilePatterns :many
SELECT id, pattern, updated_at, creator FROM file_patterns
LIMIT 50
`

func (q *Queries) GetAllFilePatterns(ctx context.Context) ([]FilePattern, error) {
	rows, err := q.query(ctx, q.getAllFilePatternsStmt, getAllFilePatterns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FilePattern
	for rows.Next() {
		var i FilePattern
		if err := rows.Scan(
			&i.ID,
			&i.Pattern,
			&i.UpdatedAt,
			&i.Creator,
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
