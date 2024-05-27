-- name: BatchInsertRepoDependencies :exec
INSERT INTO repos_to_dependencies (url, dependency_id, first_use_data, last_use_data) VALUES (  
  $1,  
  unnest($2::int[]),  
  unnest($3::bigint[]),  
  unnest($4::bigint[]) 
)
ON CONFLICT (url, dependency_id) DO UPDATE 
SET last_use_data = excluded.last_use_data, 
first_use_data = excluded.first_use_data;

-- name: InitializeRepoDependencies :exec
INSERT INTO repos_to_dependencies (url, dependency_id) VALUES (  
 $1,  
  unnest($2::int[]) 
)
ON CONFLICT (url, dependency_id) DO NOTHING;

-- name: GetRepoDependencies :many
SELECT 
d.dependency_name,
rd.first_use_data,
rd.last_use_data
FROM dependencies d
LEFT JOIN repos_to_dependencies rd ON d.internal_id = rd.dependency_id
WHERE d.dependency_name = $1 AND rd.url = $2 AND d.dependency_file = ANY($3::text[]);





