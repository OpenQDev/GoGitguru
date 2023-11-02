package server

import (
	"database/sql"
	"strings"
	"time"

	"github.com/OpenQDev/GoGitguru/database"

	"github.com/OpenQDev/GoGitguru/util/githubRest"
)

type GitguruRepo struct {
	GithubRestID    int32  `json:"github_rest_id"`
	GithubGraphqlID string `json:"github_graphql_id"`
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Private         bool   `json:"private"`
	OwnerLogin      string `json:"owner_login"`
	OwnerAvatarUrl  string `json:"owner_avatar_url"`
	Description     string `json:"description"`
	Homepage        string `json:"homepage"`
	Fork            bool   `json:"fork"`
	ForksCount      int32  `json:"forks_count"`
	Archived        bool   `json:"archived"`
	Disabled        bool   `json:"disabled"`
	License         string `json:"license"`
	Language        string `json:"language"`
	StargazersCount int32  `json:"stargazers_count"`
	WatchersCount   int32  `json:"watchers_count"`
	OpenIssuesCount int32  `json:"open_issues_count"`
	HasIssues       bool   `json:"has_issues"`
	HasDiscussions  bool   `json:"has_discussions"`
	HasProjects     bool   `json:"has_projects"`
	CreatedAt       int32  `json:"created_at"`
	UpdatedAt       int32  `json:"updated_at"`
	PushedAt        int32  `json:"pushed_at"`
	Visibility      string `json:"visibility"`
	Size            int32  `json:"size"`
	DefaultBranch   string `json:"default_branch"`
}

func RestRepoToDatabaseParams(repo githubRest.GithubRestRepo) database.InsertGithubRepoParams {
	createdAt, _ := time.Parse(time.RFC3339, repo.CreatedAt)
	createdAtUnix := createdAt.Unix()
	updatedAt, _ := time.Parse(time.RFC3339, repo.UpdatedAt)
	updatedAtUnix := updatedAt.Unix()
	pushedAt, _ := time.Parse(time.RFC3339, repo.PushedAt)
	pushedAtUnix := pushedAt.Unix()

	return database.InsertGithubRepoParams{
		GithubRestID:    int32(repo.GithubRestID),
		GithubGraphqlID: repo.GithubGraphqlID,
		Url:             repo.URL,
		Name:            repo.Name,
		FullName:        strings.ToLower(repo.FullName),
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
		CreatedAt:       sql.NullInt32{Int32: int32(createdAtUnix), Valid: true},
		UpdatedAt:       sql.NullInt32{Int32: int32(updatedAtUnix), Valid: true},
		PushedAt:        sql.NullInt32{Int32: int32(pushedAtUnix), Valid: true},
		Visibility:      sql.NullString{String: repo.Visibility, Valid: true},
		Size:            sql.NullInt32{Int32: int32(repo.Size), Valid: true},
		DefaultBranch:   sql.NullString{String: repo.DefaultBranch, Valid: true},
	}
}

func RestRepoToGitguruRepo(repo githubRest.GithubRestRepo) GitguruRepo {
	createdAt, _ := time.Parse(time.RFC3339, repo.CreatedAt)
	createdAtUnix := createdAt.Unix()
	updatedAt, _ := time.Parse(time.RFC3339, repo.UpdatedAt)
	updatedAtUnix := updatedAt.Unix()
	pushedAt, _ := time.Parse(time.RFC3339, repo.PushedAt)
	pushedAtUnix := pushedAt.Unix()

	return GitguruRepo{
		GithubRestID:    int32(repo.GithubRestID),
		GithubGraphqlID: repo.GithubGraphqlID,
		Name:            repo.Name,
		FullName:        strings.ToLower(repo.FullName),
		Private:         repo.Private,
		OwnerLogin:      repo.Owner.Login,
		OwnerAvatarUrl:  repo.Owner.AvatarURL,
		Description:     repo.Description,
		Homepage:        repo.Homepage,
		Fork:            repo.Fork,
		ForksCount:      int32(repo.ForksCount),
		Archived:        repo.Archived,
		Disabled:        repo.Disabled,
		License:         repo.License.Name,
		Language:        repo.Language,
		StargazersCount: int32(repo.StargazersCount),
		WatchersCount:   int32(repo.WatchersCount),
		OpenIssuesCount: int32(repo.OpenIssuesCount),
		HasIssues:       repo.HasIssues,
		HasDiscussions:  repo.HasDiscussions,
		HasProjects:     repo.HasProjects,
		CreatedAt:       int32(createdAtUnix),
		UpdatedAt:       int32(updatedAtUnix),
		PushedAt:        int32(pushedAtUnix),
		Visibility:      repo.Visibility,
		Size:            int32(repo.Size),
		DefaultBranch:   repo.DefaultBranch,
	}
}

func DatabaseRepoToGitguruRepo(params database.GithubRepo) GitguruRepo {
	return GitguruRepo{
		GithubRestID:    int32(params.GithubRestID),
		GithubGraphqlID: params.GithubGraphqlID,
		Name:            params.Name,
		FullName:        strings.ToLower(params.FullName),
		Private:         params.Private.Bool,
		OwnerLogin:      params.OwnerLogin,
		OwnerAvatarUrl:  params.OwnerAvatarUrl.String,
		Description:     params.Description.String,
		Homepage:        params.Homepage.String,
		Fork:            params.Fork.Bool,
		ForksCount:      int32(params.ForksCount.Int32),
		Archived:        params.Archived.Bool,
		Disabled:        params.Disabled.Bool,
		License:         params.License.String,
		Language:        params.Language.String,
		StargazersCount: int32(params.StargazersCount.Int32),
		WatchersCount:   int32(params.WatchersCount.Int32),
		OpenIssuesCount: int32(params.OpenIssuesCount.Int32),
		HasIssues:       params.HasIssues.Bool,
		HasDiscussions:  params.HasDiscussions.Bool,
		HasProjects:     params.HasProjects.Bool,
		CreatedAt:       params.CreatedAt.Int32,
		UpdatedAt:       params.UpdatedAt.Int32,
		PushedAt:        params.PushedAt.Int32,
		Visibility:      params.Visibility.String,
		Size:            int32(params.Size.Int32),
		DefaultBranch:   params.DefaultBranch.String,
	}
}
