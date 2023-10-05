package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"main/internal/database"
	"net/http"
	"time"
)

func (apiConfig *ApiConfig) HandlerRepoCommits(w http.ResponseWriter, r *http.Request) {
	var input struct {
		RepoURL string `json:"repo_url"`
		Since   string `json:"since"`
		Until   string `json:"until"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		RespondWithError(w, 400, "Invalid request body.")
		return
	}

	if input.RepoURL == "" {
		RespondWithError(w, 400, "Missing repo URL.")
		return
	}

	layout := "2006-01-02T15:04:05Z" // ISO 8601 format
	var since, until time.Time

	if input.Since != "" {
		since, err = time.Parse(layout, input.Since)
		if err != nil {
			RespondWithError(w, 400, "Invalid 'since' format.")
			return
		}
	} else {
		since = time.Now().AddDate(0, 0, -30)
	}

	if input.Until != "" {
		until, err = time.Parse(layout, input.Until)
		if err != nil {
			RespondWithError(w, 400, "Invalid 'until' format.")
			return
		}
	} else {
		until = time.Now()
	}

	// Fetch commits from database
	commits, err := apiConfig.DB.GetCommitsWithAuthorInfo(context.Background(), database.GetCommitsWithAuthorInfoParams{
		RepoUrl:      sql.NullString{String: input.RepoURL, Valid: true},
		AuthorDate:   sql.NullInt64{Int64: since.Unix(), Valid: true},
		AuthorDate_2: sql.NullInt64{Int64: until.Unix(), Valid: true},
	})

	if err != nil {
		RespondWithError(w, 500, "Failed to fetch commits.")
		return
	}

	RespondWithJSON(w, 200, commits)
}
