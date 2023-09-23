package main

import (
	"encoding/json"
	"fmt"
	"main/internal/pkg/logger"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string) {

	type errResponse struct {
		Error string `json:"error"`
	}

	respondWithJSON(w, code, errResponse{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)

	if err != nil {
		logger.LogError("failed to marshall JSON response: %v", payload)
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("failed to marshall JSON response: %v", payload)))
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
