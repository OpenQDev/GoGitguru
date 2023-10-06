package server

import "net/http"

func (apiCfg *ApiConfig) HandlerDependencyHistory(w http.ResponseWriter, r *http.Request) {
	RespondWithJSON(w, 202, struct{}{})
}
