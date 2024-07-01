-- name: GetAllFilePatterns :many
SELECT * FROM file_patterns
LIMIT 50;

-- name: BulkUpsertFilePatterns :exec
INSERT INTO file_patterns (pattern, updated_at, creator) VALUES (unnest(sqlc.arg(patterns)::text[]), sqlc.arg(updated_at)::bigint, sqlc.arg(creator)::text)
ON CONFLICT DO NOTHING;