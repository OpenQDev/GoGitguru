package reposync

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"slices"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-git/go-git/v5/plumbing/object"
)

type UsersToRepoUrl struct {
	AuthorEmails     []string
	FirstCommitDates []int64
	LastCommitDates  []int64
}

func GetObjectsFromCommitList(params GitLogParams, commitList []*object.Commit, numberOfCommits int, currentDependencies []database.GetRepoDependenciesByURLRow) (database.BatchInsertRepoDependenciesParams, database.BulkInsertCommitsParams, UsersToRepoUrl, int, error) {
	// sync this from the db
	repoDir := filepath.Join(params.prefixPath, params.organization, params.repo)
	dependencyHistoryObject := database.BatchInsertRepoDependenciesParams{
		Url:             params.repoUrl,
		Firstusedates:   []int64{},
		Lastusedates:    []int64{},
		Dependencynames: []string{},
		Filenames:       []string{},
		UpdatedAt:       sql.NullInt64{Int64: now().Unix(), Valid: true},
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
	for commitCount, commit := range commitList {
		if commitCount >= numberOfCommits {
			println("commit count is greater than or equal to number of commits check", commitCount, numberOfCommits)
			break
		}
		if commitCount < numberOfCommits {
			if commitCount%commitWindow == 0 {
				err = CheckCommitForDependencies(commit, repoDir, &dependencyHistoryObject)
				if err != nil {
					return dependencyHistoryObject, commitObject, usersToRepoUrl, 0, err
				}
			}
			AddCommitToCommitObject(commit, &commitObject, commitCount)
			alreadyHasEmail := slices.Contains(usersToRepoUrl.AuthorEmails, commit.Author.Email)
			if alreadyHasEmail {

				for index, email := range usersToRepoUrl.AuthorEmails {
					if email == commit.Author.Email {

						if commit.Author.When.Unix() < usersToRepoUrl.FirstCommitDates[index] {
							usersToRepoUrl.FirstCommitDates[index] = commit.Author.When.Unix()
						}
						if commit.Author.When.Unix() > usersToRepoUrl.LastCommitDates[index] {
							usersToRepoUrl.LastCommitDates[index] = commit.Author.When.Unix()
						}
						break
					}
				}
			}
			if !alreadyHasEmail {
				println("adding email to users to repo url", commit.Author.When.Unix())
				usersToRepoUrl.AuthorEmails = append(usersToRepoUrl.AuthorEmails, commit.Author.Email)
				usersToRepoUrl.FirstCommitDates = append(usersToRepoUrl.FirstCommitDates, commit.Author.When.Unix())
				usersToRepoUrl.LastCommitDates = append(usersToRepoUrl.LastCommitDates, commit.Author.When.Unix())
			}
		}

	}
	c := commitList[len(commitList)-1]
	// always check last commit last
	fmt.Printf("Commit number %d: %s\n", len(commitList)-1, c.Hash)
	err = CheckCommitForDependencies(c, repoDir, &dependencyHistoryObject)
	return dependencyHistoryObject, commitObject, usersToRepoUrl, numberOfCommits, err
}
