package reposync

import (
	"context"
	"fmt"
	"path/filepath"
	"slices"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/OpenQDev/GoGitguru/util/gitutil"
)

type GitLogParams struct {
	prefixPath     string
	organization   string
	repo           string
	repoUrl        string
	fromCommitDate time.Time
	db             *database.Queries
}

func getFirstAndLastUserCommitByEmail(usersToReposObject UsersToRepoUrl, emails []string, firstCommitDate int64, lastCommitDate int64) (int64, int64) {
	for index, userEmail := range usersToReposObject.AuthorEmails {
		if slices.Contains(emails, userEmail) {

			resultFirstCommitDate := firstCommitDate
			resultLastCommitDate := lastCommitDate

			if firstCommitDate > usersToReposObject.FirstCommitDates[index] || firstCommitDate == 0 {
				resultFirstCommitDate = usersToReposObject.FirstCommitDates[index]
			}
			if lastCommitDate < usersToReposObject.LastCommitDates[index] || lastCommitDate == 0 {
				resultLastCommitDate = usersToReposObject.LastCommitDates[index]
			}
			return resultFirstCommitDate, resultLastCommitDate
		}
	}
	return firstCommitDate, lastCommitDate
}

func getUpsertRepoByIdsParams(params GitLogParams, usersToReposObject UsersToRepoUrl) database.UpsertRepoToUserByIdParams {

	internalIdsWithEmails, err := params.db.GetGithubUserByCommitEmail(context.Background(), usersToReposObject.AuthorEmails)
	if err != nil {
		fmt.Println("Error getting internal ids with emails", err)
	}

	insertByIdParams := database.UpsertRepoToUserByIdParams{
		InternalIds:      []int32{},
		Url:              params.repoUrl,
		FirstCommitDates: []int64{},
		LastCommitDates:  []int64{},
	}

	for _, internalIdWithEmail := range internalIdsWithEmails {

		emails := internalIdWithEmail.Emails

		if err != nil {
			fmt.Println("Error getting emails from internal id", err)
		}

		alreadyHas := slices.Contains(insertByIdParams.InternalIds, internalIdWithEmail.InternalID)
		if alreadyHas {
			for insertByIdParamIndex, insertByIdParam := range insertByIdParams.InternalIds {
				if insertByIdParam == internalIdWithEmail.InternalID {
					currentFirstCommit := insertByIdParams.FirstCommitDates[insertByIdParamIndex]
					currentLastCommit := insertByIdParams.LastCommitDates[insertByIdParamIndex]
					firstCommit, lastCommit := getFirstAndLastUserCommitByEmail(usersToReposObject, emails, currentFirstCommit, currentLastCommit)
					insertByIdParams.FirstCommitDates[insertByIdParamIndex] = firstCommit
					insertByIdParams.LastCommitDates[insertByIdParamIndex] = lastCommit
				}
			}

		} else {
			insertByIdParams.InternalIds = append(insertByIdParams.InternalIds, internalIdWithEmail.InternalID)
			firstCommitDate, lastCommitDate := getFirstAndLastUserCommitByEmail(usersToReposObject, emails, 0, 0)
			insertByIdParams.FirstCommitDates = append(insertByIdParams.FirstCommitDates, firstCommitDate)
			insertByIdParams.LastCommitDates = append(insertByIdParams.LastCommitDates, lastCommitDate)
		}
	}
	return insertByIdParams

}

// from commitDate should be the date of the last commit that was synced for the repository or any of the dependencies.

func StoreGitLogsAndDepsHistoryForRepo(params GitLogParams) (int, error) {

	repoDir := filepath.Join(params.prefixPath, params.organization, params.repo)

	commitList, err := CreateCommitList(repoDir)

	if err != nil {
		return 0, fmt.Errorf("error getting commit list %s: %s", params.repoUrl, err)
	}

	numberOfCommitsToSync, err := gitutil.GetNumberOfCommits(params.prefixPath, params.organization, params.repo, params.fromCommitDate)
	if err != nil {
		return 0, fmt.Errorf("error getting number of commits for %s: %s", params.repoUrl, err)
	}

	currentDependencies, err := params.db.GetRepoDependenciesByURL(context.Background(), params.repoUrl)

	if err != nil {
		return 0, fmt.Errorf("error getting current dependencies for %s: %s", params.repoUrl, err)
	}

	dependencyHistoryObjects, commitObject, usersToReposObject, numberOfCommitsToSync, err := GetObjectsFromCommitList(params, commitList, numberOfCommitsToSync, currentDependencies)

	if err != nil {
		return 0, fmt.Errorf("error getting structs from commit list %s: %s", params.repoUrl, err)
	}

	insertByIdParams := getUpsertRepoByIdsParams(params, usersToReposObject)

	err = params.db.UpsertRepoToUserById(context.Background(), insertByIdParams)

	if err != nil {
		return 0, fmt.Errorf("error storing users to repo for %s: %s", params.repoUrl, err)
	}

	err = params.db.BatchInsertRepoDependencies(context.Background(), dependencyHistoryObjects)

	if err != nil {
		fmt.Printf("error storing dependency history for %s: %s", params.repoUrl, err)
		fmt.Println("dependencyHistoryObjects", dependencyHistoryObjects)
	}

	err = params.db.BulkInsertCommits(context.Background(), commitObject)

	if err != nil {
		return 0, fmt.Errorf("error storing commits for %s: %s", params.repoUrl, err)
	}
	return numberOfCommitsToSync, nil
}
