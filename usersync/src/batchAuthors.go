package usersync

import "sort"

type BatchAuthor struct {
	RepoURL            string
	AuthorCommitTuples []AuthorCommitTuple
}

type BatchAuthors = []BatchAuthor

func generateBatchAuthors(repoUrlToAuthorsMap RepoToAuthorCommitTuples, batchSize int) BatchAuthors {
	var result BatchAuthors

	// Get the keys and sort them
	keys := make([]string, 0, len(repoUrlToAuthorsMap.Repos))
	for k := range repoUrlToAuthorsMap.Repos {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Iterate over the sorted keys
	for _, repoUrl := range keys {
		authors := repoUrlToAuthorsMap.Repos[repoUrl]
		for i := 0; i < len(authors); i += batchSize {
			end := i + batchSize
			if end > len(authors) {
				end = len(authors)
			}

			batch := authors[i:end]
			batchAuthor := BatchAuthor{
				RepoURL:            repoUrl,
				AuthorCommitTuples: batch,
			}
			result = append(result, batchAuthor)
		}
	}

	return result
}
