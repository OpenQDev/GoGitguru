package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"main/internal/database"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type HandlerGithubUserByLoginRequest struct{}
type HandlerGithubUserByLoginResponse struct{}

func (apiConfig *ApiConfig) HandlerGithubUserByLogin(w http.ResponseWriter, r *http.Request) {
	githubAccessToken := r.Header.Get("GH-Authorization")

	if githubAccessToken == "" {
		RespondWithError(w, http.StatusUnauthorized, "You must provide a GitHub access token.")
		return
	}

	login := chi.URLParam(r, "login")

	initialDbUser, err := apiConfig.DB.GetGithubUser(context.Background(), login)

	if err == nil {
		RespondWithJSON(w, 200, initialDbUser)
		return
	} else {
		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api.github.com/users/"+login, nil)

		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create request: %s", err))
			return
		}

		req.Header.Add("Authorization", "token "+githubAccessToken)
		resp, err := client.Do(req)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to make request: %s", err))
			return
		}

		defer resp.Body.Close()

		var user User
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to decode response.: %s", err))
			return
		}

		layout := "2006-01-02T15:04:05Z" // ISO 8601 format
		createdAt, err := time.Parse(layout, user.CreatedAt)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse CreatedAt: %s", err))
			return
		}
		updatedAt, err := time.Parse(layout, user.UpdatedAt)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to parse UpdatedAt: %s", err))
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
			RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to insert user into database: %s", err))
			return
		}

		RespondWithJSON(w, http.StatusOK, user)
	}
}
