package handlers

import (
	"net/http"
	"os"
)

func (apiCfg *ApiConfig) HandlerVersion(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("VERSION")
	type gitguruVersion struct {
		Version string `json:"version"`
	}
	RespondWithJSON(w, 200, gitguruVersion{Version: version})
}
