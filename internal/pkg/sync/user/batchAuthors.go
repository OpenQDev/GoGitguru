package usersync

func BatchAuthors(repoUrlToAuthorsMap map[string][]string, batchSize int) []interface{} {
	var batches []interface{}

	for repoUrl, authorList := range repoUrlToAuthorsMap {
		for i := 0; i < len(authorList); i += batchSize {
			end := i + batchSize
			if end > len(authorList) {
				end = len(authorList)
			}
			batch := authorList[i:end]
			batches = append(batches, []interface{}{repoUrl, batch})
		}
	}

	return batches
}
