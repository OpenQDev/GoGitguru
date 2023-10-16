package gitutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Author struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	User  struct {
		GithubRestID    int    `json:"github_rest_id"`
		GithubGraphqlID string `json:"github_graphql_id"`
		Login           string `json:"login"`
		Name            string `json:"name"`
		Email           string `json:"email"`
		AvatarURL       string `json:"avatar_url"`
		Company         string `json:"company"`
		Location        string `json:"location"`
		Hireable        bool   `json:"hireable"`
		Bio             string `json:"bio"`
		Blog            string `json:"blog"`
		TwitterUsername string `json:"twitter_username"`
		Followers       struct {
			TotalCount int `json:"totalCount"`
		} `json:"followers"`
		Following struct {
			TotalCount int `json:"totalCount"`
		} `json:"following"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"user"`
}

type Commit struct {
	Author Author `json:"author"`
}

type CommitAuthorsResponse struct {
	Data *struct {
		RateLimit struct {
			Limit   int    `json:"limit"`
			Used    int    `json:"used"`
			ResetAt string `json:"resetAt"`
		} `json:"rateLimit"`
		Repository map[string]Commit `json:"repository"`
	} `json:"data"`
	Errors *[]struct {
		Path       []string `json:"path"`
		Extensions struct {
			Code         string `json:"code"`
			TypeName     string `json:"typeName"`
			ArgumentName string `json:"argumentName"`
		} `json:"extensions"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Message string `json:"message"`
	} `json:"errors"`
}

type GraphQLPayload struct {
	Query string `json:"query"`
}

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
