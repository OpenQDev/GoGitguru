package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"main/internal/database"
	"main/internal/pkg/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

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
	var repos []RestRepo
	for {
		requestUrl := fmt.Sprintf("%s/users/%s/repos?per_page=100&page=%d", apiConfig.GithubRestAPIBaseUrl, owner, page)
		logger.LogGreenDebug("calling %s", requestUrl)

		req, err := http.NewRequest("GET", requestUrl, nil)
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

		// Create a new reader with the body bytes for the json decoder
		resp = PrintResponseBody(resp)

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

	for _, repo := range repos {

		params := ConvertRestRepoToInsertParams(repo)

		_, err := apiConfig.DB.InsertGithubRepo(context.Background(), params)
		if err != nil {
			RespondWithError(w, 500, fmt.Sprintf("failed to insert repo into database: %s", err))
			return
		}
	}

	RespondWithJSON(w, 200, repos)
}

func ConvertRestRepoToInsertParams(repo RestRepo) database.InsertGithubRepoParams {
	createdAt, _ := time.Parse(time.RFC3339, repo.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, repo.UpdatedAt)
	pushedAt, _ := time.Parse(time.RFC3339, repo.PushedAt)

	return database.InsertGithubRepoParams{
		GithubRestID:    int32(repo.ID),
		GithubGraphqlID: repo.NodeID,
		Url:             repo.URL,
		Name:            repo.Name,
		FullName:        repo.FullName,
		Private:         sql.NullBool{Bool: repo.Private, Valid: true},
		OwnerLogin:      repo.Owner.Login,
		OwnerAvatarUrl:  sql.NullString{String: repo.Owner.AvatarURL, Valid: true},
		Description:     sql.NullString{String: repo.Description, Valid: true},
		Homepage:        sql.NullString{String: repo.Homepage, Valid: repo.Homepage != ""},
		Fork:            sql.NullBool{Bool: repo.Fork, Valid: true},
		ForksCount:      sql.NullInt32{Int32: int32(repo.ForksCount), Valid: true},
		Archived:        sql.NullBool{Bool: repo.Archived, Valid: true},
		Disabled:        sql.NullBool{Bool: repo.Disabled, Valid: true},
		License:         sql.NullString{String: repo.License.Name, Valid: repo.License.Name != ""},
		Language:        sql.NullString{String: repo.Language, Valid: repo.Language != ""},
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
}
