package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Repo struct {
	GithubRestID    int    `json:"github_rest_id"`
	GithubGraphqlID string `json:"github_graphql_id"`
	URL             string `json:"url"`
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Private         bool   `json:"private"`
	OwnerLogin      string `json:"owner_login"`
	OwnerAvatarURL  string `json:"owner_avatar_url"`
	Description     string `json:"description"`
	Homepage        string `json:"homepage"`
	Fork            bool   `json:"fork"`
	ForksCount      int    `json:"forks_count"`
	Archived        bool   `json:"archived"`
	Disabled        bool   `json:"disabled"`
	License         string `json:"license"`
	Language        string `json:"language"`
	StargazersCount int    `json:"stargazers_count"`
	WatchersCount   int    `json:"watchers_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
	HasIssues       bool   `json:"has_issues"`
	HasDiscussions  bool   `json:"has_discussions"`
	HasProjects     bool   `json:"has_projects"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	PushedAt        string `json:"pushed_at"`
	Visibility      string `json:"visibility"`
	Size            int    `json:"size"`
	DefaultBranch   string `json:"default_branch"`
}

func (apiCfg *ApiConfig) HandlerGithubReposByOwner(w http.ResponseWriter, r *http.Request) {
	owner := strings.TrimPrefix(r.URL.Path, "/repos/github/")
	githubAccessToken := r.Header.Get("GH-Authorization")

	if githubAccessToken == "" {
		RespondWithError(w, 400, "You must provide a GitHub access token.")
		return
	}

	var repos []Repo
	page := 1

	for {
		requestURL := fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&page=%d", owner, page)
		req, _ := http.NewRequest("GET", requestURL, nil)
		req.Header.Set("Authorization", "Bearer "+githubAccessToken)

		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.StatusCode != 200 {
			RespondWithError(w, 500, "Error fetching data from GitHub.")
			return
		}

		body, _ := io.ReadAll(resp.Body)
		var jsonRepos []Repo
		json.Unmarshal(body, &jsonRepos)

		repos = append(repos, jsonRepos...)

		if len(jsonRepos) < 100 {
			break
		} else {
			page++
		}
	}

	// TODO: Insert data into the database

	RespondWithJSON(w, 200, repos)
}
