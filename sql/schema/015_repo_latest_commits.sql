-- +goose Up

CREATE TABLE repo_latest_commit (
		repo_name TEXT NOT NULL PRIMARY KEY,
		last_push_event_time TIMESTAMPTZ NOT NULL
)

-- +goose Down
DROP TABLE repo_latest_commit;