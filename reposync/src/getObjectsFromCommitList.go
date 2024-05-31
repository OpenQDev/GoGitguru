package reposync

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetObjectsFromCommitList(params GitLogParams, dependencies []DependencyWithUpdatedTime, commitList []*object.Commit, numberOfCommits NumberOfCommits) (database.BatchInsertRepoDependenciesParams, database.BulkInsertCommitsParams, error) {

	// always check last commit
	c := commitList[len(commitList)-1]
	fmt.Printf("Commit number %d: %s\n", len(commitList)-1, c.Hash)

	dependencyHistoryObject := database.BatchInsertRepoDependenciesParams{
		Url: params.repoUrl,
		Updatedat: sql.NullInt64{
			Int64: time.Now().Unix(),
			Valid: true,
		},
		Dependencyids: make([]int32, 0),
		Firstusedates: make([]int64, 0),
		Lastusedates:  make([]int64, 0),
	}
	dependencyHistoryObject, err := CheckCommitForDependencies(c, dependencyHistoryObject, dependencies, numberOfCommits.ToCheckByDependency, 0)
	commitWindow := getCommitWindow(len(commitList))

	commitObject := database.BulkInsertCommitsParams{
		Commithashes:   make([]string, numberOfCommits.ToSync),
		Authors:        make([]string, numberOfCommits.ToSync),
		Authoremails:   make([]string, numberOfCommits.ToSync),
		Authordates:    make([]int64, numberOfCommits.ToSync),
		Committerdates: make([]int64, numberOfCommits.ToSync),
		Messages:       make([]string, numberOfCommits.ToSync),
		Fileschanged:   make([]int32, numberOfCommits.ToSync),
	}
	fmt.Println("assemble arrays for and check for dependencies", params.repoUrl)
	for commitCount, commit := range commitList {
		if commitCount >= numberOfCommits.GreatestToCheckForDependency && commitCount >= numberOfCommits.ToSync {
			println("commit count is greater than or equal to number of commits check", commitCount, numberOfCommits.GreatestToCheckForDependency, numberOfCommits.ToSync)
			break
		}
		if commitCount%commitWindow == 0 && numberOfCommits.GreatestToCheckForDependency < commitCount {
			dependencyHistoryObject, err = CheckCommitForDependencies(commit, dependencyHistoryObject, dependencies, numberOfCommits.ToCheckByDependency, commitCount)
		}
		if commitCount < numberOfCommits.ToSync {

			AddCommitToCommitObject(commit, &commitObject, commitCount)
		}

	}

	return dependencyHistoryObject, commitObject, err
}
