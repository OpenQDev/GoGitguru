package sync

func BatchAuthors(repoUrlToAuthorsMap map[string][]string, batchSize int) [][]interface{} {
	var result [][]interface{}

	for repoUrl, authors := range repoUrlToAuthorsMap {
		for i := 0; i < len(authors); i += batchSize {
			end := i + batchSize
			if end > len(authors) {
				end = len(authors)
			}

			batch := authors[i:end]
			result = append(result, []interface{}{repoUrl, batch})
		}
	}

	return result
}
