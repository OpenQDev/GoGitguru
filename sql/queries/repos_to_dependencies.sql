-- name: BatchInsertRepoDependencies :exec
INSERT INTO repos_to_dependencies (url, dependency_id, first_use_data, last_use_data) VALUES (  
  unnest($1::varchar[]),  
  unnest($2::int[]),  
  unnest($3::bigint[]),  
  unnest($4::bigint[]) 
)
ON CONFLICT (url, dependency_id) DO UPDATE 
SET last_use_data = excluded.last_use_data, 
first_use_data = excluded.first_use_data;


