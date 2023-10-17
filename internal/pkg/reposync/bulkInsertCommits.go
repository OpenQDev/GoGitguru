package reposync

import (
	"context"
	"main/internal/database"
)

func BulkInsertCommits(
	db *database.Queries,
	commitHash []string,
	author []string,
	authorEmail []string,
	authorDate []int64,
	committerDate []int64,
	message []string,
	insertions []int32,
	deletions []int32,
	filesChanged []int32,
	repoUrls []string,
) error {
	params := database.BulkInsertCommitsParams{
		Column1:  commitHash,
		Column2:  author,
		Column3:  authorEmail,
		Column4:  authorDate,
		Column5:  committerDate,
		Column6:  message,
		Column7:  insertions,
		Column8:  deletions,
		Column9:  filesChanged,
		Column10: repoUrls,
	}

	return db.BulkInsertCommits(context.Background(), params)
}
