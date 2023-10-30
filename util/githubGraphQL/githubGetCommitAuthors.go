package githubGraphQL

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GithubGetCommitAuthors(query string, ghAccessToken string, githubGraphQLBaseUrl string) (GithubGraphQLCommitAuthorsResponse, error) {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("token %s", ghAccessToken),
		"Content-Type":  "application/json",
	}

	url := githubGraphQLBaseUrl
	fmt.Println(url)

	req, err := createGraphQLRequest(url, query, headers)
	if err != nil {
		return GithubGraphQLCommitAuthorsResponse{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return GithubGraphQLCommitAuthorsResponse{}, err
	}

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		fmt.Printf("Request failed with status code %d.\n", res.StatusCode)
		return GithubGraphQLCommitAuthorsResponse{}, fmt.Errorf("request failed with status code %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var jsonData GithubGraphQLCommitAuthorsResponse
	json.Unmarshal(body, &jsonData)

	return jsonData, nil
}
