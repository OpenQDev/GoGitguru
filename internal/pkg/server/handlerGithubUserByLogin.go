package server

import (
	"context"
	"encoding/json"
	"main/internal/database"
	"net/http"

	"github.com/go-chi/chi"
)

type User struct {
	InternalID      int    `json:"internal_id"`
	GithubRestID    int    `json:"github_rest_id"`
	GithubGraphqlID string `json:"github_graphql_id"`
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

	// Fetch user from database
	user, err := apiConfig.DB.GetUserByLogin(context.Background(), login)
	if err != nil {
		RespondWithError(w, 500, "Failed to fetch user from database.")
		return
	}

	// If user is not found in database, fetch from GitHub API
	if user == nil {
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

		// Insert the user into the database using sqlc generated methods
		params := database.InsertUserParams{
			// Fill in the fields based on your database schema
		}

		_, err = apiConfig.DB.InsertUser(context.Background(), params)
		if err != nil {
			RespondWithError(w, 500, "Failed to insert user into database.")
			return
		}
	}

	RespondWithJSON(w, 200, user)
}
