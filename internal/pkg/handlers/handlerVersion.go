package handlers

import (
	"main/internal/pkg/util"
	"net/http"
	"os"
)

func (apiCfg *ApiConfig) HandlerVersion(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("VERSION")
	type gitguruVersion struct {
		Version string `json:"version"`
	}
	util.RespondWithJSON(w, 200, gitguruVersion{Version: version})
}
