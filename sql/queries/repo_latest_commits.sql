-- name: GetRepoLatestCommit :one
SELECT * FROM repo_latest_commit WHERE repo_name = $1;
