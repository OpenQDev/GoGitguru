package sync

import (
	"context"
	"fmt"
	"main/internal/database"
	"main/internal/pkg/logger"
	"time"
)

type UserSync struct {
	CommitHash string
	Author     struct {
		Email   string
		NotNull bool
	}
	Repo struct {
		URL     string
		NotNull bool
	}
}

func StartSyncingUser(
	db *database.Queries,
	prefixPath string,
	concurrency int,
	timeBetweenSyncs time.Duration) {
	/*
	   [
	     {
	       65062be663cc004b77ca8a3b13255bc5efa42f25
	       {andrew@openq.dev true}
	       {https://github.com/openqdev/openq-workflows true}
	     }
	   ]
	   **/

	newCommitAuthorsRaw, err := db.GetLatestUncheckedCommitPerAuthor(context.Background())

	if len(newCommitAuthorsRaw) == 0 {
		logger.LogBlue("No new authors to process.")
		return
	}

	newCommitAuthors := convertToUserSync(newCommitAuthorsRaw)

	if err != nil {
		logger.LogError("errerrerr", err)
	}

	fmt.Println(newCommitAuthors)

	logger.LogBlue("identifying %d new authors", len(newCommitAuthors))

	repoUrlToAuthorsMap := getRepoToAuthorsMap(newCommitAuthors)

	fmt.Println(repoUrlToAuthorsMap)
}

func convertToUserSync(newCommitAuthorsRaw []database.GetLatestUncheckedCommitPerAuthorRow) []UserSync {
	var newCommitAuthors []UserSync

	for _, author := range newCommitAuthorsRaw {
		newCommitAuthors = append(newCommitAuthors, UserSync{
			CommitHash: author.CommitHash,
			Author: struct {
				Email   string
				NotNull bool
			}{
				Email:   author.AuthorEmail.String,
				NotNull: author.AuthorEmail.Valid,
			},
			Repo: struct {
				URL     string
				NotNull bool
			}{
				URL:     author.RepoUrl.String,
				NotNull: author.RepoUrl.Valid,
			},
		})
	}

	return newCommitAuthors
}

/*
Creates batches of authors grouped by repo_url.

Let'say we have a list of 500 unknown authors (ordered by repo url).
The first 250 authors belong to repo A, the next 150 belong to B,
the next 90 belong to C, and the last 10 belong to D, E, F, G, ..., M.

The function will return 16 batches:
A: 100, 100, 50, B: 100, 50, C: 90, D: 1, E: 1, F: 1, ..., M: 1

So identifying 500 authors requires 16 API calls instead of 500 because we
can group authors by repo url in the GraphQL query.
*/
func getRepoToAuthorsMap(newCommitAuthors []UserSync) map[string][]string {
	// Create a map with repoUrl as key and array of authors as value
	repoAuthorsMap := make(map[string][]string)

	for _, author := range newCommitAuthors {
		if author.Repo.NotNull {
			repoAuthorsMap[author.Repo.URL] = append(repoAuthorsMap[author.Repo.URL], author.Author.Email)
		}
	}

	return repoAuthorsMap
}
