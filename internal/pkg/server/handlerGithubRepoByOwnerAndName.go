package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type Repo struct {
	GithubRestID    int    `json:"id"`
	GithubGraphqlID string `json:"node_id"`
	URL             string `json:"url"`
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Private         bool   `json:"private"`
	Owner           struct {
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
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

func (apiConfig *ApiConfig) HandlerGithubRepoByOwnerAndName(w http.ResponseWriter, r *http.Request) {
	githubAccessToken := r.Header.Get("GH-Authorization")

	if githubAccessToken == "" {
		RespondWithError(w, 400, "You must provide a GitHub access token.")
		return
	}

	owner := chi.URLParam(r, "owner")
	name := chi.URLParam(r, "name")

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+owner+"/"+name, nil)
	if err != nil {
		RespondWithError(w, 500, "Failed to create request.")
		return
	}

	req.Header.Add("Authorization", "token "+githubAccessToken)
	resp, err := client.Do(req)
	if err != nil {
		RespondWithError(w, 500, "Failed to make request.")
		return
	}

	defer resp.Body.Close()

	var repo Repo
	err = json.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		RespondWithError(w, 500, "Failed to decode response.")
		return
	}

	RespondWithJSON(w, 200, repo)
}
