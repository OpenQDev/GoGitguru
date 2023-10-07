package sync

import (
	"context"
	"database/sql"
	"main/internal/database"
	"main/internal/pkg/gitutil"
	"main/internal/pkg/logger"
	"time"
)

func StartSyncingUser(
	db *database.Queries,
	prefixPath string,
	concurrency int,
	timeBetweenSyncs time.Duration,
	ghAccessToken string,
	batchSize int,
) {
	newCommitAuthorsRaw, err := db.GetLatestUncheckedCommitPerAuthor(context.Background())
	if err != nil {
		logger.LogError("errerrerr", err)
	}

	logger.LogGreenDebug("new commit authors to check: %s", newCommitAuthorsRaw)

	if len(newCommitAuthorsRaw) == 0 {
		logger.LogBlue("No new authors to process.")
		return
	}

	logger.LogBlue("identifying %d new authors", len(newCommitAuthorsRaw))

	// Convert to database object to local type
	newCommitAuthors := ConvertToUserSync(newCommitAuthorsRaw)
	logger.LogGreenDebug("newCommitAuthors", newCommitAuthors)

	// Create map of repoUrl -> []authors
	repoUrlToAuthorsMap := GetRepoToAuthorsMap(newCommitAuthors)
	logger.LogGreenDebug("repoUrlToAuthorsMap", repoUrlToAuthorsMap)

	// Create batches of repos for GraphQL query
	repoToAuthorBatches := GenerateBatchAuthors(repoUrlToAuthorsMap, 2)
	logger.LogGreenDebug("repoToAuthorBatches", repoToAuthorBatches)

	// Get info for each batch
	for _, repoToAuthorBatch := range repoToAuthorBatches {
		logger.LogGreenDebug("%s", repoToAuthorBatch.RepoURL)

		commits, err := IdentifyRepoAuthorsBatch(repoToAuthorBatch.RepoURL, repoToAuthorBatch.Tuples, ghAccessToken)

		if commits == nil {
			logger.LogError("commits is nil")
			continue
		}

		commitValues := make([]gitutil.Commit, 0, len(*commits))

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

			_, err = db.InsertRestIdToEmail(context.Background(), database.InsertRestIdToEmailParams{
				RestID: sql.NullInt32{Int32: int32(restId), Valid: true},
				Email:  author.Email,
			})
			if err != nil {
				logger.LogError("error occured while inserting RestID to Email: %s", err)
			}

			createdAt, err := time.Parse(time.RFC3339, author.User.CreatedAt)
			if err != nil {
				logger.LogError("error parsing time: %s", err)
			}
			updatedAt, err := time.Parse(time.RFC3339, author.User.UpdatedAt)
			if err != nil {
				logger.LogError("error parsing time: %s", err)
			}

			_, err = db.InsertUser(context.Background(), database.InsertUserParams{
				GithubRestID:    int32(author.User.GithubRestID),
				GithubGraphqlID: author.User.GithubGraphqlID,
				Login:           author.User.Login,
				Name:            sql.NullString{String: author.User.Name, Valid: true},
				Email:           sql.NullString{String: author.User.Email, Valid: true},
				AvatarUrl:       sql.NullString{String: author.User.AvatarURL, Valid: true},
				Company:         sql.NullString{String: author.User.Company, Valid: true},
				Location:        sql.NullString{String: author.User.Location, Valid: true},
				Bio:             sql.NullString{String: author.User.Bio, Valid: true},
				Blog:            sql.NullString{String: author.User.Blog, Valid: true},
				Hireable:        sql.NullBool{Bool: author.User.Hireable, Valid: true},
				TwitterUsername: sql.NullString{String: author.User.TwitterUsername, Valid: true},
				Followers:       sql.NullInt32{Int32: int32(author.User.Followers.TotalCount), Valid: true},
				Following:       sql.NullInt32{Int32: int32(author.User.Following.TotalCount), Valid: true},
				Type:            "User",
				CreatedAt:       sql.NullTime{Time: createdAt, Valid: true},
				UpdatedAt:       sql.NullTime{Time: updatedAt, Valid: true},
			})

			if err != nil {
				logger.LogError("error occured while inserting author: %s", err)
			}

		}
	}
}
