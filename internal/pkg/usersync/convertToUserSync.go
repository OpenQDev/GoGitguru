package usersync

import "main/internal/database"

func ConvertDatabaseObjectToUserSync(newCommitAuthorsRaw []database.GetLatestUncheckedCommitPerAuthorRow) []UserSync {
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
