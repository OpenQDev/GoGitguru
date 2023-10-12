package server

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/pkg/logger"
	"net/http"

	"github.com/go-chi/chi"
)

type HandlerGithubReposByOwnerRequest struct{}
type HandlerGithubReposByOwnerResponse struct{}

func (apiConfig *ApiConfig) HandlerGithubReposByOwner(w http.ResponseWriter, r *http.Request) {
	githubAccessToken := r.Header.Get("GH-Authorization")

	if githubAccessToken == "" {
		RespondWithError(w, 401, "You must provide a GitHub access token.")
		return
	}

	owner := chi.URLParam(r, "owner")
	logger.LogGreenDebug("getting repos for owner: %s", owner)

	client := &http.Client{}
	page := 1
	var repos []GithubRestRepo
	for {
		requestUrl := fmt.Sprintf("%s/users/%s/repos?per_page=100&page=%d", apiConfig.GithubRestAPIBaseUrl, owner, page)
		logger.LogGreenDebug("calling %s", requestUrl)

		req, err := http.NewRequest("GET", requestUrl, nil)
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

		// Create a new reader with the body bytes for the json decoder
		resp = PrintResponseBody(resp)

		var restReposResponse []GithubRestRepo
		err = json.NewDecoder(resp.Body).Decode(&restReposResponse)
		if err != nil {
			RespondWithError(w, 500, "Failed to decode response.")
			return
		}

		repos = append(repos, restReposResponse...)
		if len(restReposResponse) < 100 {
			break
		}

		page++
	}

	for _, repo := range repos {

		params := ConvertRestRepoToInsertParams(repo)

		_, err := apiConfig.DB.InsertGithubRepo(context.Background(), params)
		if err != nil {
			RespondWithError(w, 500, fmt.Sprintf("failed to insert repo into database: %s", err))
			return
		}
	}

	RespondWithJSON(w, 200, repos)
}
