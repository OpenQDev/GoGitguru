package usersync

import (
	"context"
	"database/sql"
	"slices"
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
			internalIds := []int32{}
			lastCommitDates := []int64{}
			firstCommitDates := []int64{}

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
				userCommits, err := db.GetFirstAndLastCommit(context.Background(), sql.NullString{String: author.Email, Valid: true})
				if err != nil {
					logger.LogError("error occured while getting first and last commit for user %s: %s", author.Email, err)
				}
				alreadySet := slices.Contains(internalIds, internal_id)
				if alreadySet {

					//change the first and last commit dates if the current commit is earlier or later than the current first and last commit dates
					for index, id := range internalIds {
						if id == internal_id {
							if userCommits.FirstCommitDate.(int64) < firstCommitDates[index] {
								firstCommitDates[index] = userCommits.FirstCommitDate.(int64)
							}
							if userCommits.LastCommitDate.(int64) > lastCommitDates[index] {
								lastCommitDates[index] = userCommits.LastCommitDate.(int64)
							}
						}
					}
				}

				if !alreadySet {
					internalIds = append(internalIds, internal_id)
					lastCommitDates = append(lastCommitDates, userCommits.LastCommitDate.(int64))
					firstCommitDates = append(firstCommitDates, userCommits.FirstCommitDate.(int64))
				}

			}
			UpsertRepoToUserByIdParams := database.UpsertRepoToUserByIdParams{
				Url:              repoToAuthorBatch.RepoURL,
				InternalIds:      internalIds,
				FirstCommitDates: firstCommitDates,
				LastCommitDates:  lastCommitDates,
			}

			err = SyncUserDependencies(db)
			if err != nil {
				logger.LogFatalRedAndExit("error syncing dependencies: %s", err)
				return
			}
			err = db.UpsertRepoToUserById(context.Background(), UpsertRepoToUserByIdParams)
			if err != nil {
				logger.LogError("error occured while upserting repo to user by id: %s", err)
			}
		}

	}

}
