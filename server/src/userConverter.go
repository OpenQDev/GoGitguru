package server

import (
	"database/sql"
	"strings"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
)

func ConvertServerUserToInsertUserParams(user User) database.InsertUserParams {
	createdAt, _ := time.Parse(time.RFC3339, user.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, user.UpdatedAt)

	return database.InsertUserParams{
		GithubRestID:    int32(user.GithubRestID),
		GithubGraphqlID: user.GithubGraphqlID,
		Login:           strings.ToLower(user.Login), // we lowercase the login before user exists check, and before insertion
		Name:            sql.NullString{String: user.Name, Valid: true},
		Email:           sql.NullString{String: user.Email, Valid: true},
		AvatarUrl:       sql.NullString{String: user.AvatarURL, Valid: true},
		Company:         sql.NullString{String: user.Company, Valid: true},
		Location:        sql.NullString{String: user.Location, Valid: true},
		Bio:             sql.NullString{String: user.Bio, Valid: true},
		Blog:            sql.NullString{String: user.Blog, Valid: true},
		Hireable:        sql.NullBool{Bool: user.Hireable, Valid: true},
		TwitterUsername: sql.NullString{String: user.TwitterUsername, Valid: true},
		Followers:       sql.NullInt32{Int32: int32(user.Followers), Valid: true},
		Following:       sql.NullInt32{Int32: int32(user.Following), Valid: true},
		Type:            user.Type,
		CreatedAt:       sql.NullTime{Time: createdAt, Valid: true},
		UpdatedAt:       sql.NullTime{Time: updatedAt, Valid: true},
	}
}

func ConvertDatabaseInsertUserParamsToServerUser(params database.GithubUser) User {
	return User{
		InternalID:      int(params.InternalID),
		GithubRestID:    int(params.GithubRestID),
		GithubGraphqlID: params.GithubGraphqlID,
		Login:           params.Login,
		Name:            params.Name.String,
		Email:           params.Email.String,
		AvatarURL:       params.AvatarUrl.String,
		Company:         params.Company.String,
		Location:        params.Location.String,
		Bio:             params.Bio.String,
		Blog:            params.Blog.String,
		Hireable:        params.Hireable.Bool,
		TwitterUsername: params.TwitterUsername.String,
		Followers:       int(params.Followers.Int32),
		Following:       int(params.Following.Int32),
		Type:            params.Type,
		CreatedAt:       params.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:       params.UpdatedAt.Time.Format(time.RFC3339),
	}
}
