-- name: BatchInsertRepoDependencies :exec
WITH new_dependencies AS (
  INSERT INTO dependencies (dependency_file, dependency_name)
  SELECT
    DISTINCT dependency_file,
    dependency_name
  FROM (
    SELECT
      unnest(sqlc.arg(filenames)::text[]) AS dependency_file,
      unnest(sqlc.arg(dependencyNames)::text[]) AS dependency_name
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
    d.internal_id AS dependency_id,
    s.url,
    s.firstUseDates AS first_use_data,
    s.lastUseDates AS last_use_data
  FROM (
    SELECT
      sqlc.arg(url) AS url,
      unnest(sqlc.arg(filenames)::text[]) AS dependency_file,
      unnest(sqlc.arg(dependencyNames)::text[]) AS dependency_name,
unnest(COALESCE(sqlc.arg(firstUseDates)::bigint[], ARRAY[]::bigint[]))AS firstUseDates,
unnest(COALESCE(sqlc.arg(lastUseDates)::bigint[], ARRAY[]::bigint[])) AS lastUseDates
  ) s
  JOIN all_dependencies d ON d.dependency_file = s.dependency_file AND d.dependency_name = s.dependency_name
)
INSERT INTO repos_to_dependencies (
  url, 
  dependency_id, 
  first_use_data,
  last_use_data
) 
SELECT DISTINCT
  url, 
  dependency_id,  
  first_use_data,
  last_use_data
FROM
  dependency_ids
ON CONFLICT (url, dependency_id) DO UPDATE 
SET 
  first_use_data = EXCLUDED.first_use_data,
  last_use_data = EXCLUDED.last_use_data;



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


-- name: GetRepoDependenciesByURL :many
SELECT 
d.dependency_name,
d.dependency_file,
rd.first_use_data,
rd.last_use_data
FROM dependencies d
LEFT JOIN repos_to_dependencies rd ON d.internal_id = rd.dependency_id
WHERE rd.url = $1;







