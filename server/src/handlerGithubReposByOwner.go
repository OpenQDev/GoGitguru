package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/OpenQDev/GoGitguru/util/githubRest"
	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/OpenQDev/GoGitguru/util/marshaller"

	"github.com/go-chi/chi"
)

type HandlerGithubReposByOwnerRequest struct{}
type HandlerGithubReposByOwnerResponse = []githubRest.GithubRestRepo

func (apiConfig *ApiConfig) HandlerGithubReposByOwner(w http.ResponseWriter, r *http.Request) {
	githubAccessToken := r.Header.Get("GH-Authorization")

	if githubAccessToken == "" {
		RespondWithError(w, http.StatusUnauthorized, "You must provide a GitHub access token.")
		return
	}

	owner := chi.URLParam(r, "owner")
	logger.LogGreenDebug("getting repos for owner: %s", owner)

	client := &http.Client{}
	page := 1
	var repos []githubRest.GithubRestRepo
	for {
		requestUrl := fmt.Sprintf("%s/users/%s/repos?per_page=100&page=%d", apiConfig.GithubRestAPIBaseUrl, owner, page)
		logger.LogGreenDebug("calling %s", requestUrl)

		req, err := http.NewRequest("GET", requestUrl, nil)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create request: %s", err))
			return
		}

		req.Header.Add("Authorization", "token "+githubAccessToken)
		resp, err := client.Do(req)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to make request %s: %s", requestUrl, err))
			return
		}

		if resp.StatusCode == http.StatusNotFound {
			RespondWithError(w, http.StatusNotFound, "GitHub owner not found.")
			return
		}

		var restReposResponse []githubRest.GithubRestRepo
		err = marshaller.ReaderToType(resp.Body, &restReposResponse)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to decode response from %s to []GithubRestRepo: %s", requestUrl, err))
			return
		}

		repos = append(repos, restReposResponse...)
		if len(restReposResponse) < 100 {
			break
		}

		page++
	}

	responseRepos := []GitguruRepo{}
	for _, repo := range repos {
		logger.LogBlue("inserting repo %s", repo.Name)

		params := RestRepoToDatabaseParams(repo)

		_, err := apiConfig.DB.InsertGithubRepo(context.Background(), params)
		if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
			fmt.Println("IN DB ERROR", err)
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to insert repo into database: %s", err))
			return
		}

		responseRepos = append(responseRepos, RestRepoToGitguruRepo(repo))
	}

	RespondWithJSON(w, 200, responseRepos)
}
