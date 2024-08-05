-- name: BatchInsertRepoDependencies :many
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
    s.updated_at::bigint as updated_at,
    d.internal_id AS dependency_id,
    s.url,
    s.firstUseDates AS first_use_date,
    s.lastUseDates AS last_use_date
  FROM (
    SELECT
  $1 as updated_at,
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
  last_use_date = EXCLUDED.last_use_date,
  updated_at = EXCLUDED.updated_at
RETURNING url;



-- name: InitializeRepoDependencies :exec
INSERT INTO repos_to_dependencies (url, dependency_id) VALUES (  
 $1,  
  unnest($2::int[]) 
)
ON CONFLICT (url, dependency_id) DO NOTHING;

-- name: GetRepoDependencies :many
SELECT 
d.dependency_name,
rd.first_use_date,
rd.last_use_date,
rd.url
FROM dependencies d
LEFT JOIN repos_to_dependencies rd ON d.internal_id = rd.dependency_id
WHERE d.dependency_name = $1 AND rd.url =  ANY($2::text[]) AND d.dependency_file = ANY($3::text[]);


-- name: GetRepoDependenciesByURL :many
SELECT 
d.dependency_name,
d.dependency_file,
rd.first_use_date,
rd.last_use_date
FROM dependencies d
LEFT JOIN repos_to_dependencies rd ON d.internal_id = rd.dependency_id
WHERE rd.url = $1;







-- name: SwitchReposRelationToSimple :exec
UPDATE users_to_dependencies
SET dependency_id =  (
  SELECT internal_id FROM dependencies d1
WHERE d1.dependency_name  = sqlc.arg(dependency_file)
)
WHERE dependency_id IN (
  SELECT internal_id FROM dependencies d2
WHERE d2.dependency_name  LIKE sqlc.arg(dependency_file_like)

) ;