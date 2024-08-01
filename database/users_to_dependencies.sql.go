// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: users_to_dependencies.sql

package database

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
)

const bulkInsertUserDependencies = `-- name: BulkInsertUserDependencies :exec
INSERT INTO users_to_dependencies (user_id, dependency_id, first_use_date, last_use_date, updated_at) VALUES (  
  unnest($1::int[]),  
  unnest($2::int[]),  
  unnest($3::bigint[]),  
  unnest($4::bigint[]),
  $5::bigint)

ON CONFLICT (user_id, dependency_id) DO UPDATE
SET last_use_date = excluded.last_use_date,
first_use_date = excluded.first_use_date,
resync_all = false
RETURNING user_id, dependency_id
`

type BulkInsertUserDependenciesParams struct {
	UserID       []int32 `json:"user_id"`
	DependencyID []int32 `json:"dependency_id"`
	FirstUseDate []int64 `json:"first_use_date"`
	LastUseDate  []int64 `json:"last_use_date"`
	UpdatedAt    int64   `json:"updated_at"`
}

func (q *Queries) BulkInsertUserDependencies(ctx context.Context, arg BulkInsertUserDependenciesParams) error {
	_, err := q.exec(ctx, q.bulkInsertUserDependenciesStmt, bulkInsertUserDependencies,
		pq.Array(arg.UserID),
		pq.Array(arg.DependencyID),
		pq.Array(arg.FirstUseDate),
		pq.Array(arg.LastUseDate),
		arg.UpdatedAt,
	)
	return err
}

const getAllUserDependenciesByUser = `-- name: GetAllUserDependenciesByUser :many
SELECT gu.internal_id, gu.login, ud.first_use_date, ud.last_use_date, d.dependency_file, d.dependency_name  FROM (SELECT internal_id, login FROM github_users WHERE login = $1) gu
JOIN users_to_dependencies ud ON internal_id = ud.user_id
JOIN dependencies d ON d.internal_id = ud.dependency_id
`

type GetAllUserDependenciesByUserRow struct {
	InternalID     int32         `json:"internal_id"`
	Login          string        `json:"login"`
	FirstUseDate   sql.NullInt64 `json:"first_use_date"`
	LastUseDate    sql.NullInt64 `json:"last_use_date"`
	DependencyFile string        `json:"dependency_file"`
	DependencyName string        `json:"dependency_name"`
}

func (q *Queries) GetAllUserDependenciesByUser(ctx context.Context, login string) ([]GetAllUserDependenciesByUserRow, error) {
	rows, err := q.query(ctx, q.getAllUserDependenciesByUserStmt, getAllUserDependenciesByUser, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUserDependenciesByUserRow
	for rows.Next() {
		var i GetAllUserDependenciesByUserRow
		if err := rows.Scan(
			&i.InternalID,
			&i.Login,
			&i.FirstUseDate,
			&i.LastUseDate,
			&i.DependencyFile,
			&i.DependencyName,
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

const getUserDependenciesByUpdatedAt = `-- name: GetUserDependenciesByUpdatedAt :many

SELECT s.internal_id as user_id, 
MIN(s.first_use_date_result) as first_use_date,
MAX(s.last_use_date_result) as last_use_date,
 s.dependency_id

 FROM (

SELECT GREATEST(rd.first_use_date , MIN(c.committer_date)) as first_use_date_result, 
	 LEAST(rd.last_use_date, MAX(c.committer_date)) as last_use_date_result, 
	 rd.url, 
	 gu.internal_id, rd.dependency_id
FROM repos_to_dependencies rd
LEFT JOIN commits c ON c.repo_url = rd.url
LEFT JOIN github_user_rest_id_author_emails guriae ON guriae.email = c.author_email
LEFT JOIN github_users gu ON gu.github_rest_id = guriae.rest_id
WHERE (rd.updated_at > $1 OR rd.updated_at IS NULL  ) AND
 gu.internal_id > 0
GROUP BY gu.internal_id, rd.dependency_id, rd.url
) s
LEFT JOIN users_to_dependencies ud ON s.internal_id = ud.user_id AND s.dependency_id = ud.dependency_id

   
GROUP BY s.internal_id, s.dependency_id
`

type GetUserDependenciesByUpdatedAtRow struct {
	UserID       sql.NullInt32 `json:"user_id"`
	FirstUseDate interface{}   `json:"first_use_date"`
	LastUseDate  interface{}   `json:"last_use_date"`
	DependencyID int32         `json:"dependency_id"`
}

func (q *Queries) GetUserDependenciesByUpdatedAt(ctx context.Context, updatedAt sql.NullInt64) ([]GetUserDependenciesByUpdatedAtRow, error) {
	rows, err := q.query(ctx, q.getUserDependenciesByUpdatedAtStmt, getUserDependenciesByUpdatedAt, updatedAt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserDependenciesByUpdatedAtRow
	for rows.Next() {
		var i GetUserDependenciesByUpdatedAtRow
		if err := rows.Scan(
			&i.UserID,
			&i.FirstUseDate,
			&i.LastUseDate,
			&i.DependencyID,
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

const getUserDependenciesByUser = `-- name: GetUserDependenciesByUser :many
SELECT ud.first_use_date,
ud.last_use_date,
ud.dependency_id,
ud.user_id
FROM users_to_dependencies ud
WHERE (ud.user_id, ud.dependency_id) IN
(SELECT unnest($1::int[]), unnest($2::int[]))
`

type GetUserDependenciesByUserParams struct {
	UserIds       []int32 `json:"user_ids"`
	DependencyIds []int32 `json:"dependency_ids"`
}

type GetUserDependenciesByUserRow struct {
	FirstUseDate sql.NullInt64 `json:"first_use_date"`
	LastUseDate  sql.NullInt64 `json:"last_use_date"`
	DependencyID int32         `json:"dependency_id"`
	UserID       int32         `json:"user_id"`
}

func (q *Queries) GetUserDependenciesByUser(ctx context.Context, arg GetUserDependenciesByUserParams) ([]GetUserDependenciesByUserRow, error) {
	rows, err := q.query(ctx, q.getUserDependenciesByUserStmt, getUserDependenciesByUser, pq.Array(arg.UserIds), pq.Array(arg.DependencyIds))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserDependenciesByUserRow
	for rows.Next() {
		var i GetUserDependenciesByUserRow
		if err := rows.Scan(
			&i.FirstUseDate,
			&i.LastUseDate,
			&i.DependencyID,
			&i.UserID,
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

const switchUsersRelationToSimple = `-- name: SwitchUsersRelationToSimple :exec
UPDATE users_to_dependencies
SET dependency_id =  (
  SELECT internal_id FROM  dependencies  d
WHERE d.dependency_name  = $1
)
WHERE dependency_id IN (
  SELECT internal_id FROM dependencies  d2
WHERE d2.dependency_name  LIKE $2

)
`

type SwitchUsersRelationToSimpleParams struct {
	DependencyFile     string `json:"dependency_file"`
	DependencyFileLike string `json:"dependency_file_like"`
}

func (q *Queries) SwitchUsersRelationToSimple(ctx context.Context, arg SwitchUsersRelationToSimpleParams) error {
	_, err := q.exec(ctx, q.switchUsersRelationToSimpleStmt, switchUsersRelationToSimple, arg.DependencyFile, arg.DependencyFileLike)
	return err
}
