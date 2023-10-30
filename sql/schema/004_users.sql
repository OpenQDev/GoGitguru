-- +goose Up

CREATE TABLE github_users (
    internal_id SERIAL PRIMARY KEY,
    github_rest_id INT NOT NULL,
    github_graphql_id VARCHAR(60) NOT NULL,
    login VARCHAR(120) NOT NULL,
    name VARCHAR(120),
    email VARCHAR(255),
    avatar_url VARCHAR(150),
    company VARCHAR(120),
    location VARCHAR(120),
    bio TEXT,
    blog VARCHAR(150),
    hireable BOOLEAN,
    twitter_username VARCHAR(120),
    followers INT,
    following INT,
    type VARCHAR(30) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    UNIQUE(github_rest_id),
    UNIQUE(github_graphql_id)
);

-- +goose Down

DROP TABLE github_users;


				author.User.GithubRestID,
				author.User.GithubGraphqlID,
				author.User.Login,
				author.User.Name,
				author.User.Email,
				author.User.AvatarURL,
				author.User.Company,
				author.User.Location,

								author.User.GithubRestID,
				author.User.GithubGraphqlID,
				author.User.Login,
				author.User.Name,
				author.User.Email,
				author.User.AvatarURL,
				author.User.Company,
				author.User.Location,