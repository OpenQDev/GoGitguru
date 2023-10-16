package githubGraphQL

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GithubGetCommitAuthors(query string, ghAccessToken string) (CommitAuthorsResponse, error) {
	// Queries the GitHub GraphQL API.
	// Parameters:
	// - query (str): The GraphQL query to execute.
	// - gh_access_token (str): The GitHub access token to use for authentication.
	// Returns:
	// - dict: The JSON response from the API.
	// - error: If the request failed.

	headers := map[string]string{
		"Authorization": fmt.Sprintf("token %s", ghAccessToken),
		"Content-Type":  "application/json",
	}

	url := "https://api.github.com/graphql"

	payloadObj := GraphQLPayload{
		Query: query,
	}

	payloadBytes, err := json.Marshal(payloadObj)
	if err != nil {
		return CommitAuthorsResponse{}, err
	}

	payload := bytes.NewReader(payloadBytes)

	req, _ := http.NewRequest("POST", url, payload)

	for key, value := range headers {
		req.Header.Add(key, value)
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
