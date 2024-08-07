package usersync

import (
	"github.com/OpenQDev/GoGitguru/database"
)

func convertDatabaseObjectToUserSync(newCommitAuthorsRaw []database.GetLatestUncheckedCommitPerAuthorRow) []UserSync {
	var newCommitAuthors []UserSync

	for _, author := range newCommitAuthorsRaw {
		var authorEmail string
		if author.AuthorEmail.Valid {
			authorEmail = author.AuthorEmail.String
		}

		var repoUrl string
		if author.RepoUrl.Valid {
			repoUrl = author.RepoUrl.String
		}
		if !author.GithubUserEmail.Valid {
			newCommitAuthors = append(newCommitAuthors, UserSync{
				CommitHash:  author.CommitHash,
				AuthorEmail: authorEmail,
				RepoUrl:     repoUrl,
			})
		}
	}

	return newCommitAuthors
}
