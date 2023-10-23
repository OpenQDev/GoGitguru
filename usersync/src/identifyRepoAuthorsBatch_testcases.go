package usersync

import (
	"github.com/OpenQDev/GoGitguru/util/githubGraphQL"
)

type IdentifyRepoAuthorsBatchTestCase struct {
	title            string
	repoUrl          string
	authorCommitList []AuthorCommitTuple
	authorized       bool
	expectedOutput   map[string]githubGraphQL.GithubGraphQLCommit
}

const repoUrl = "https://github.com/OpenQDev/OpenQ-Workflows"

func identifyRepoAuthorsBatchTest1() IdentifyRepoAuthorsBatchTestCase {
	const AUTHOR_BATCH = "AUTHOR_BATCH"

	user := githubGraphQL.GithubGraphQLUser{
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
	}

	author := githubGraphQL.GithubGraphQLAuthor{
		Name:  "FlacoJones",
		Email: "andrew@openq.dev",
		User:  user,
	}

	expectedOutput := make(map[string]githubGraphQL.GithubGraphQLCommit)
	expectedOutput["commit_0"] = githubGraphQL.GithubGraphQLCommit{Author: author}
	expectedOutput["commit_1"] = githubGraphQL.GithubGraphQLCommit{Author: author}

	return IdentifyRepoAuthorsBatchTestCase{
		title:   AUTHOR_BATCH,
		repoUrl: repoUrl,
		authorCommitList: []AuthorCommitTuple{
			{
				Author:     "abc123@email.com",
				CommitHash: "commitHash",
			},
		},
		authorized:     true,
		expectedOutput: expectedOutput,
	}
}

func IdentifyRepoAuthorsBatchTestCases() []IdentifyRepoAuthorsBatchTestCase {
	return []IdentifyRepoAuthorsBatchTestCase{
		identifyRepoAuthorsBatchTest1(),
	}
}
