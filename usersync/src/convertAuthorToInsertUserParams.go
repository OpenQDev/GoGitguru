package usersync

import (
	"database/sql"
	"strings"
	"time"

	"github.com/OpenQDev/GoGitguru/database"
)

func convertAuthorToInsertUserParams(author GithubGraphQLAuthor, createdAt time.Time, updatedAt time.Time) database.InsertUserParams {
	authorParams := database.InsertUserParams{
		GithubRestID:    int32(author.User.GithubRestID),
		GithubGraphqlID: author.User.GithubGraphqlID,
		Login:           strings.ToLower(author.User.Login),
		Name:            sql.NullString{String: author.User.Name, Valid: true},
		Email:           sql.NullString{String: author.User.Email, Valid: true},
		AvatarUrl:       sql.NullString{String: author.User.AvatarURL, Valid: true},
		Company:         sql.NullString{String: author.User.Company, Valid: true},
		Location:        sql.NullString{String: author.User.Location, Valid: true},
		Hireable:        sql.NullBool{Bool: author.User.Hireable, Valid: true},
		Bio:             sql.NullString{String: author.User.Bio, Valid: true},
		Blog:            sql.NullString{String: author.User.Blog, Valid: true},
		TwitterUsername: sql.NullString{String: author.User.TwitterUsername, Valid: true},
		Followers:       sql.NullInt32{Int32: int32(author.User.Followers.TotalCount), Valid: true},
		Following:       sql.NullInt32{Int32: int32(author.User.Following.TotalCount), Valid: true},
		Type:            "User",
		CreatedAt:       sql.NullTime{Time: createdAt, Valid: true},
		UpdatedAt:       sql.NullTime{Time: updatedAt, Valid: true},
	}
	return authorParams
}
