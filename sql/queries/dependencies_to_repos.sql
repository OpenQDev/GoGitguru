-- name: BulkInsertRepoDependencyInfo :exec
INSERT INTO dependencies_to_repos (github_repo_internal_id, dependency_name, first_commit_date, date_added, date_removed) VALUES (  
	unnest($1::int[]),  
	unnest($2::varchar[]),
	unnest($3::bigint[]),
	unnest($4::bigint[]),
	unnest($5::bigint[])
)
ON CONFLICT (dependency_name) DO NOTHING
RETURNING *;

-- name: QueryBulkRepoDependencyInfo :many
SELECT * FROM dependencies_to_repos WHERE github_repo_internal_id = ANY($1::int[]) AND dependency_name = $2	;
