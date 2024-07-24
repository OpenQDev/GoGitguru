package usersync

import (
	"context"
	"fmt"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/logger"
)

type UserSync struct {
	CommitHash  string
	AuthorEmail string
	RepoUrl     string
}

func StartUserSyncing(
	db *database.Queries,
	prefixPath string,
	ghAccessToken string,
	batchSize int,
	githubGraphQLUrl string,
) {
	newCommitAuthorsRaw, err := getNewCommitAuthors(db)
	fmt.Println(newCommitAuthorsRaw)

	if err != nil {
		logger.LogFatalRedAndExit("error getting new commit authors to process: %s", err)
		return
	}

	if newCommitAuthorsRaw != nil {
		fmt.Println(newCommitAuthorsRaw[0].AuthorEmail)
		logger.LogBlue("identifying %d new authors", len(newCommitAuthorsRaw))

		// Convert to database object to local type
		newCommitAuthors := convertDatabaseObjectToUserSync(newCommitAuthorsRaw)

		// Create map of repoUrl -> []authors
		repoUrlToAuthorsMap := getRepoToAuthorsMap(newCommitAuthors)

		// Create batches of repos for GraphQL query
		repoToAuthorBatches := generateBatchAuthors(repoUrlToAuthorsMap, batchSize)

		// Get info for each batch
		for _, repoToAuthorBatch := range repoToAuthorBatches {

			githubGraphQLCommitAuthorsMap, err := identifyRepoAuthorsBatch(repoToAuthorBatch.RepoURL, repoToAuthorBatch.AuthorCommitTuples, ghAccessToken, githubGraphQLUrl)

			if err != nil {
				logger.LogError("error occured while identifying authors: %s", err)
			}

			logger.LogGreenDebug("successfully fetched info for batch %s", repoToAuthorBatch.RepoURL)

			if githubGraphQLCommitAuthorsMap == nil {
				logger.LogError("commits is nil")
				continue
			}

			githubGraphQLCommitAuthors := make([]GithubGraphQLCommit, 0, len(githubGraphQLCommitAuthorsMap))

			for _, commitAuthor := range githubGraphQLCommitAuthorsMap {
				githubGraphQLCommitAuthors = append(githubGraphQLCommitAuthors, commitAuthor)
			}

			upsertRepoToUserByIdParams := database.UpsertRepoToUserByIdParams{
				Url: repoToAuthorBatch.RepoURL,
			}

			for _, commitAuthor := range githubGraphQLCommitAuthors {
				author := commitAuthor.Author

				err := insertIntoRestIdToUser(author, db)
				if err != nil {
					logger.LogError("error occured while inserting author RestID %s to Email %s: %s", author.User.GithubRestID, author.Email, err)
				}

				result, err := db.CheckGithubUserIdExists(context.Background(), author.User.GithubRestID)
				if err != nil {
					logger.LogError("error checking if github user exists: %s", err)
				}
				// TODO update their for that specific repo.
				if !result {
					logger.LogBlue("inserting github user %s", author.Name)
					err := insertGithubUser(author, db)
					if err != nil {
						logger.LogError("error occured while inserting github user %s with RestId %s: %s", author.User.Login, author.User.GithubRestID, err)
					} else {
						logger.LogGreen("user %s inserted!", author.Name)
					}

				}

				internal_id, err := db.GetGithubUserByRestId(context.Background(), author.User.GithubRestID)

				if err != nil {
					logger.LogError("error occured while getting GetGithubUserByRestId: %s", err)
				}

				err = GetReposToUsers(db, &upsertRepoToUserByIdParams, internal_id, author)

				if err != nil {
					logger.LogError("error occured while getting repos to users: %s", err)
				}
			}

			if err != nil {
				logger.LogError("error occured while getting repos to users: %s", err)
			}

			err = db.UpsertRepoToUserById(context.Background(), upsertRepoToUserByIdParams)
			if err != nil {
				logger.LogError("error occured while upserting repo to user by id: %s", err)
			}
		}

		// update to hasSynced
		newCommitEmails := make([]string, 0, len(newCommitAuthorsRaw))
		for _, commitAuthor := range newCommitAuthorsRaw {

			newCommitEmails = append(newCommitEmails, commitAuthor.AuthorEmail.String)
		}
		err = db.SetAllCommitsToChecked(context.Background(), newCommitEmails)
		if err != nil {
			logger.LogError("error occured while setting all commits to checked: %s", err)
			return
		}

	} else {
		fmt.Println("no new commits to sync")
		return
	}
}
