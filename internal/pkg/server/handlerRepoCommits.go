package server

import (
	"context"
	"database/sql"
	"fmt"
	"main/internal/database"
	"main/internal/pkg/server/util"
	"net/http"
	"time"
)

type HandlerRepoCommitsRequest struct {
	RepoURL string `json:"repo_url"`
	Since   string `json:"since"`
	Until   string `json:"until"`
}

type HandlerRepoCommitsResponse struct{}

func (apiConfig *ApiConfig) HandlerRepoCommits(w http.ResponseWriter, r *http.Request) {
	var handlerRepoCommitsRequest HandlerRepoCommitsRequest
	err := util.ReaderToType(r.Body, &handlerRepoCommitsRequest)

	if err != nil {
		RespondWithError(w, 400, "Invalid request body.")
		return
	}

	if handlerRepoCommitsRequest.RepoURL == "" {
		RespondWithError(w, 400, "Missing repo URL.")
		return
	}

	since, err := time.Parse(time.RFC3339, handlerRepoCommitsRequest.Since)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse time for since (%s): %s", handlerRepoCommitsRequest.Since, err))
		return
	}

	until, err := time.Parse(time.RFC3339, handlerRepoCommitsRequest.Until)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse time for until (%s): %s", handlerRepoCommitsRequest.Until, err))
		return
	}

	commitsWithAuthorInfo, err := apiConfig.DB.GetCommitsWithAuthorInfo(context.Background(), database.GetCommitsWithAuthorInfoParams{
		RepoUrl:      sql.NullString{String: handlerRepoCommitsRequest.RepoURL, Valid: true},
		AuthorDate:   sql.NullInt64{Int64: since.Unix(), Valid: true},
		AuthorDate_2: sql.NullInt64{Int64: until.Unix(), Valid: true},
	})

	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Failed to fetch commits: %s", err))
		return
	}

	RespondWithJSON(w, 200, commitsWithAuthorInfo)
}
