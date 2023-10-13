package server

import (
	"database/sql"
	"main/internal/database"
	"time"
)

func ConvertGithubRestRepoToInsertGithubRepoParams(repo GithubRestRepo) database.InsertGithubRepoParams {
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
