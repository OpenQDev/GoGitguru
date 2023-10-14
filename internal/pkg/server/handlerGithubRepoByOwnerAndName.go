package server

import (
	"context"
	"fmt"
	"main/internal/pkg/server/util"
	"net/http"

	"github.com/go-chi/chi"
)

type HandlerGithubRepoByOwnerAndNameRequest struct{}
type HandlerGithubRepoByOwnerAndNameResponse = GithubRestRepo

func (apiConfig *ApiConfig) HandlerGithubRepoByOwnerAndName(w http.ResponseWriter, r *http.Request) {
	githubAccessToken := r.Header.Get("GH-Authorization")

	if githubAccessToken == "" {
		RespondWithError(w, http.StatusUnauthorized, "You must provide a GitHub access token.")
		return
	}

	owner := chi.URLParam(r, "owner")
	name := chi.URLParam(r, "name")

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

	var repo GithubRestRepo
	util.ReaderToType(resp.Body, &repo)

	params := ConvertGithubRestRepoToInsertGithubRepoParams(repo)

	_, err = apiConfig.DB.InsertGithubRepo(context.Background(), params)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("failed to insert repo into database: %s", err))
		return
	}

	RespondWithJSON(w, 200, repo)
}
