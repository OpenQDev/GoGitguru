package reposync

import (
	"fmt"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func GetObjectsFromCommitList(params GitLogParams, commitList []*object.Commit, numberOfCommits int, currentDependencies []database.GetRepoDependenciesByURLRow) (database.BatchInsertRepoDependenciesParams, database.BulkInsertCommitsParams, error) {
	// sync this from the db
	dependencyHistoryObject := database.BatchInsertRepoDependenciesParams{
		Url:             params.repoUrl,
		Firstusedates:   []int64{},
		Lastusedates:    []int64{},
		Dependencynames: []string{},
		Filenames:       []string{},
	}
	for _, dep := range currentDependencies {
		dependencyHistoryObject.Dependencynames = append(dependencyHistoryObject.Dependencynames, dep.DependencyName)
		dependencyHistoryObject.Filenames = append(dependencyHistoryObject.Filenames, dep.DependencyFile)
		dependencyHistoryObject.Firstusedates = append(dependencyHistoryObject.Firstusedates, dep.FirstUseData.Int64)
		dependencyHistoryObject.Lastusedates = append(dependencyHistoryObject.Lastusedates, dep.LastUseData.Int64)
	}
	commitWindow := GetCommitWindow(len(commitList))

	commitObject := database.BulkInsertCommitsParams{
		Commithashes:   []string{},
		Authors:        []string{},
		Authoremails:   []string{},
		Authordates:    []int64{},
		Committerdates: []int64{},
		Messages:       []string{},
		Fileschanged:   []int32{},
	}
	var err error
	// start from first commit that hasn't been synced
	for commitCount, commit := range commitList {

		fmt.Printf("commit date: %v\n", commit.Committer.When.Unix())
		if commitCount >= numberOfCommits {
			println("commit count is greater than or equal to number of commits check", commitCount, numberOfCommits)
			break
		}
		if commitCount < numberOfCommits {
			if commitCount%commitWindow == 0 {
				err = CheckCommitForDependencies(commit, &dependencyHistoryObject)
				if err != nil {
					return dependencyHistoryObject, commitObject, err
				}
			}
			AddCommitToCommitObject(commit, &commitObject, commitCount)
		}

	}
	c := commitList[len(commitList)-1]
	// always check last commit last
	fmt.Printf("Commit number %d: %s\n", len(commitList)-1, c.Hash)
	err = CheckCommitForDependencies(c, &dependencyHistoryObject)

	return dependencyHistoryObject, commitObject, err
}
