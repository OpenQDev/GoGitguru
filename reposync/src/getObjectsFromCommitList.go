package reposync

import (
	"database/sql"
	"path/filepath"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/lib"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type UsersToRepoUrl struct {
	AuthorEmails     []string
	FirstCommitDates []int64
	LastCommitDates  []int64
}

func GetObjectsFromCommitList(params GitLogParams, commitList []*object.Commit, numberOfCommits int, currentDependencies []database.GetRepoDependenciesByURLRow, dependencyFiles []string) (database.BatchInsertRepoDependenciesParams, database.BulkInsertCommitsParams, UsersToRepoUrl, int, error) {
	// sync this from the db

	repoDir := filepath.Join(params.prefixPath, params.organization, params.repo)
	dependencyHistoryObject := database.BatchInsertRepoDependenciesParams{
		Url:             params.repoUrl,
		Firstusedates:   []int64{},
		Lastusedates:    []int64{},
		Dependencynames: []string{},
		Filenames:       []string{},
		UpdatedAt:       sql.NullInt64{Int64: lib.Now().Unix(), Valid: true},
	}

	for _, dep := range currentDependencies {
		dependencyHistoryObject.Dependencynames = append(dependencyHistoryObject.Dependencynames, dep.DependencyName)
		dependencyHistoryObject.Filenames = append(dependencyHistoryObject.Filenames, dep.DependencyFile)
		dependencyHistoryObject.Firstusedates = append(dependencyHistoryObject.Firstusedates, dep.FirstUseDate.Int64)
		dependencyHistoryObject.Lastusedates = append(dependencyHistoryObject.Lastusedates, dep.LastUseDate.Int64)
	}
	commitWindow := GetCommitWindow(len(commitList))

	usersToRepoUrl := UsersToRepoUrl{
		AuthorEmails: []string{},
	}

	commitObject := database.BulkInsertCommitsParams{
		Repourl: sql.NullString{
			String: params.repoUrl,
			Valid:  true,
		},
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
	firstViewableCommitIndex := len(commitList) - numberOfCommits
	for commitIndex, commit := range commitList {

		if commitIndex >= firstViewableCommitIndex {
			if commitIndex%commitWindow == 0 {
				err = CheckCommitForDependencies(commit, repoDir, &dependencyHistoryObject, dependencyFiles)
				if err != nil {
					return dependencyHistoryObject, commitObject, usersToRepoUrl, 0, err
				}
			}

			AddCommitToCommitObject(commit, &commitObject, commitIndex)
			AddFirstLastCommitDateByEmail(&usersToRepoUrl, commit)

		}

	}
	c := commitList[len(commitList)-1]
	// always check last commit last
	err = CheckCommitForDependencies(c, repoDir, &dependencyHistoryObject, dependencyFiles)

	return dependencyHistoryObject, commitObject, usersToRepoUrl, numberOfCommits, err
}
