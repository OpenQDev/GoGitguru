package sync

import (
	"context"
	"database/sql"
	"fmt"
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
		commitValues := make([]gitutil.Commit, 0, len(*commits))

		for _, value := range *commits {
			commitValues = append(commitValues, value)
		}

		fmt.Println("commitValues", commitValues)

		if err != nil {
			logger.LogError("error occured while identifying authors: %s", err)
		}

		logger.LogGreenDebug("successfully fetched info for batch %s", repoToAuthorBatch.RepoURL)

		logger.LogGreenDebug("got the following info: %v", commits)

		for _, commit := range commitValues {
			author := commit.Author

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
				Name:            sql.NullString{String: author.User.Name, Valid: author.User.Name != ""},
				Email:           sql.NullString{String: author.User.Email, Valid: author.User.Email != ""},
				AvatarUrl:       sql.NullString{String: author.User.AvatarURL, Valid: author.User.AvatarURL != ""},
				Company:         sql.NullString{String: *author.User.Company, Valid: author.User.Company != nil},
				Location:        sql.NullString{String: *author.User.Location, Valid: author.User.Location != nil},
				Bio:             sql.NullString{String: author.User.Bio, Valid: author.User.Bio != ""},
				Blog:            sql.NullString{String: *author.User.Blog, Valid: author.User.Blog != nil},
				Hireable:        sql.NullBool{Bool: author.User.Hireable, Valid: true},
				TwitterUsername: sql.NullString{String: *author.User.TwitterUsername, Valid: author.User.TwitterUsername != nil},
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
