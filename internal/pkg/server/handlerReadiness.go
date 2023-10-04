package server

import (
	"net/http"
)

func (apiCfg *ApiConfig) HandlerReadiness(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, 200, struct{}{})
}
