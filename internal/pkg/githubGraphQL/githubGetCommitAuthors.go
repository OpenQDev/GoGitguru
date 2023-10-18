package githubGraphQL

import (
	"encoding/json"
	"fmt"
	"io"
	"main/internal/pkg/server"
	"net/http"
)

func GithubGetCommitAuthors(query string, ghAccessToken string, apiCfg server.ApiConfig) (CommitAuthorsResponse, error) {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("token %s", ghAccessToken),
		"Content-Type":  "application/json",
	}

	url := apiCfg.GithubGraphQLBaseUrl

	req, err := createGraphQLRequest(url, query, headers)
	if err != nil {
		return CommitAuthorsResponse{}, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return CommitAuthorsResponse{}, err
	}

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		fmt.Printf("Request failed with status code %d.\n", res.StatusCode)
		return CommitAuthorsResponse{}, fmt.Errorf("request failed with status code %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var jsonData CommitAuthorsResponse
	json.Unmarshal(body, &jsonData)

	return jsonData, nil
}
