package server

import (
	"context"
	"fmt"
	"main/internal/pkg/logger"
	"main/internal/pkg/server/util"
	"net/http"

	"github.com/go-chi/chi"
)

type HandlerGithubReposByOwnerRequest struct{}
type HandlerGithubReposByOwnerResponse = []GithubRestRepo

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
			RespondWithError(w, 500, fmt.Sprintf("Failed to create request: %s", err))
			return
		}

		req.Header.Add("Authorization", "token "+githubAccessToken)
		resp, err := client.Do(req)
		if err != nil {
			RespondWithError(w, 500, fmt.Sprintf("Failed to make request %s: %s", requestUrl, err))
			return
		}

		var restReposResponse []GithubRestRepo
		err = util.ReaderToType(resp.Body, &restReposResponse)
		if err != nil {
			RespondWithError(w, 500, fmt.Sprintf("Failed to decode response from %s to []GithubRestRepo: %s", requestUrl, err))
			return
		}

		repos = append(repos, restReposResponse...)
		if len(restReposResponse) < 100 {
			break
		}

		page++
	}

	for _, repo := range repos {

		params := ConvertGithubRestRepoToInsertGithubRepoParams(repo)

		_, err := apiConfig.DB.InsertGithubRepo(context.Background(), params)
		if err != nil {
			RespondWithError(w, 500, fmt.Sprintf("failed to insert repo into database: %s", err))
			return
		}
	}

	RespondWithJSON(w, 200, repos)
}
