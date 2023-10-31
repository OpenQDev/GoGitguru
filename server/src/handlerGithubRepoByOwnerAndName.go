package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/OpenQDev/GoGitguru/util/githubRest"
	"github.com/OpenQDev/GoGitguru/util/marshaller"

	"github.com/go-chi/chi"
)

type HandlerGithubRepoByOwnerAndNameRequest struct{}
type HandlerGithubRepoByOwnerAndNameResponse = githubRest.GithubRestRepo

func (apiConfig *ApiConfig) HandlerGithubRepoByOwnerAndName(w http.ResponseWriter, r *http.Request) {
	githubAccessToken := r.Header.Get("GH-Authorization")

	if githubAccessToken == "" {
		RespondWithError(w, http.StatusUnauthorized, "You must provide a GitHub access token.")
		return
	}

	owner := chi.URLParam(r, "owner")
	name := chi.URLParam(r, "name")

	fullName := fmt.Sprintf("%s/%s", strings.ToLower(owner), strings.ToLower(name))

	// Check if the repo already exists in the database
	repoExists, err := apiConfig.DB.CheckGithubRepoExists(context.Background(), fullName)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if repoExists {
		repo, err := apiConfig.DB.GetGithubRepo(context.Background(), fullName)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, ConvertGithubDatabaseRepoToGitguruRepo(repo))
		return
	}

	url := fmt.Sprintf("%s/repos/%s/%s", apiConfig.GithubRestAPIBaseUrl, owner, name)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		RespondWithError(w, 500, "failed to create request.")
		return
	}

	req.Header.Add("Authorization", "token "+githubAccessToken)
	resp, err := client.Do(req)
	if err != nil {
		RespondWithError(w, 500, "failed to make request.")
		return
	}

	defer resp.Body.Close()

	var repo githubRest.GithubRestRepo
	marshaller.ReaderToType(resp.Body, &repo)

	params := ConvertGithubRestRepoToInsertGithubRepoParams(repo)

	_, err = apiConfig.DB.InsertGithubRepo(context.Background(), params)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		RespondWithError(w, 500, fmt.Sprintf("failed to insert repo into database: %s", err))
		return
	}

	RespondWithJSON(w, 200, ConvertGithubRestRepoToGitguruRepo(repo))
}

type GitguruRepo struct {
	GithubRestID    int32     `json:"github_rest_id"`
	GithubGraphqlID string    `json:"github_graphql_id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Private         bool      `json:"private"`
	OwnerLogin      string    `json:"owner_login"`
	OwnerAvatarUrl  string    `json:"owner_avatar_url"`
	Description     string    `json:"description"`
	Homepage        string    `json:"homepage"`
	Fork            bool      `json:"fork"`
	ForksCount      int32     `json:"forks_count"`
	Archived        bool      `json:"archived"`
	Disabled        bool      `json:"disabled"`
	License         string    `json:"license"`
	Language        string    `json:"language"`
	StargazersCount int32     `json:"stargazers_count"`
	WatchersCount   int32     `json:"watchers_count"`
	OpenIssuesCount int32     `json:"open_issues_count"`
	HasIssues       bool      `json:"has_issues"`
	HasDiscussions  bool      `json:"has_discussions"`
	HasProjects     bool      `json:"has_projects"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	PushedAt        time.Time `json:"pushed_at"`
	Visibility      string    `json:"visibility"`
	Size            int32     `json:"size"`
	DefaultBranch   string    `json:"default_branch"`
}
