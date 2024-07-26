package server

import (
	"net/http"

	"github.com/OpenQDev/GoGitguru/util/logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func StartServer(apiCfg ApiConfig, portString string, originUrl string) {
	// Initialize a primary Chi router
	// This is where global middleware will be attached
	router := chi.NewRouter()

	// Configure CORS
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{originUrl},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Initialize an application router for your actual routes
	v1Router := chi.NewRouter()

	// UTIL
	v1Router.Get("/healthz", apiCfg.HandlerHealth)
	v1Router.Get("/version", apiCfg.HandlerVersion)

	// REPOSITORY
	v1Router.Post("/add", apiCfg.HandlerAdd)
	v1Router.Post("/add-dependency-pattern", apiCfg.HandlerAddDependencyPattern)
	v1Router.Get("/repos/github/{owner}/{name}", apiCfg.HandlerGithubRepoByOwnerAndName)
	v1Router.Get("/repos/github/{owner}", apiCfg.HandlerGithubReposByOwner)
	v1Router.Post("/repos/commits", apiCfg.HandlerRepoCommits)
	v1Router.Post("/repos/authors", apiCfg.HandlerRepoAuthors)

	// USER
	v1Router.Get("/users/github/{login}", apiCfg.HandlerGithubUserByLogin)
	v1Router.Post("/users/github/{login}/commits", apiCfg.HandlerGithubUserCommits)

	// DEPENDENCY HISTORY
	v1Router.Post("/dependency-history", apiCfg.HandlerDependencyHistory)

	// DEPENDENCY HISTORY
	v1Router.Post("/status", apiCfg.HandlerStatus)

	// MISCELLANEOUS
	v1Router.Get("/get-next-repo-url", apiCfg.HandlerGetNextRepoUrl)
	v1Router.Post("/first-commit", apiCfg.HandlerFirstCommit)

	// Mounting sub-router "v1Router" to the primary router "router" so CORS applies to all routes
	router.Mount("/", v1Router)

	// Turn the Chi router into a an http.Server pointer
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	// Start server on port "portSting"
	logger.LogBlue("server starting on port %v", portString)
	srverr := srv.ListenAndServe()

	if srverr != nil {
		logger.LogFatalRedAndExit("the gitguru server encountered an error: %s", srverr)
	}
}
