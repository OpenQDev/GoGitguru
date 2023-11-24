-- name: InsertDependencies :one

INSERT INTO dependencies (
	dependency_name,
	dependency_files
) VALUES (
    $1, $2
)
ON CONFLICT (dependency_name) DO NOTHING
RETURNING *;
