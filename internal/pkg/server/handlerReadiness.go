package server

import (
	"net/http"
)

type HandlerReadinessResponse struct{}

func (apiCfg *ApiConfig) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, 200, HandlerReadinessResponse{})
}
