package sync

// Create a map with repoUrl as key and array of authors as value
func GetRepoToAuthorsMap(newCommitAuthors []UserSync) map[string][]string {
	repoAuthorsMap := make(map[string][]string)

	for _, author := range newCommitAuthors {
		if author.Repo.NotNull {
			repoAuthorsMap[author.Repo.URL] = append(repoAuthorsMap[author.Repo.URL], author.Author.Email)
		}
	}

	return repoAuthorsMap
}
