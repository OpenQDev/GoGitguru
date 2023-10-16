package usersync

type AuthorCommitTuple struct {
	Author     string
	CommitHash string
}

// Create a map with repoUrl as key and array of authors as value
func GetRepoToAuthorsMap(newCommitAuthors []UserSync) map[string][]AuthorCommitTuple {
	repoAuthorsMap := make(map[string][]AuthorCommitTuple)

	for _, author := range newCommitAuthors {
		if author.Repo.NotNull {
			repoAuthorsMap[author.Repo.URL] = append(repoAuthorsMap[author.Repo.URL], AuthorCommitTuple{Author: author.Author.Email, CommitHash: author.CommitHash})
		}
	}

	return repoAuthorsMap
}
