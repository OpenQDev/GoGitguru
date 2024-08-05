-- name: GetUserDependenciesByUpdatedAt :many

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
WHERE rd.url = $1 AND
 gu.internal_id > 0
GROUP BY gu.internal_id, rd.dependency_id, rd.url
) s
LEFT JOIN users_to_dependencies ud ON s.internal_id = ud.user_id AND s.dependency_id = ud.dependency_id

   
GROUP BY s.internal_id, s.dependency_id;


-- name: GetUserDependenciesByUser :many
SELECT ud.first_use_date,
ud.last_use_date,
ud.dependency_id,
ud.user_id,
ud.resync_all
FROM users_to_dependencies ud
WHERE (ud.user_id, ud.dependency_id) IN
(SELECT unnest(sqlc.arg(user_ids)::int[]), unnest(sqlc.arg(dependency_ids)::int[]));


-- name: BulkInsertUserDependencies :exec
INSERT INTO users_to_dependencies (user_id, dependency_id, first_use_date, last_use_date, updated_at) VALUES (  
  unnest(sqlc.arg(user_id)::int[]),  
  unnest(sqlc.arg(dependency_id)::int[]),  
  unnest(sqlc.arg(first_use_date)::bigint[]),  
  unnest(sqlc.arg(last_use_date)::bigint[]),
  sqlc.arg(updated_at)::bigint)

ON CONFLICT (user_id, dependency_id) DO UPDATE
SET last_use_date = excluded.last_use_date,
first_use_date = excluded.first_use_date,
updated_at = excluded.updated_at,
resync_all = false
RETURNING user_id, dependency_id;


-- name: SwitchUsersRelationToSimple :exec
UPDATE users_to_dependencies
SET dependency_id =  (
  SELECT internal_id FROM  dependencies  d
WHERE d.dependency_name  = sqlc.arg(dependency_file)
)
WHERE dependency_id IN (
  SELECT internal_id FROM dependencies  d2
WHERE d2.dependency_name  LIKE sqlc.arg(dependency_file_like)

);

-- name: GetAllUserDependenciesByUser :many
SELECT gu.internal_id, gu.login, ud.first_use_date, ud.last_use_date, d.dependency_file, d.dependency_name  FROM (SELECT internal_id, login FROM github_users WHERE login = $1) gu
JOIN users_to_dependencies ud ON internal_id = ud.user_id
JOIN dependencies d ON d.internal_id = ud.dependency_id;

