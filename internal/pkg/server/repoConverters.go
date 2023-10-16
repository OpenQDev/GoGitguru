package server

import (
	"database/sql"
	"main/internal/database"
	"main/internal/pkg/githubRest"
	"time"
)

func ConvertGithubRestRepoToInsertGithubRepoParams(repo githubRest.GithubRestRepo) database.InsertGithubRepoParams {
	createdAt, _ := time.Parse(time.RFC3339, repo.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, repo.UpdatedAt)
	pushedAt, _ := time.Parse(time.RFC3339, repo.PushedAt)

	return database.InsertGithubRepoParams{
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
		License:         sql.NullString{String: repo.License.Name, Valid: true},
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
}

func ConvertDatabaseGithubRepoToGithubRestRepo(params database.GithubRepo) githubRest.GithubRestRepo {
	return githubRest.GithubRestRepo{
		GithubRestID:    int(params.GithubRestID),
		GithubGraphqlID: params.GithubGraphqlID,
		URL:             params.Url,
		Name:            params.Name,
		FullName:        params.FullName,
		Private:         params.Private.Bool,
		Owner: struct {
			Login      string `json:"login"`
			ID         int    `json:"id"`
			NodeID     string `json:"node_id"`
			AvatarURL  string `json:"avatar_url"`
			GravatarID string `json:"gravatar_id"`
			URL        string `json:"url"`
		}{
			Login:     params.OwnerLogin,
			AvatarURL: params.OwnerAvatarUrl.String,
		},
		Description: params.Description.String,
		Homepage:    params.Homepage.String,
		Fork:        params.Fork.Bool,
		ForksCount:  int(params.ForksCount.Int32),
		Archived:    params.Archived.Bool,
		Disabled:    params.Disabled.Bool,
		License: struct {
			Key  string `json:"key"`
			Name string `json:"name"`
		}{
			Name: params.License.String,
		},
		Language:        params.Language.String,
		StargazersCount: int(params.StargazersCount.Int32),
		WatchersCount:   int(params.WatchersCount.Int32),
		OpenIssuesCount: int(params.OpenIssuesCount.Int32),
		HasIssues:       params.HasIssues.Bool,
		HasDiscussions:  params.HasDiscussions.Bool,
		HasProjects:     params.HasProjects.Bool,
		CreatedAt:       params.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:       params.UpdatedAt.Time.Format(time.RFC3339),
		PushedAt:        params.PushedAt.Time.Format(time.RFC3339),
		Visibility:      params.Visibility.String,
		Size:            int(params.Size.Int32),
		DefaultBranch:   params.DefaultBranch.String,
	}
}
