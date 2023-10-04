package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"main/internal/database"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type User struct {
	InternalID      int    `json:"internal_id"`
	GithubRestID    int    `json:"id"`
	GithubGraphqlID string `json:"node_id"`
	Login           string `json:"login"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	AvatarURL       string `json:"avatar_url"`
	Company         string `json:"company"`
	Location        string `json:"location"`
	Bio             string `json:"bio"`
	Blog            string `json:"blog"`
	Hireable        bool   `json:"hireable"`
	TwitterUsername string `json:"twitter_username"`
	Followers       int    `json:"followers"`
	Following       int    `json:"following"`
	Type            string `json:"type"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

func (apiConfig *ApiConfig) HandlerGithubUserByLogin(w http.ResponseWriter, r *http.Request) {
	githubAccessToken := r.Header.Get("GH-Authorization")

	if githubAccessToken == "" {
		RespondWithError(w, 400, "You must provide a GitHub access token.")
		return
	}

	login := chi.URLParam(r, "login")

	// Fetch initialDbUser from database
	initialDbUser, err := apiConfig.DB.GetGithubUser(context.Background(), login)

	if err == nil {
		RespondWithJSON(w, 200, initialDbUser)
		return
	} else {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api.github.com/users/"+login, nil)

		if err != nil {
			RespondWithError(w, 500, "Failed to create request.")
			return
		}

		req.Header.Add("Authorization", "token "+githubAccessToken)
		resp, err := client.Do(req)
		if err != nil {
			RespondWithError(w, 500, "Failed to make request.")
			return
		}

		defer resp.Body.Close()

		var user User
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			RespondWithError(w, 500, "Failed to decode response.")
			return
		}

		layout := "2006-01-02T15:04:05Z" // ISO 8601 format
		createdAt, err := time.Parse(layout, user.CreatedAt)
		if err != nil {
			RespondWithError(w, 500, "Failed to parse CreatedAt.")
			return
		}
		updatedAt, err := time.Parse(layout, user.UpdatedAt)
		if err != nil {
			RespondWithError(w, 500, "Failed to parse UpdatedAt.")
			return
		}

		// Insert the user into the database using sqlc generated methods
		params := database.InsertUserParams{
			GithubRestID:    int32(user.GithubRestID),
			GithubGraphqlID: user.GithubGraphqlID,
			Login:           user.Login,
			Name:            sql.NullString{String: user.Name, Valid: user.Name != ""},
			Email:           sql.NullString{String: user.Email, Valid: user.Email != ""},
			AvatarUrl:       sql.NullString{String: user.AvatarURL, Valid: user.AvatarURL != ""},
			Company:         sql.NullString{String: user.Company, Valid: user.Company != ""},
			Location:        sql.NullString{String: user.Location, Valid: user.Location != ""},
			Bio:             sql.NullString{String: user.Bio, Valid: user.Bio != ""},
			Blog:            sql.NullString{String: user.Blog, Valid: user.Blog != ""},
			Hireable:        sql.NullBool{Bool: user.Hireable, Valid: true},
			TwitterUsername: sql.NullString{String: user.TwitterUsername, Valid: user.TwitterUsername != ""},
			Followers:       sql.NullInt32{Int32: int32(user.Followers), Valid: true},
			Following:       sql.NullInt32{Int32: int32(user.Following), Valid: true},
			Type:            user.Type,
			CreatedAt:       sql.NullTime{Time: createdAt, Valid: true},
			UpdatedAt:       sql.NullTime{Time: updatedAt, Valid: true},
		}

		_, err = apiConfig.DB.InsertUser(context.Background(), params)
		if err != nil {
			RespondWithError(w, 500, "Failed to insert user into database.")
			return
		}

		RespondWithJSON(w, 200, user)
	}
}
