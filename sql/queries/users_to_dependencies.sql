-- name: BatchInsertUserDependencies :exec
INSERT INTO user_to_dependencies (user_id, dependency_id, first_use_data, last_use_data) VALUES (  
  $1,  
  unnest($2::int[]),  
  unnest($3::bigint[]),  
  unnest($4::bigint[]) 
)
ON CONFLICT (user_id, dependency_id) DO UPDATE 
SET last_use_data = excluded.last_use_data, 
first_use_data = excluded.first_use_data;

-- name: InitializeUserDependencies :exec
INSERT INTO user_to_dependencies ( user_id, dependency_id)
SELECT  internal_id, unnest($2::int[])
FROM github_users
WHERE login = $1
ON CONFLICT (user_id, dependency_id) DO NOTHING;

-- name: GetUserDependenciesByUpdatedAt :many

WITH computed_dates AS (
    SELECT
        gu.login,
        du.user_id,
        du.dependency_id,
        CASE
            WHEN subquery.last_use_date = 0 THEN NULL
            ELSE LEAST(subquery.last_commit_date, subquery.last_use_date)
        END AS last_use_date,
        CASE
            WHEN subquery.first_use_date = 0 THEN NULL
            ELSE GREATEST(subquery.first_commit_date, subquery.first_use_date)
        END AS first_use_date
    FROM
        user_to_dependencies du
    LEFT JOIN
        dependencies d ON du.dependency_id = d.internal_id
    LEFT JOIN
        github_users gu ON du.user_id = gu.internal_id
    LEFT JOIN
        github_user_rest_id_author_emails gure ON gu.github_rest_id = gure.rest_id
    LEFT JOIN (
        SELECT
            rd.dependency_id,
            rd.url,
            MIN(c1.author_date) AS first_commit_date,
            MAX(c1.author_date) AS last_commit_date,
            MIN(rd.first_use_data) AS first_use_date,
            MAX(rd.last_use_data) AS last_use_date
        FROM
            repos_to_dependencies rd
        LEFT JOIN
            commits c1 ON c1.repo_url = rd.url
        GROUP BY
            rd.dependency_id,
            rd.url
    ) subquery ON subquery.dependency_id = d.internal_id
    LEFT JOIN
        commits first_commit ON first_commit.author_date = subquery.first_commit_date
    LEFT JOIN
        commits last_commit ON last_commit.author_date = subquery.last_commit_date
    WHERE du.updated_at IS NULL OR du.updated_at < $1
)
SELECT
    login,
    user_id,
    dependency_id,
    MIN(first_use_date) AS earliest_first_use_date,
    MAX(last_use_date) AS latest_last_use_date
FROM
    computed_dates
GROUP BY
    login,
    user_id,
    dependency_id;

-- name: BulkInsertUserDependencies :many
INSERT INTO user_to_dependencies (user_id, dependency_id, first_use_data, last_use_data, updated_at) VALUES (  
  unnest($1::int[]),  
  unnest($2::int[]),  
  unnest($3::bigint[]),  
  unnest($4::bigint[]),
  $5
)
ON CONFLICT (user_id, dependency_id) DO UPDATE
SET last_use_data = excluded.last_use_data,
first_use_data = excluded.first_use_data,
updated_at = excluded.updated_at
RETURNING user_id, dependency_id, updated_at;
