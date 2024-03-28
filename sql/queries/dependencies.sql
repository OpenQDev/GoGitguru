-- name: InsertDependency :one

INSERT INTO dependencies (
    dependency_name,
    dependency_file
) VALUES (
    $1, $2
)
ON CONFLICT (dependency_name, dependency_file) DO NOTHING
RETURNING *;

-- name: GetDependency :one
SELECT * FROM dependencies WHERE dependency_name = $1 AND dependency_file = $2;

-- name: GetDependencies :many
SELECT * FROM dependencies LIMIT 10;

-- name: GetDependenciesByNames :many
SELECT * FROM dependencies WHERE dependency_name = ANY($1);


