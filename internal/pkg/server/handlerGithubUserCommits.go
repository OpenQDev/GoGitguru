package server

import (
	"context"
	"database/sql"
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

type HandlerGithubUserCommitsResponse = []CommitWithAuthorInfo

func (apiConfig *ApiConfig) HandlerGithubUserCommits(w http.ResponseWriter, r *http.Request) {
	githubAccessToken := r.Header.Get("GH-Authorization")

	if githubAccessToken == "" {
		RespondWithError(w, http.StatusUnauthorized, "You must provide a GitHub access token.")
		return
	}

	login := chi.URLParam(r, "login")

	var body HandlerGithubUserCommitsRequest
	err := util.ReaderToType(r.Body, &body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to read body of request: %s", err))
		return
	}

	since, err := time.Parse(time.RFC3339, body.Since)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse time for since (%s): %s", body.Since, err))
		return
	}

	until, err := time.Parse(time.RFC3339, body.Until)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("failed to parse time for until (%s): %s", body.Until, err))
		return
	}

	var commits []database.GetAllUserCommitsRow

	params := database.GetAllUserCommitsParams{
		AuthorDate:   sql.NullInt64{Int64: since.Unix(), Valid: true},
		AuthorDate_2: sql.NullInt64{Int64: until.Unix(), Valid: true},
		Login:        login,
	}

	commits, err = apiConfig.DB.GetAllUserCommits(context.Background(), params)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("failed to fetch commits from database: %s", err))
		return
	}

	RespondWithJSON(w, http.StatusOK, ConvertGetAllUserCommitsRowToAuthorsToCommitWithAuthorInfo(commits))
}

func ConvertGetAllUserCommitsRowToAuthorsToCommitWithAuthorInfo(rows []database.GetAllUserCommitsRow) []CommitWithAuthorInfo {
	var commits []CommitWithAuthorInfo
	for _, row := range rows {
		commit := CommitWithAuthorInfo{
			CommitHash:      row.CommitHash,
			Author:          row.Author.String,
			AuthorEmail:     row.AuthorEmail.String,
			AuthorDate:      row.AuthorDate.Int64,
			CommitterDate:   row.CommitterDate.Int64,
			Message:         row.Message.String,
			Insertions:      row.Insertions.Int32,
			Deletions:       row.Deletions.Int32,
			LinesChanged:    row.LinesChanged.Int32,
			FilesChanged:    row.FilesChanged.Int32,
			RepoUrl:         row.RepoUrl.String,
			RestID:          row.RestID,
			Email:           row.Email,
			InternalID:      row.InternalID,
			GithubRestID:    row.GithubRestID,
			GithubGraphqlID: row.GithubGraphqlID,
			Login:           row.Login,
			Name:            row.Name.String,
			Email_2:         row.Email_2.String,
			AvatarUrl:       row.AvatarUrl.String,
			Company:         row.Company.String,
			Location:        row.Location.String,
			Bio:             row.Bio.String,
			Blog:            row.Blog.String,
			Hireable:        row.Hireable.Bool,
			TwitterUsername: row.TwitterUsername.String,
			Followers:       row.Followers.Int32,
			Following:       row.Following.Int32,
			Type:            row.Type,
			CreatedAt:       row.CreatedAt.Time,
			UpdatedAt:       row.UpdatedAt.Time,
		}
		commits = append(commits, commit)
	}
	return commits
}
