package server

import (
	"encoding/json"
	"main/internal/database"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type RequestBody struct {
	RepoUrls []string `json:"repo_urls"`
	Since    string   `json:"since"`
	Until    string   `json:"until"`
}

func (apiConfig *ApiConfig) HandlerGithubUserCommits(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "login")

	var body RequestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = time.Parse(time.RFC3339, body.Since)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = time.Parse(time.RFC3339, body.Until)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var commits []database.Commit

	RespondWithJSON(w, 200, commits)
}
