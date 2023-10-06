package sync

import "fmt"

func BatchAuthors(repoUrlToAuthorsMap map[string][]string, batchSize int) []map[string][]string {

	var myMap [][]map[string][]string

	for repoUrl, authorList := range repoUrlToAuthorsMap {
		for i := 0; i < len(authorList); i += batchSize {
			itemMap := make(map[string][]string)

			authorBatch := authorList[i:min(i+batchSize, len(authorList))]
			fmt.Println(authorBatch)

			itemMap[repoUrl] = authorBatch
			myMap = append(myMap, itemMap)
		}
	}

	return myMap
}
