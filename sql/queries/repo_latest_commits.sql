-- name: UpsertRepoLatestCommit :exec

INSERT INTO repo_latest_commit (repo_name, last_push_event_time)
VALUES ($1, $2)
ON CONFLICT (repo_name) 
DO UPDATE SET 
    last_push_event_time = EXCLUDED.last_push_event_time;
