package server

import (
	"fmt"
	"main/internal/database"
	"main/internal/pkg/server/util"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type HandlerGithubUserCommitsRequest struct {
	RepoUrls []string `json:"repo_urls"`
	Since    string   `json:"since"`
	Until    string   `json:"until"`
}
type HandlerGithubUserCommitsResponse struct{}

func (apiConfig *ApiConfig) HandlerGithubUserCommits(w http.ResponseWriter, r *http.Request) {
	_ = chi.URLParam(r, "login")

	var body HandlerGithubUserCommitsRequest
	err := util.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		return
	}

	_, err = time.Parse(time.RFC3339, body.Since)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse time for since (%s): %s", body.Since, err))
		return
	}

	_, err = time.Parse(time.RFC3339, body.Until)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse time for until (%s): %s", body.Until, err))
		return
	}

	var commits []database.Commit

	RespondWithJSON(w, http.StatusOK, commits)
}
