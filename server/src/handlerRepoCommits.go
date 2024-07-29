package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/marshaller"
)

type HandlerRepoCommitsRequest struct {
	RepoURL string `json:"repo_url"`
	Since   string `json:"since"`
	Until   string `json:"until"`
}

type HandlerRepoCommitsResponse struct{}

func (apiConfig *ApiConfig) HandlerRepoCommits(w http.ResponseWriter, r *http.Request) {
	var handlerRepoCommitsRequest HandlerRepoCommitsRequest
	err := marshaller.ReaderToType(r.Body, &handlerRepoCommitsRequest)

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

	params := database.GetCommitsWithAuthorInfoParams{
		RepoUrl:      sql.NullString{String: handlerRepoCommitsRequest.RepoURL, Valid: true},
		AuthorDate:   sql.NullInt64{Int64: since.Unix(), Valid: true},
		AuthorDate_2: sql.NullInt64{Int64: until.Unix(), Valid: true},
	}

	commitsWithAuthorInfo, err := apiConfig.DB.GetCommitsWithAuthorInfo(context.Background(), params)

	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Failed to fetch commits: %s", err))
		return
	}

	RespondWithJSON(w, 200, ConvertGetCommitsWithAuthorInfoRowToAuthorsToCommitWithAuthorInfo(commitsWithAuthorInfo))
}

func ConvertGetCommitsWithAuthorInfoRowToAuthorsToCommitWithAuthorInfo(rows []database.GetCommitsWithAuthorInfoRow) []CommitWithAuthorInfo {
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
			HasCheckedUser:  row.HasCheckedUser.Bool,
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

type CommitWithAuthorInfo struct {
	CommitHash      string    `json:"commit_hash"`
	Author          string    `json:"author"`
	AuthorEmail     string    `json:"author_email"`
	AuthorDate      int64     `json:"author_date"`
	CommitterDate   int64     `json:"committer_date"`
	Message         string    `json:"message"`
	Insertions      int32     `json:"insertions"`
	Deletions       int32     `json:"deletions"`
	LinesChanged    int32     `json:"lines_changed"`
	FilesChanged    int32     `json:"files_changed"`
	RepoUrl         string    `json:"repo_url"`
	HasCheckedUser  bool      `json:"hasCheckedUser"`
	RestID          int32     `json:"rest_id"`
	Email           string    `json:"email"`
	InternalID      int32     `json:"internal_id"`
	GithubRestID    int32     `json:"github_rest_id"`
	GithubGraphqlID string    `json:"github_graphql_id"`
	Login           string    `json:"login"`
	Name            string    `json:"name"`
	Email_2         string    `json:"email_2"`
	AvatarUrl       string    `json:"avatar_url"`
	Company         string    `json:"company"`
	Location        string    `json:"location"`
	Bio             string    `json:"bio"`
	Blog            string    `json:"blog"`
	Hireable        bool      `json:"hireable"`
	TwitterUsername string    `json:"twitter_username"`
	Followers       int32     `json:"followers"`
	Following       int32     `json:"following"`
	Type            string    `json:"type"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
