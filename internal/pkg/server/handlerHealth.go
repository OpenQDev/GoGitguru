package server

import (
	"net/http"
)

type HandlerHealthResponse struct{}

func (apiCfg *ApiConfig) HandlerHealth(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, 200, HandlerHealthResponse{})
}
