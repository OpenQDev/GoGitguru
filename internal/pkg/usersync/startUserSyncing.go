package usersync

import (
	"context"
	"main/internal/database"
	"main/internal/pkg/githubGraphQL"
	"main/internal/pkg/logger"
	"time"
)

type UserSync struct {
	CommitHash  string
	AuthorEmail *string
	RepoUrl     *string
}

func StartSyncingUser(
	db *database.Queries,
	prefixPath string,
	concurrency int,
	timeBetweenSyncs time.Duration,
	ghAccessToken string,
	batchSize int,
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
	repoUrlToAuthorsMap := GetRepoToAuthorsMap(newCommitAuthors)

	// Create batches of repos for GraphQL query
	repoToAuthorBatches := GenerateBatchAuthors(repoUrlToAuthorsMap, 2)

	// Get info for each batch
	for _, repoToAuthorBatch := range repoToAuthorBatches {
		logger.LogGreenDebug("%s", repoToAuthorBatch.RepoURL)

		commits, err := IdentifyRepoAuthorsBatch(repoToAuthorBatch.RepoURL, repoToAuthorBatch.Tuples, ghAccessToken)

		if commits == nil {
			logger.LogError("commits is nil")
			continue
		}

		commitValues := make([]githubGraphQL.Commit, 0, len(*commits))

		for _, value := range *commits {
			commitValues = append(commitValues, value)
		}

		if err != nil {
			logger.LogError("error occured while identifying authors: %s", err)
		}

		logger.LogGreenDebug("successfully fetched info for batch %s", repoToAuthorBatch.RepoURL)

		logger.LogGreenDebug("got the following info: %v", commits)

		for _, commit := range commitValues {
			restId := commit.Author.User.GithubRestID
			author := commit.Author

			var params database.InsertRestIdToEmailParams
			if restId == 0 {
				params = database.InsertRestIdToEmailParams{
					Email: author.Email,
				}
			} else {
				params = database.InsertRestIdToEmailParams{
					RestID: int32(restId),
					Email:  author.Email,
				}
			}

			_, err = db.InsertRestIdToEmail(context.Background(), params)
			if err != nil {
				logger.LogError("error occured while inserting RestID to Email: %s", err)
			}

			createdAt, err := time.Parse(time.RFC3339, author.User.CreatedAt)
			if err != nil && !createdAt.IsZero() {
				logger.LogError("error parsing time: %s", err)
			}

			updatedAt, err := time.Parse(time.RFC3339, author.User.UpdatedAt)
			if err != nil && !createdAt.IsZero() {
				logger.LogError("error parsing time: %s", err)
			}

			authorParams := convertAuthorToInsertUserParams(author, createdAt, updatedAt)

			_, err = db.InsertUser(context.Background(), authorParams)

			if err != nil {
				logger.LogError("error occured while inserting author: %s", err)
			}

		}
	}
}
