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

type HandlerGithubUserCommitsResponse = struct{}

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

	var commits []database.GetUserCommitsForReposRow
	// if len(body.RepoUrls) > 0 {
	// 	repoUrls := "{" + strings.Join(body.RepoUrls, ",") + "}"

	// 	params := database.GetUserCommitsForReposParams{
	// 		AuthorDate:   sql.NullInt64{Int64: since.Unix(), Valid: true},
	// 		AuthorDate_2: sql.NullInt64{Int64: until.Unix(), Valid: true},
	// 		Login:        login,
	// 		RepoUrl:      sql.NullString{String: repoUrls, Valid: true},
	// 	}
	// 	rawCommits, err := apiConfig.DB.GetUserCommitsForRepos(context.Background(), params)
	// 	if err != nil {
	// 		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to call GetUserCommitsForRepos: %s", err))
	// 		return
	// 	}

	// 	commits, ok := rawCommits.([]database.GetUserCommitsForReposRow)
	// 	if !ok {

	// 	}
	// } else {
	// 	params := database.GetAllUserCommitsParams{
	// 		AuthorDate:   sql.NullInt64{Int64: since.Unix(), Valid: true},
	// 		AuthorDate_2: sql.NullInt64{Int64: until.Unix(), Valid: true},
	// 		Login:        login,
	// 	}
	// 	commits, err = apiConfig.DB.GetAllUserCommits(context.Background(), params)
	// }

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to fetch commits from database: %s", err))
		return
	}

	RespondWithJSON(w, http.StatusOK, commits)
}

func ConvertUserCommitsForReposToAllUserCommits(input database.GetUserCommitsForReposRow) database.GetAllUserCommitsRow {
	return database.GetAllUserCommitsRow{
		CommitHash:      input.CommitHash,
		Author:          input.Author,
		AuthorEmail:     input.AuthorEmail,
		AuthorDate:      input.AuthorDate,
		CommitterDate:   input.CommitterDate,
		Message:         input.Message,
		Insertions:      input.Insertions,
		Deletions:       input.Deletions,
		LinesChanged:    input.LinesChanged,
		FilesChanged:    input.FilesChanged,
		RepoUrl:         input.RepoUrl,
		RestID:          input.RestID,
		Email:           input.Email,
		InternalID:      input.InternalID,
		GithubRestID:    input.GithubRestID,
		GithubGraphqlID: input.GithubGraphqlID,
		Login:           input.Login,
		Name:            input.Name,
		Email_2:         input.Email_2,
		AvatarUrl:       input.AvatarUrl,
		Company:         input.Company,
		Location:        input.Location,
		Bio:             input.Bio,
		Blog:            input.Blog,
		Hireable:        input.Hireable,
		TwitterUsername: input.TwitterUsername,
		Followers:       input.Followers,
		Following:       input.Following,
		Type:            input.Type,
		CreatedAt:       input.CreatedAt,
		UpdatedAt:       input.UpdatedAt,
	}
}
