package main

import (
	"net/http"
	"os"
)

func (apiCfg *apiConfig) handler_version(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("VERSION")
	type gitguruVersion struct {
		Version string `json:"version"`
	}
	respondWithJSON(w, 200, gitguruVersion{Version: version})
}
