-- +goose Up

CREATE TABLE file_patterns (
    id SERIAL PRIMARY KEY,
    pattern TEXT NOT NULL,
    updated_at INTEGER NOT NULL,
    creator TEXT NOT NULL
);