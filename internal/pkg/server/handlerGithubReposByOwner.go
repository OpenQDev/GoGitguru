package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (apiConfig *ApiConfig) HandlerGithubReposByOwner(w http.ResponseWriter, r *http.Request) {
	githubAccessToken := r.Header.Get("GH-Authorization")

	if githubAccessToken == "" {
		RespondWithError(w, 400, "You must provide a GitHub access token.")
		return
	}

	owner := chi.URLParam(r, "owner")

	client := &http.Client{}
	page := 1
	var repos []RestRepo
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/users/%s/repos?per_page=100&page=%d", owner, page), nil)
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

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		_ = string(bodyBytes)

		// Create a new reader with the body bytes for the json decoder
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var restReposResponse []RestRepo
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

	// var insertedRepos []GithubRepo
	// for _, repo := range repos {
	// 	// Convert Repo to InsertGithubRepoParams and insert into database...
	// 	// Refer to the code block from filePath: internal/pkg/server/handlerGithubRepoByOwnerAndName.go
	// 	// startLine: 64
	// 	// endLine: 95
	// 	insertedRepo, err := apiConfig.DB.InsertGithubRepo(context.Background(), params)
	// 	if err != nil {
	// 		RespondWithError(w, 500, "Failed to insert repo into database.")
	// 		return
	// 	}
	// 	insertedRepos = append(insertedRepos, insertedRepo)
	// }

	RespondWithJSON(w, 200, repos)
}
