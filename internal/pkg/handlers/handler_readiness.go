package handlers

import (
	"main/internal/pkg/util"
	"net/http"
)

func (apiCfg *ApiConfig) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	util.RespondWithJSON(w, 200, struct{}{})
}
