package usersync

import (
	"context"
	"strings"

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

	if err != nil {
		logger.LogFatalRedAndExit("error getting new commit authors to process: %s", err)
		return
	}
	if newCommitAuthorsRaw != nil {

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

			UpsertRepoToUserByIdParams := database.UpsertRepoToUserByIdParams{
				Url: repoToAuthorBatch.RepoURL,
			}

			for _, commitAuthor := range githubGraphQLCommitAuthors {
				author := commitAuthor.Author

				err := insertIntoRestIdToUser(author, db)
				if err != nil {
					logger.LogError("error occured while inserting author RestID %s to Email %s: %s", author.User.GithubRestID, author.Email, err)
				}

				internal_id, err := db.CheckGithubUserExists(context.Background(), strings.ToLower(author.User.Login))
				if err != nil {
					logger.LogBlue("github user does not exist: %s", err)
				}
				// TODO update their for that specific repo.
				if internal_id == 0 {
					logger.LogBlue("inserting github user %s", author.Name)
					internal_id, err = insertGithubUser(author, db)
					if err != nil {
						logger.LogError("error occured while inserting github user %s with RestId %s: %s", author.User.Login, author.User.GithubRestID, err)
					} else {
						logger.LogGreen("user %s inserted!", author.Name)
					}
				}

				err = GetReposToUsers(db, &UpsertRepoToUserByIdParams, internal_id, author)
				if err != nil {
					logger.LogError("error occured while getting repos to users: %s", err)
				}
			}
			if err != nil {
				logger.LogError("error occured while getting repos to users: %s", err)
			}

			err = db.UpsertRepoToUserById(context.Background(), UpsertRepoToUserByIdParams)
			if err != nil {
				logger.LogError("error occured while upserting repo to user by id: %s", err)
			}
		}

	}

}
