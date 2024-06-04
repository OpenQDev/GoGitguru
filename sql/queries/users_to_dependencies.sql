-- name: BatchInsertUserDependencies :exec
INSERT INTO user_to_dependencies (user_id, dependency_id, first_use_date, last_use_date) VALUES (  
  $1,  
  unnest($2::int[]),  
  unnest($3::bigint[]),  
  unnest($4::bigint[]) 
)
ON CONFLICT (user_id, dependency_id) DO UPDATE 
SET last_use_date = excluded.last_use_date, 
first_use_date = excluded.first_use_date;

-- name: InitializeUserDependencies :exec
INSERT INTO user_to_dependencies ( user_id, dependency_id)
SELECT  internal_id, unnest($2::int[])
FROM github_users
WHERE login = $1
ON CONFLICT (user_id, dependency_id) DO NOTHING;

-- name: GetUserDependenciesByUpdatedAt :many

SELECT s.internal_id as user_id, MIN(first_use_date) as first_use_date,
CASE
            WHEN MIN(last_use_date) = 0 THEN NULL
            ELSE MAX(last_use_date)
        END AS last_use_date,
        s.dependency_id

 FROM (

SELECT GREATEST(first_use_date , MIN(c.committer_date)) as first_use_date_result,  LEAST(last_use_date, MAX(c.committer_date)) as last_use_date_result, rd.url, gu.internal_id, rd.dependency_id
FROM repos_to_dependencies rd
LEFT JOIN commits c ON c.repo_url = rd.url
LEFT JOIN github_user_rest_id_author_emails guriae ON guriae.email = c.author_email
LEFT JOIN github_users gu ON gu.github_rest_id = guriae.rest_id
GROUP BY gu.internal_id, rd.dependency_id, rd.url
) s
LEFT JOIN user_to_dependencies ud ON s.internal_id = ud.user_id AND s.dependency_id = ud.dependency_id
WHERE first_use_date_result <> ud.first_use_date or last_use_date_result <> ud.last_use_date
GROUP BY s.internal_id, s.dependency_id;


-- name: BulkInsertUserDependencies :many
INSERT INTO user_to_dependencies (user_id, dependency_id, first_use_date, last_use_date) VALUES (  
  unnest($1::int[]),  
  unnest($2::int[]),  
  unnest($3::bigint[]),  
  unnest($4::bigint[])
)
ON CONFLICT (user_id, dependency_id) DO UPDATE
SET last_use_date = excluded.last_use_date,
first_use_date = excluded.first_use_date
RETURNING user_id, dependency_id;
