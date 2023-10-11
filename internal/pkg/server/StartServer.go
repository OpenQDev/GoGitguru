package server

import (
	"main/internal/pkg/logger"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func StartServer(apiCfg ApiConfig, portString string, originUrl string) {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{originUrl},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", apiCfg.HandlerHealth)
	v1Router.Get("/version", apiCfg.HandlerVersion)

	v1Router.Post("/add", apiCfg.HandlerAdd)

	v1Router.Get("/repos/github/{owner}/{name}", apiCfg.HandlerGithubRepoByOwnerAndName)
	v1Router.Get("/repos/github/{owner}", apiCfg.HandlerGithubReposByOwner)

	v1Router.Get("/users/github/{login}", apiCfg.HandlerGithubUserByLogin)
	v1Router.Post("/users/github/{login}/commits", apiCfg.HandlerGithubUserCommits)
	v1Router.Post("/repos/commits", apiCfg.HandlerRepoCommits)

	v1Router.Post("/dependency-history", apiCfg.HandlerDependencyHistory)

	router.Mount("/", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	logger.LogBlue("server starting on port %v", portString)
	srverr := srv.ListenAndServe()

	if srverr != nil {
		logger.LogFatalRedAndExit("the gitguru server encountered an error: %s", srverr)
	}
}
