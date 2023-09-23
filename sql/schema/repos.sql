-- +goose Up

CREATE TABLE repos (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	organization TEXT NOT NULL,
	repository TEXT NOT NULL,
	url TEXT NOT NULL UNIQUE
);

-- +goose Down

DROP TABLE repos;