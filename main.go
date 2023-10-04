package main

import (
	"main/internal/database"
	"main/internal/pkg/handlers"
	"main/internal/pkg/logger"
	"main/internal/pkg/setup"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

func main() {
	portString, dbUrl, originUrl, debugMode := setup.ExtractAndVerifyEnvironment(".env")

	logger.SetDebugMode(debugMode)

	database, apiCfg := setup.PrepareDatabase(dbUrl)

	beginSyncingInBackground(database)

	startServer(apiCfg, portString, originUrl)
}

func startServer(apiCfg handlers.ApiConfig, portString string, originUrl string) {
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
	v1Router.Get("/healthz", apiCfg.HandlerReadiness)
	v1Router.Get("/version", apiCfg.HandlerVersion)
	v1Router.Post("/add", apiCfg.HandlerAdd)

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

// Initialize periodic syncing in the background
func beginSyncingInBackground(database *database.Queries) {
	logger.LogBlue("Beginning sync for all repo urls...")
	go startSyncing(database, "repos", 10, 10*time.Second)
}
