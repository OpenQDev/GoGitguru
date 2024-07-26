package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/OpenQDev/GoGitguru/database"
	"github.com/go-chi/chi"
)

type HandlerUserDependenciesResponse []database.GetAllUserDependenciesByUserRow

func (apiCfg *ApiConfig) HandlerUserDependencies(w http.ResponseWriter, r *http.Request) {

	login := strings.ToLower(chi.URLParam(r, "login"))
	if login == "" {
		RespondWithError(w, http.StatusBadRequest, "login parameter is required")
		return
	}

	userDeps, err := apiCfg.DB.GetAllUserDependenciesByUser(r.Context(), login)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error in GetAllUserDependenciesByUser: %s", err))
		return
	}

	RespondWithJSON(w, 202, userDeps)
}
