package usersync

import (
	"main/internal/database"
	"main/internal/pkg/githubGraphQL"
	"main/internal/pkg/logger"
	"main/internal/pkg/server"
	"time"
)

type UserSync struct {
	CommitHash  string
	AuthorEmail string
	RepoUrl     string
}

func StartSyncingUser(
	db *database.Queries,
	prefixPath string,
	concurrency int,
	timeBetweenSyncs time.Duration,
	ghAccessToken string,
	batchSize int,
	apiCfg server.ApiConfig,
) {
	newCommitAuthorsRaw, err := getNewCommitAuthors(db)
	if err != nil {
		logger.LogFatalRedAndExit("error getting new commit authors to process: %s", err)
		return
	}
	if newCommitAuthorsRaw == nil {
		logger.LogBlue("no new authors to sync")
		return
	}

	logger.LogBlue("identifying %d new authors", len(newCommitAuthorsRaw))

	// Convert to database object to local type
	newCommitAuthors := convertDatabaseObjectToUserSync(newCommitAuthorsRaw)

	// Create map of repoUrl -> []authors
	repoUrlToAuthorsMap := getRepoToAuthorsMap(newCommitAuthors)

	// Create batches of repos for GraphQL query
	repoToAuthorBatches := generateBatchAuthors(repoUrlToAuthorsMap, 2)

	// Get info for each batch
	for _, repoToAuthorBatch := range repoToAuthorBatches {
		logger.LogGreenDebug("%s", repoToAuthorBatch.RepoURL)

		githubGraphQLCommitAuthorsMap, err := identifyRepoAuthorsBatch(repoToAuthorBatch.RepoURL, repoToAuthorBatch.AuthorCommitTuples, ghAccessToken, apiCfg)

		logger.LogGreenDebug("successfully fetched info for batch %s", repoToAuthorBatch.RepoURL)

		if githubGraphQLCommitAuthorsMap == nil {
			logger.LogError("commits is nil")
			continue
		}

		githubGraphQLCommitAuthors := make([]githubGraphQL.GithubGraphQLCommit, 0, len(githubGraphQLCommitAuthorsMap))

		for _, commitAuthor := range githubGraphQLCommitAuthorsMap {
			githubGraphQLCommitAuthors = append(githubGraphQLCommitAuthors, commitAuthor)
		}

		if err != nil {
			logger.LogError("error occured while identifying authors: %s", err)
		}

		logger.LogGreenDebug("got the following info: %v", githubGraphQLCommitAuthorsMap)

		for _, commitAuthor := range githubGraphQLCommitAuthors {
			author := commitAuthor.Author

			err := insertIntoRestIdToUser(author, db)
			if err != nil {
				logger.LogError("error occured while inserting RestID to Email: %s", err)
			}

			err = insertGithubUser(author, db)
			if err != nil {
				logger.LogError("error occured while inserting author: %s", err)
			}

		}
	}
}
