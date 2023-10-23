package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/OpenQDev/GoGitguru/util/marshaller"

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

	userExists, err := apiConfig.DB.CheckGithubUserExists(context.Background(), login)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if userExists {
		user, err := apiConfig.DB.GetGithubUser(context.Background(), login)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		RespondWithJSON(w, http.StatusOK, ConvertDatabaseInsertUserParamsToServerUser(user))
		return
	}

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
	err = marshaller.ReaderToType(resp.Body, &user)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to decode response.: %s", err))
		return
	}

	params := ConvertServerUserToInsertUserParams(user)

	_, err = apiConfig.DB.InsertUser(context.Background(), params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to insert user into database: %s", err))
		return
	}

	RespondWithJSON(w, http.StatusOK, user)
}
