package usersync

import "fmt"

type IdentifyRepoAuthorsBatchTestCase struct {
	title            string
	repoUrl          string
	authorCommitList []AuthorCommitTuple
	authorized       bool
	expectedOutput   map[string]GithubGraphQLCommit
}

const repoUrl = "https://github.com/OpenQDev/OpenQ-Workflows"

func identifyRepoAuthorsBatchTest1() IdentifyRepoAuthorsBatchTestCase {
	const AUTHOR_BATCH = "AUTHOR_BATCH"

	user := GithubGraphQLUser{
		GithubRestID:    93455288,
		GithubGraphqlID: "U_kgDOBZIDuA",
		Login:           "FlacoJones",
		Name:            "AndrewOBrien",
		Email:           "",
		AvatarURL:       "https://avatars.githubusercontent.com/u/93455288?u=fd1fb04b6ff2bf397f8353eafffc3bfb4bd66e84&v=4",
		Company:         "",
		Location:        "",
		Hireable:        false,
		Bio:             "builder at OpenQ",
		Blog:            "",
		TwitterUsername: "",
		Followers: struct {
			TotalCount int `json:"totalCount"`
		}{
			TotalCount: 13,
		},
		Following: struct {
			TotalCount int `json:"totalCount"`
		}{
			TotalCount: 0,
		},
		CreatedAt: "2021-10-30T23:43:10Z",
		UpdatedAt: "2024-07-17T15:21:38Z",
	}

	author := GithubGraphQLAuthor{
		Name:  "FlacoJones",
		Email: "andrew@openq.dev",
		User:  user,
	}

	expectedOutput := make(map[string]GithubGraphQLCommit)
	expectedOutput["commit_0"] = GithubGraphQLCommit{Author: author}
	expectedOutput["commit_1"] = GithubGraphQLCommit{Author: author}
	fmt.Println(expectedOutput["commit_0"])

	return IdentifyRepoAuthorsBatchTestCase{
		title:   AUTHOR_BATCH,
		repoUrl: repoUrl,
		authorCommitList: []AuthorCommitTuple{
			{
				Author:     "abc123@email.com",
				CommitHash: "65062be663cc004b77ca8a3b13255bc5efa42f25",
			},
			{
				Author:     "abc123@email.com",
				CommitHash: "65062be663cc004b77ca8a3b13255bc5efa42f25",
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
