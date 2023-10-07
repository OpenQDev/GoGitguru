package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"main/internal/database"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

func (apiConfig *ApiConfig) HandlerGithubReposByOwner(w http.ResponseWriter, r *http.Request) {
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
		RespondWithError(w, 500, "Failed to create request.")
		return
	}

	req.Header.Add("Authorization", "token "+githubAccessToken)
	resp, err := client.Do(req)
	if err != nil {
		RespondWithError(w, 500, "Failed to make request.")
		return
	}

	defer resp.Body.Close()

	var repo Repo
	err = json.NewDecoder(resp.Body).Decode(&repo)
	if err != nil {
		RespondWithError(w, 500, "Failed to decode response.")
		return
	}

	layout := "2006-01-02T15:04:05Z" // ISO 8601 format
	createdAt, err := time.Parse(layout, repo.CreatedAt)
	if err != nil {
		RespondWithError(w, 500, "Failed to parse CreatedAt.")
		return
	}
	updatedAt, err := time.Parse(layout, repo.UpdatedAt)
	if err != nil {
		RespondWithError(w, 500, "Failed to parse UpdatedAt.")
		return
	}
	pushedAt, err := time.Parse(layout, repo.PushedAt)
	if err != nil {
		RespondWithError(w, 500, "Failed to parse PushedAt.")
		return
	}

	// Insert the repo into the database using sqlc generated methods
	params := database.InsertGithubRepoParams{
		GithubRestID:    int32(repo.GithubRestID),
		GithubGraphqlID: repo.GithubGraphqlID,
		Url:             repo.URL,
		Name:            repo.Name,
		FullName:        repo.FullName,
		Private:         sql.NullBool{Bool: repo.Private, Valid: true},
		OwnerLogin:      repo.Owner.Login,
		OwnerAvatarUrl:  sql.NullString{String: repo.Owner.AvatarURL, Valid: true},
		Description:     sql.NullString{String: repo.Description, Valid: true},
		Homepage:        sql.NullString{String: repo.Homepage, Valid: true},
		Fork:            sql.NullBool{Bool: repo.Fork, Valid: true},
		ForksCount:      sql.NullInt32{Int32: int32(repo.ForksCount), Valid: true},
		Archived:        sql.NullBool{Bool: repo.Archived, Valid: true},
		Disabled:        sql.NullBool{Bool: repo.Disabled, Valid: true},
		License:         sql.NullString{String: repo.License, Valid: true},
		Language:        sql.NullString{String: repo.Language, Valid: true},
		StargazersCount: sql.NullInt32{Int32: int32(repo.StargazersCount), Valid: true},
		WatchersCount:   sql.NullInt32{Int32: int32(repo.WatchersCount), Valid: true},
		OpenIssuesCount: sql.NullInt32{Int32: int32(repo.OpenIssuesCount), Valid: true},
		HasIssues:       sql.NullBool{Bool: repo.HasIssues, Valid: true},
		HasDiscussions:  sql.NullBool{Bool: repo.HasDiscussions, Valid: true},
		HasProjects:     sql.NullBool{Bool: repo.HasProjects, Valid: true},
		CreatedAt:       sql.NullTime{Time: createdAt, Valid: true},
		UpdatedAt:       sql.NullTime{Time: updatedAt, Valid: true},
		PushedAt:        sql.NullTime{Time: pushedAt, Valid: true},
		Visibility:      sql.NullString{String: repo.Visibility, Valid: true},
		Size:            sql.NullInt32{Int32: int32(repo.Size), Valid: true},
		DefaultBranch:   sql.NullString{String: repo.DefaultBranch, Valid: true},
	}

	_, err = apiConfig.DB.InsertGithubRepo(context.Background(), params)
	if err != nil {
		RespondWithError(w, 500, "Failed to insert repo into database.")
		return
	}

	RespondWithJSON(w, 200, repo)
}
