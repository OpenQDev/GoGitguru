-- +goose Up

CREATE TABLE github_repos (
    internal_id SERIAL PRIMARY KEY,
    github_rest_id INT NOT NULL,
    github_graphql_id VARCHAR(60) NOT NULL,
    url VARCHAR(150) NOT NULL,
    name VARCHAR(120) NOT NULL,
    full_name VARCHAR(150) NOT NULL,
    private BOOLEAN,
    owner_login VARCHAR(120) NOT NULL,
    owner_avatar_url VARCHAR(150),
    description TEXT,
    homepage VARCHAR(150),
    fork BOOLEAN,
    forks_count INT,
    archived BOOLEAN,
    disabled BOOLEAN,
    license VARCHAR(120),
    language VARCHAR(60),
    stargazers_count INT,
    watchers_count INT,
    open_issues_count INT,
    has_issues BOOLEAN,
    has_discussions BOOLEAN,
    has_projects BOOLEAN,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    pushed_at TIMESTAMP,
    visibility VARCHAR(30),
    size INT,
    default_branch VARCHAR(60),
    UNIQUE(github_rest_id),
    UNIQUE(github_graphql_id)
);

-- +goose Down

DROP TABLE github_repos;