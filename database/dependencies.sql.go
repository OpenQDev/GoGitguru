// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: dependencies.sql

package database

import (
	"context"

	"github.com/lib/pq"
)

const bulkInsertDependencies = `-- name: BulkInsertDependencies :many

INSERT INTO dependencies (
    dependency_name,
    dependency_file
) VALUES (
    $1, unnest($2::varchar[])
)
ON CONFLICT (dependency_name, dependency_file) DO UPDATE SET dependency_name = excluded.dependency_name
RETURNING internal_id
`

type BulkInsertDependenciesParams struct {
	DependencyName string   `json:"dependency_name"`
	Column2        []string `json:"column_2"`
}

func (q *Queries) BulkInsertDependencies(ctx context.Context, arg BulkInsertDependenciesParams) ([]int32, error) {
	rows, err := q.query(ctx, q.bulkInsertDependenciesStmt, bulkInsertDependencies, arg.DependencyName, pq.Array(arg.Column2))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var internal_id int32
		if err := rows.Scan(&internal_id); err != nil {
			return nil, err
		}
		items = append(items, internal_id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDependencies = `-- name: GetDependencies :many
SELECT internal_id, dependency_name, dependency_file FROM dependencies LIMIT 10
`

func (q *Queries) GetDependencies(ctx context.Context) ([]Dependency, error) {
	rows, err := q.query(ctx, q.getDependenciesStmt, getDependencies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Dependency
	for rows.Next() {
		var i Dependency
		if err := rows.Scan(&i.InternalID, &i.DependencyName, &i.DependencyFile); err != nil {
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

const getDependenciesByNames = `-- name: GetDependenciesByNames :many
SELECT internal_id, dependency_name, dependency_file FROM dependencies WHERE dependency_name = ANY($1)
`

func (q *Queries) GetDependenciesByNames(ctx context.Context, dependencyName string) ([]Dependency, error) {
	rows, err := q.query(ctx, q.getDependenciesByNamesStmt, getDependenciesByNames, dependencyName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Dependency
	for rows.Next() {
		var i Dependency
		if err := rows.Scan(&i.InternalID, &i.DependencyName, &i.DependencyFile); err != nil {
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

const getDependency = `-- name: GetDependency :one
SELECT internal_id, dependency_name, dependency_file FROM dependencies WHERE dependency_name = $1 AND dependency_file = $2
`

type GetDependencyParams struct {
	DependencyName string `json:"dependency_name"`
	DependencyFile string `json:"dependency_file"`
}

func (q *Queries) GetDependency(ctx context.Context, arg GetDependencyParams) (Dependency, error) {
	row := q.queryRow(ctx, q.getDependencyStmt, getDependency, arg.DependencyName, arg.DependencyFile)
	var i Dependency
	err := row.Scan(&i.InternalID, &i.DependencyName, &i.DependencyFile)
	return i, err
}
