package gitutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GithubGraphQL(query string, ghAccessToken string) (map[string]interface{}, error) {
	// Queries the GitHub GraphQL API.
	// Parameters:
	// - query (str): The GraphQL query to execute.
	// - gh_access_token (str): The GitHub access token to use for authentication.
	// Returns:
	// - dict: The JSON response from the API.
	// - error: If the request failed.

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", ghAccessToken),
		"Content-Type":  "application/json",
	}

	url := "https://api.github.com/graphql"

	payload := strings.NewReader(fmt.Sprintf(`{"query": "%s"}`, query))
	req, _ := http.NewRequest("POST", url, payload)

	for key, value := range headers {
		req.Header.Add(key, value)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if !(res.StatusCode >= 200 && res.StatusCode < 300) {
		fmt.Printf("Request failed with status code %d.\n", res.StatusCode)
		return nil, fmt.Errorf("request failed with status code %d", res.StatusCode)
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var jsonData map[string]interface{}
	json.Unmarshal(body, &jsonData)
	return jsonData["data"].(map[string]interface{}), nil
}
