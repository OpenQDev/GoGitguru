package server

import (
	"net/http"
	"os"
)

type HandlerVersionResponse struct {
	Version string `json:"version"`
}

func (apiCfg *ApiConfig) HandlerVersion(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("VERSION")
	RespondWithJSON(w, 200, HandlerVersionResponse{Version: version})
}
