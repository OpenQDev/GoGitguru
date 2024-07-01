-- name: BulkInsertDependencies :many

INSERT INTO dependencies (
    dependency_name,
    dependency_file
) VALUES (
    $1, unnest($2::varchar[])
)
ON CONFLICT (dependency_name, dependency_file) DO UPDATE SET dependency_name = excluded.dependency_name
RETURNING internal_id;

-- name: GetDependency :one
SELECT * FROM dependencies WHERE dependency_name = $1 AND dependency_file = $2;

-- name: GetDependencies :many
SELECT d.dependency_name,
d.dependency_file,
d.internal_id,
rd.updated_at
 FROM repos_to_dependencies rd
LEFT JOIN dependencies d ON rd.dependency_id = d.internal_id
WHERE rd.url = $1;

-- name: GetDependenciesByNames :many
SELECT * FROM dependencies WHERE dependency_name =ANY($1::text[]);


-- name: GetDependenciesByFiles :many

SELECT dependency_name
FROM dependencies
WHERE dependency_file LIKE sqlc.arg('dependency_file_like')
  AND dependency_file != sqlc.arg('dependency_file')
GROUP BY dependency_name; 

-- name: UpsertMissingDependencies :exec

INSERT INTO dependencies (
    dependency_name,
    dependency_file
) VALUES (
    unnest(sqlc.arg(dependency_names)::text[]), $1
)
ON CONFLICT DO NOTHING
RETURNING internal_id;

-- name: DeleteUnusedDependencies :exec
DELETE FROM dependencies WHERE dependency_file LIKE sqlc.arg('dependency_file_like') AND dependency_file != sqlc.arg('dependency_file');