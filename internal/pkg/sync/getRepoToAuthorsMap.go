package sync

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
func GetRepoToAuthorsMap(newCommitAuthors []UserSync) map[string][]string {
	// Create a map with repoUrl as key and array of authors as value
	repoAuthorsMap := make(map[string][]string)

	for _, author := range newCommitAuthors {
		if author.Repo.NotNull {
			repoAuthorsMap[author.Repo.URL] = append(repoAuthorsMap[author.Repo.URL], author.Author.Email)
		}
	}

	return repoAuthorsMap
}
