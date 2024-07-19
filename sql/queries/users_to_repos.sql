
-- name: UpsertRepoToUserById :exec
INSERT INTO users_to_repo_urls (url, user_id, first_commit_date, last_commit_date)
SELECT $1, unnest(sqlc.arg(internal_ids)::int[]), unnest(sqlc.arg(first_commit_dates)::bigint[]), unnest(sqlc.arg(last_commit_dates)::bigint[])
ON CONFLICT (url, user_id) DO UPDATE
    SET last_commit_date = GREATEST(users_to_repo_urls.last_commit_date, EXCLUDED.last_commit_date),
    first_commit_date = LEAST(users_to_repo_urls.first_commit_date, EXCLUDED.first_commit_date)
    RETURNING url, user_id;
