package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

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
		RespondWithJSON(w, http.StatusOK, DatabaseRepoToGitguruRepo(repo))
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

	params := RestRepoToDatabaseParams(repo)

	_, err = apiConfig.DB.InsertGithubRepo(context.Background(), params)
	if err != nil && !strings.Contains(err.Error(), "no rows in result set") {
		RespondWithError(w, 500, fmt.Sprintf("failed to insert repo into database: %s", err))
		return
	}

	RespondWithJSON(w, 200, RestRepoToGitguruRepo(repo))
}
