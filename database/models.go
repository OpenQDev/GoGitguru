// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"
)

type RepoStatus string

const (
	RepoStatusPending     RepoStatus = "pending"
	RepoStatusQueued      RepoStatus = "queued"
	RepoStatusSyncingRepo RepoStatus = "syncing_repo"
	RepoStatusSynced      RepoStatus = "synced"
	RepoStatusFailed      RepoStatus = "failed"
	RepoStatusNotListed   RepoStatus = "not_listed"
)

func (e *RepoStatus) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = RepoStatus(s)
	case string:
		*e = RepoStatus(s)
	default:
		return fmt.Errorf("unsupported scan type for RepoStatus: %T", src)
	}
	return nil
}

type NullRepoStatus struct {
	RepoStatus RepoStatus `json:"repo_status"`
	Valid      bool       `json:"valid"` // Valid is true if RepoStatus is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullRepoStatus) Scan(value interface{}) error {
	if value == nil {
		ns.RepoStatus, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.RepoStatus.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullRepoStatus) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.RepoStatus), nil
}

type Commit struct {
	CommitHash    string         `json:"commit_hash"`
	Author        sql.NullString `json:"author"`
	AuthorEmail   sql.NullString `json:"author_email"`
	AuthorDate    sql.NullInt64  `json:"author_date"`
	CommitterDate sql.NullInt64  `json:"committer_date"`
	Message       sql.NullString `json:"message"`
	Insertions    sql.NullInt32  `json:"insertions"`
	Deletions     sql.NullInt32  `json:"deletions"`
	LinesChanged  sql.NullInt32  `json:"lines_changed"`
	FilesChanged  sql.NullInt32  `json:"files_changed"`
	RepoUrl       sql.NullString `json:"repo_url"`
}

type Dependency struct {
	InternalID     int32  `json:"internal_id"`
	DependencyName string `json:"dependency_name"`
	DependencyFile string `json:"dependency_file"`
}

type GithubRepo struct {
	InternalID      int32          `json:"internal_id"`
	GithubRestID    int32          `json:"github_rest_id"`
	GithubGraphqlID string         `json:"github_graphql_id"`
	Url             string         `json:"url"`
	Name            string         `json:"name"`
	FullName        string         `json:"full_name"`
	Private         sql.NullBool   `json:"private"`
	OwnerLogin      string         `json:"owner_login"`
	OwnerAvatarUrl  sql.NullString `json:"owner_avatar_url"`
	Description     sql.NullString `json:"description"`
	Homepage        sql.NullString `json:"homepage"`
	Fork            sql.NullBool   `json:"fork"`
	ForksCount      sql.NullInt32  `json:"forks_count"`
	Archived        sql.NullBool   `json:"archived"`
	Disabled        sql.NullBool   `json:"disabled"`
	License         sql.NullString `json:"license"`
	Language        sql.NullString `json:"language"`
	StargazersCount sql.NullInt32  `json:"stargazers_count"`
	WatchersCount   sql.NullInt32  `json:"watchers_count"`
	OpenIssuesCount sql.NullInt32  `json:"open_issues_count"`
	HasIssues       sql.NullBool   `json:"has_issues"`
	HasDiscussions  sql.NullBool   `json:"has_discussions"`
	HasProjects     sql.NullBool   `json:"has_projects"`
	CreatedAt       sql.NullInt32  `json:"created_at"`
	UpdatedAt       sql.NullInt32  `json:"updated_at"`
	PushedAt        sql.NullInt32  `json:"pushed_at"`
	Visibility      sql.NullString `json:"visibility"`
	Size            sql.NullInt32  `json:"size"`
	DefaultBranch   sql.NullString `json:"default_branch"`
}

type GithubUser struct {
	InternalID      int32          `json:"internal_id"`
	GithubRestID    int32          `json:"github_rest_id"`
	GithubGraphqlID string         `json:"github_graphql_id"`
	Login           string         `json:"login"`
	Name            sql.NullString `json:"name"`
	Email           sql.NullString `json:"email"`
	AvatarUrl       sql.NullString `json:"avatar_url"`
	Company         sql.NullString `json:"company"`
	Location        sql.NullString `json:"location"`
	Bio             sql.NullString `json:"bio"`
	Blog            sql.NullString `json:"blog"`
	Hireable        sql.NullBool   `json:"hireable"`
	TwitterUsername sql.NullString `json:"twitter_username"`
	Followers       sql.NullInt32  `json:"followers"`
	Following       sql.NullInt32  `json:"following"`
	Type            string         `json:"type"`
	CreatedAt       sql.NullTime   `json:"created_at"`
	UpdatedAt       sql.NullTime   `json:"updated_at"`
}

type GithubUserRestIDAuthorEmail struct {
	RestID int32  `json:"rest_id"`
	Email  string `json:"email"`
}

type RepoUrl struct {
	Url       string       `json:"url"`
	Status    RepoStatus   `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at"`
}

type ReposToDependency struct {
	Url          string        `json:"url"`
	DependencyID int32         `json:"dependency_id"`
	FirstUseDate sql.NullInt64 `json:"first_use_date"`
	LastUseDate  sql.NullInt64 `json:"last_use_date"`
	UpdatedAt    sql.NullInt64 `json:"updated_at"`
}

type UserToDependency struct {
	DependencyID int32         `json:"dependency_id"`
	UserID       int32         `json:"user_id"`
	FirstUseDate sql.NullInt64 `json:"first_use_date"`
	LastUseDate  sql.NullInt64 `json:"last_use_date"`
	UpdatedAt    sql.NullInt64 `json:"updated_at"`
}
