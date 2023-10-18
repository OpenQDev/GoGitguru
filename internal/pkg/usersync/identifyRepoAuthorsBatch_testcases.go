package usersync

import (
	"main/internal/pkg/githubGraphQL"
)

type IdentifyRepoAuthorsBatchTestCase struct {
	title            string
	repoUrl          string
	authorCommitList []AuthorCommitTuple
	authorized       bool
	expectedOutput   *map[string]githubGraphQL.Commit
}

const repoUrl = "OpenQ-Workflows"

func identifyRepoAuthorsBatchTest1() IdentifyRepoAuthorsBatchTestCase {
	const AUTHOR_BATCH = "AUTHOR_BATCH"

	author := githubGraphQL.Author{
		Name:  "FlacoJones",
		Email: "andrew@openq.dev",
		User: struct {
			GithubRestID    int    `json:"github_rest_id"`
			GithubGraphqlID string `json:"github_graphql_id"`
			Login           string `json:"login"`
			Name            string `json:"name"`
			Email           string `json:"email"`
			AvatarURL       string `json:"avatar_url"`
			Company         string `json:"company"`
			Location        string `json:"location"`
			Hireable        bool   `json:"hireable"`
			Bio             string `json:"bio"`
			Blog            string `json:"blog"`
			TwitterUsername string `json:"twitter_username"`
			Followers       struct {
				TotalCount int `json:"totalCount"`
			} `json:"followers"`
			Following struct {
				TotalCount int `json:"totalCount"`
			} `json:"following"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		}{
			GithubRestID:    93455288,
			GithubGraphqlID: "U_kgDOBZIDuA",
			Login:           "FlacoJones",
			Name:            "AndrewOBrien",
			Email:           "",
			AvatarURL:       "https://avatars.githubusercontent.com/u/93455288?u=fd1fb04b6ff2bf397f8353eafffc3bfb4bd66e84\u0026v=4",
			Company:         "",
			Location:        "",
			Hireable:        false,
			Bio:             "builder at OpenQ",
			Blog:            "",
			TwitterUsername: "",
			Followers: struct {
				TotalCount int `json:"totalCount"`
			}{
				TotalCount: 12,
			},
			Following: struct {
				TotalCount int `json:"totalCount"`
			}{
				TotalCount: 0,
			},
			CreatedAt: "2021-10-30T23:43:10Z",
			UpdatedAt: "2023-10-10T15:52:33Z",
		},
	}

	expectedOutput := make(map[string]githubGraphQL.Commit)
	expectedOutput["commitHash"] = githubGraphQL.Commit{Author: author}

	return IdentifyRepoAuthorsBatchTestCase{
		title:   AUTHOR_BATCH,
		repoUrl: repoUrl,
		authorCommitList: []AuthorCommitTuple{
			AuthorCommitTuple{
				Author:     "abc123@email.com",
				CommitHash: "commitHash",
			},
		},
		authorized:     true,
		expectedOutput: &expectedOutput,
	}
}

func IdentifyRepoAuthorsBatchTestCases() []IdentifyRepoAuthorsBatchTestCase {
	return []IdentifyRepoAuthorsBatchTestCase{
		identifyRepoAuthorsBatchTest1(),
	}
}
