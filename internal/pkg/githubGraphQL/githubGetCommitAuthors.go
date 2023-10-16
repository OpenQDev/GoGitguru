package githubGraphQL

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GithubGetCommitAuthors(query string, ghAccessToken string) (CommitAuthorsResponse, error) {
	headers := map[string]string{
		"Authorization": fmt.Sprintf("token %s", ghAccessToken),
		"Content-Type":  "application/json",
	}

	url := "https://api.github.com/graphql"

	req, err := createRequest(url, query, headers)
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

func createRequest(url string, query string, headers map[string]string) (*http.Request, error) {
	payloadObj := GraphQLPayload{Query: query}
	payloadBytes, err := json.Marshal(payloadObj)
	if err != nil {
		return nil, err
	}

	payload := bytes.NewReader(payloadBytes)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	return req, nil
}
