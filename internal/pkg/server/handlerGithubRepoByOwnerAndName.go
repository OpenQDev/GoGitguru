package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
)

type HandlerGithubRepoByOwnerAndNameRequest struct{}
type HandlerGithubRepoByOwnerAndNameResponse = GithubRestRepo

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
	err = json.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		RespondWithError(w, 500, "failed to decode response.")
		return
	}

	// Insert the repo into the database using sqlc generated methods
	params := ConvertGithubRestRepoToInsertGithubRepoParams(repo)

	_, err = apiConfig.DB.InsertGithubRepo(context.Background(), params)
	if err != nil {
		RespondWithError(w, 500, "failed to insert repo into database.")
		return
	}

	RespondWithJSON(w, 200, repo)
}
