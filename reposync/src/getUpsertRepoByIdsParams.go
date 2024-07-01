package reposync

import (
	"context"
	"fmt"
	"slices"

	"github.com/OpenQDev/GoGitguru/database"
)

func GetUpsertRepoByIdsParams(params GitLogParams, usersToReposObject UsersToRepoUrl) database.UpsertRepoToUserByIdParams {

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
					firstCommit, lastCommit := GetFirstAndLastUserCommitByEmail(usersToReposObject, emails, currentFirstCommit, currentLastCommit)
					insertByIdParams.FirstCommitDates[insertByIdParamIndex] = firstCommit
					insertByIdParams.LastCommitDates[insertByIdParamIndex] = lastCommit
				}
			}

		} else {
			insertByIdParams.InternalIds = append(insertByIdParams.InternalIds, internalIdWithEmail.InternalID)
			firstCommitDate, lastCommitDate := GetFirstAndLastUserCommitByEmail(usersToReposObject, emails, 0, 0)
			insertByIdParams.FirstCommitDates = append(insertByIdParams.FirstCommitDates, firstCommitDate)
			insertByIdParams.LastCommitDates = append(insertByIdParams.LastCommitDates, lastCommitDate)
		}
	}
	return insertByIdParams

}
