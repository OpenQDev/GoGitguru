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

type HandlerRepoAuthorsRequest struct {
	RepoUrls []string `json:"repo_urls"`
	Since    string   `json:"since"`
	Until    string   `json:"until"`
}

type HandlerRepoAuthorsResponse struct{}

func (apiConfig *ApiConfig) HandlerRepoAuthors(w http.ResponseWriter, r *http.Request) {

	var body HandlerRepoAuthorsRequest
	err := marshaller.ReaderToType(r.Body, &body)
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

	params := database.GetRepoAuthorsInfoParams{
		RepoUrls:     body.RepoUrls,
		AuthorDate:   sql.NullInt64{Int64: since.Unix(), Valid: true},
		AuthorDate_2: sql.NullInt64{Int64: until.Unix(), Valid: true},
	}

	repoAuthorsInfo, err := apiConfig.DB.GetRepoAuthorsInfo(context.Background(), params)

	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Failed to fetch repo authors: %s", err))
		return
	}

	RespondWithJSON(w, 200, ConvertToAuthorInfo(repoAuthorsInfo))
}

func ConvertToAuthorInfo(rows []database.GetRepoAuthorsInfoRow) []AuthorInfo {
	var authors []AuthorInfo
	for _, row := range rows {
		author := AuthorInfo{
			Author:          row.Author.String,
			AuthorEmail:     row.AuthorEmail.String,
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
		}
		authors = append(authors, author)
	}
	return authors
}

type AuthorInfo struct {
	Author          string `json:"author"`
	AuthorEmail     string `json:"author_email"`
	RestID          int32  `json:"rest_id"`
	Email           string `json:"email"`
	InternalID      int32  `json:"internal_id"`
	GithubRestID    int32  `json:"github_rest_id"`
	GithubGraphqlID string `json:"github_graphql_id"`
	Login           string `json:"login"`
	Name            string `json:"name"`
	Email_2         string `json:"email_2"`
	AvatarUrl       string `json:"avatar_url"`
	Company         string `json:"company"`
	Location        string `json:"location"`
	Bio             string `json:"bio"`
	Blog            string `json:"blog"`
	Hireable        bool   `json:"hireable"`
	TwitterUsername string `json:"twitter_username"`
	Followers       int32  `json:"followers"`
	Following       int32  `json:"following"`
	Type            string `json:"type"`
}
