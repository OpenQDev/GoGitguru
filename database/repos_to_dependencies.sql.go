// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: repos_to_dependencies.sql

package database

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const batchInsertRepoDependencies = `-- name: BatchInsertRepoDependencies :exec
WITH new_dependencies AS (
  INSERT INTO dependencies (dependency_file, dependency_name)
  SELECT
    DISTINCT dependency_file,
    dependency_name
  FROM (
    SELECT
      unnest($2::text[]) AS dependency_file,
      unnest($3::text[]) AS dependency_name
  ) AS subquery
  ON CONFLICT (dependency_file, dependency_name) DO NOTHING
  RETURNING internal_id, dependency_file, dependency_name
),
all_dependencies AS (
  SELECT internal_id, dependency_file, dependency_name FROM dependencies
  UNION 
  SELECT internal_id, dependency_file, dependency_name FROM new_dependencies
),
dependency_ids AS (
  SELECT
    s.updated_at::bigint as updated_at,
    d.internal_id AS dependency_id,
    s.url,
    s.firstUseDates AS first_use_date,
    s.lastUseDates AS last_use_date
  FROM (
    SELECT
  $1 as updated_at,
      $4 AS url,
      unnest($2::text[]) AS dependency_file,
      unnest($3::text[]) AS dependency_name,
unnest(COALESCE($5::bigint[], ARRAY[]::bigint[]))AS firstUseDates,
unnest(COALESCE($6::bigint[], ARRAY[]::bigint[])) AS lastUseDates
  ) s
  JOIN all_dependencies d ON d.dependency_file = s.dependency_file AND d.dependency_name = s.dependency_name
)
INSERT INTO repos_to_dependencies (
  url,
  updated_at,
  dependency_id, 
  first_use_date,
  last_use_date
) 
SELECT DISTINCT
  url,
  updated_at,
  dependency_id,  
  first_use_date,
  last_use_date
FROM
  dependency_ids
ON CONFLICT (url, dependency_id) DO UPDATE 
SET 
  first_use_date = EXCLUDED.first_use_date,
  last_use_date = EXCLUDED.last_use_date
`

type BatchInsertRepoDependenciesParams struct {
	UpdatedAt       sql.NullInt64 `json:"updated_at"`
	Filenames       []string      `json:"filenames"`
	Dependencynames []string      `json:"dependencynames"`
	Url             string        `json:"url"`
	Firstusedates   []int64       `json:"firstusedates"`
	Lastusedates    []int64       `json:"lastusedates"`
}

func (q *Queries) BatchInsertRepoDependencies(ctx context.Context, arg BatchInsertRepoDependenciesParams) error {
	_, err := q.exec(ctx, q.batchInsertRepoDependenciesStmt, batchInsertRepoDependencies,
		arg.UpdatedAt,
		pq.Array(arg.Filenames),
		pq.Array(arg.Dependencynames),
		arg.Url,
		pq.Array(arg.Firstusedates),
		pq.Array(arg.Lastusedates),
	)
	return err
}

const getRepoDependencies = `-- name: GetRepoDependencies :many
SELECT 
d.dependency_name,
rd.first_use_date,
rd.last_use_date
FROM dependencies d
LEFT JOIN repos_to_dependencies rd ON d.internal_id = rd.dependency_id
WHERE d.dependency_name = $1 AND rd.url = $2 AND d.dependency_file = ANY($3::text[])
`

type GetRepoDependenciesParams struct {
	DependencyName string   `json:"dependency_name"`
	Url            string   `json:"url"`
	Column3        []string `json:"column_3"`
}

type GetRepoDependenciesRow struct {
	DependencyName string        `json:"dependency_name"`
	FirstUseDate   sql.NullInt64 `json:"first_use_date"`
	LastUseDate    sql.NullInt64 `json:"last_use_date"`
}

func (q *Queries) GetRepoDependencies(ctx context.Context, arg GetRepoDependenciesParams) ([]GetRepoDependenciesRow, error) {
	rows, err := q.query(ctx, q.getRepoDependenciesStmt, getRepoDependencies, arg.DependencyName, arg.Url, pq.Array(arg.Column3))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRepoDependenciesRow
	for rows.Next() {
		var i GetRepoDependenciesRow
		if err := rows.Scan(&i.DependencyName, &i.FirstUseDate, &i.LastUseDate); err != nil {
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

const getRepoDependenciesByURL = `-- name: GetRepoDependenciesByURL :many
SELECT 
d.dependency_name,
d.dependency_file,
rd.first_use_date,
rd.last_use_date
FROM dependencies d
LEFT JOIN repos_to_dependencies rd ON d.internal_id = rd.dependency_id
WHERE rd.url = $1
`

type GetRepoDependenciesByURLRow struct {
	DependencyName string        `json:"dependency_name"`
	DependencyFile string        `json:"dependency_file"`
	FirstUseDate   sql.NullInt64 `json:"first_use_date"`
	LastUseDate    sql.NullInt64 `json:"last_use_date"`
}

func (q *Queries) GetRepoDependenciesByURL(ctx context.Context, url string) ([]GetRepoDependenciesByURLRow, error) {
	rows, err := q.query(ctx, q.getRepoDependenciesByURLStmt, getRepoDependenciesByURL, url)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRepoDependenciesByURLRow
	for rows.Next() {
		var i GetRepoDependenciesByURLRow
		if err := rows.Scan(
			&i.DependencyName,
			&i.DependencyFile,
			&i.FirstUseDate,
			&i.LastUseDate,
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

const initializeRepoDependencies = `-- name: InitializeRepoDependencies :exec
INSERT INTO repos_to_dependencies (url, dependency_id) VALUES (  
 $1,  
  unnest($2::int[]) 
)
ON CONFLICT (url, dependency_id) DO NOTHING
`

type InitializeRepoDependenciesParams struct {
	Url     string  `json:"url"`
	Column2 []int32 `json:"column_2"`
}

func (q *Queries) InitializeRepoDependencies(ctx context.Context, arg InitializeRepoDependenciesParams) error {
	_, err := q.exec(ctx, q.initializeRepoDependenciesStmt, initializeRepoDependencies, arg.Url, pq.Array(arg.Column2))
	return err
}
