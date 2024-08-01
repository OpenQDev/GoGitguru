package reposync

import (
	"database/sql"
	"path/filepath"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type GetObjectsFromCommitListTestCase struct {
	name                       string
	params                     GitLogParams
	dependencyFiles            []string
	commitList                 []*object.Commit
	numberOfCommits            int
	currentDependencies        []database.GetRepoDependenciesByURLRow
	bulkInsertCommitsParams    database.BulkInsertCommitsParams
	bulkInsertDependencyParams database.BatchInsertRepoDependenciesParams
	usersToRepoUrl             UsersToRepoUrl
}

func makeCommitByReference(hashString string, author string, authorEmail string, index int) *object.Commit {

	repoDir := filepath.Join("../mock", organization, repo)
	commitList, err := CreateCommitList(repoDir)
	if err != nil {
		panic(err)
	}
	commit := commitList[index]
	commit.Hash = plumbing.NewHash(hashString)
	commit.Author.Name = author
	commit.Author.Email = authorEmail
	return commit

}

func validGetObjectsFromCommitListTest() GetObjectsFromCommitListTestCase {
	const THREE_COMMITS = "THREE_COMMITS"
	currentDependency := database.GetRepoDependenciesByURLRow{
		DependencyName: "eslint",
		DependencyFile: "package.json",
		FirstUseDate:   sql.NullInt64{Int64: 1620000000, Valid: true},
		LastUseDate:    sql.NullInt64{Int64: 0, Valid: true},
	}
	goodGetObjectsFromCommitListTestCase := GetObjectsFromCommitListTestCase{

		name: THREE_COMMITS,
		params: GitLogParams{
			prefixPath:   "../mock",
			organization: organization,
			repo:         repo,
			repoUrl:      "repoUrl",
		},
		commitList: []*object.Commit{
			makeCommitByReference("09442fceb096a56226fb528368ddf971e776057f", "DRM-Test-User", "150183211+DRM-Test-User@users.noreply.github.com", 0),
			makeCommitByReference("a7ce99317e5347735ec5349f303c7036cd007d94", "DRM-Test-User", "150183211+DRM-Test-User@users.noreply.github.com", 1),
			makeCommitByReference("32f8b288406652840a600e18d562a51661d64d99", "DRM-Test-User", "info@openq.dev", 2),
		},
		numberOfCommits: 2,
		currentDependencies: []database.GetRepoDependenciesByURLRow{
			currentDependency,
		},
		dependencyFiles: []string{
			"package.json",
			"requirements.txt",
			"pom.xml",
			"Pipfile",
			"go.mod",
			"build.gradle",
			"Gemfile",
			"Cargo.toml",
			".cabal",
			"composer.json",

			"hardhat.config",
			"truffle",
			`\/network\/`,
			`\/deployments\/`,
			"foundry.toml",
		},
		bulkInsertDependencyParams: database.BatchInsertRepoDependenciesParams{
			Url:             "repoUrl",
			Filenames:       []string{"package.json", "package.json"},
			Dependencynames: []string{"eslint", "web3"},
			Firstusedates:   []int64{1620000000, 1699383684},
			Lastusedates:    []int64{1699384512, 0},
			UpdatedAt: sql.NullInt64{
				Int64: 1609459200,
				Valid: true,
			},
		},
		bulkInsertCommitsParams: database.BulkInsertCommitsParams{
			Commithashes: []string{
				"a7ce99317e5347735ec5349f303c7036cd007d94",
				"32f8b288406652840a600e18d562a51661d64d99",
			},
			Authors: []string{
				"DRM-Test-User",
				"DRM-Test-User",
			},
			Authoremails: []string{
				"150183211+DRM-Test-User@users.noreply.github.com",
				"info@openq.dev",
			},
			Authordates:    []int64{1699383684, 1699384512},
			Committerdates: []int64{1699383684, 1699384512},
			Messages:       []string{"Create package.json", "Create BigFile.json"},
			Fileschanged:   []int32{0, 0},
			Repourl: sql.NullString{
				String: "repoUrl",
				Valid:  true,
			},
		},
		usersToRepoUrl: UsersToRepoUrl{
			AuthorEmails:     []string{"150183211+DRM-Test-User@users.noreply.github.com", "info@openq.dev"},
			FirstCommitDates: []int64{1699383684, 1699384512},
			LastCommitDates:  []int64{1699383684, 1699384512},
		},
	}

	return goodGetObjectsFromCommitListTestCase
}

func GetObjectsFromCommitListTestCases() []GetObjectsFromCommitListTestCase {
	return []GetObjectsFromCommitListTestCase{
		validGetObjectsFromCommitListTest(),
	}

}
