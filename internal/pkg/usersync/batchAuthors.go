package usersync

type BatchAuthor struct {
	RepoURL string
	Tuples  []AuthorCommitTuple
}

type BatchAuthors = []BatchAuthor

func generateBatchAuthors(repoUrlToAuthorsMap RepoToAuthorCommitTuples, batchSize int) BatchAuthors {
	var result BatchAuthors

	for repoUrl, authors := range repoUrlToAuthorsMap.Repos {
		for i := 0; i < len(authors); i += batchSize {
			end := i + batchSize
			if end > len(authors) {
				end = len(authors)
			}

			batch := authors[i:end]
			batchAuthor := BatchAuthor{
				RepoURL: repoUrl,
				Tuples:  batch,
			}
			result = append(result, batchAuthor)
		}
	}

	return result
}
