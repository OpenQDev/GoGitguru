package githubRest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GithubRestAPI(endpoint string, ghAccessToken string) (map[string]interface{}, error) {
	// Queries the GitHub REST API.
	// Parameters:
	// - endpoint (str): The REST API endpoint to hit.
	// - gh_access_token (str): The GitHub access token to use for authentication.
	// Returns:
	// - dict: The JSON response from the API.
	// - error: If the request failed.

	headers := map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", ghAccessToken),
		"Accept":        "application/vnd.github.v3+json",
	}

	url := fmt.Sprintf("https://api.github.com/%s", endpoint)

	req, _ := http.NewRequest("GET", url, nil)

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
	return jsonData, nil
}
