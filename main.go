package main

import (
	"database/sql"
	"main/internal/database"
	"main/internal/pkg/handlers"
	"main/internal/pkg/logger"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	portString, dbUrl, originUrl, debugMode := extractAndVerifyEnvironment()

	logger.SetDebugMode(debugMode)

	database, apiCfg := prepareDatabase(dbUrl)

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

func prepareDatabase(dbUrl string) (*database.Queries, handlers.ApiConfig) {
	database, err := getDatbase(dbUrl)
	if err != nil {
		logger.LogError("error getting database: %s", err)
	}

	apiCfg, err := getApiConfig(database)
	if err != nil {
		logger.LogFatalRedAndExit("can't connect to DB: %s", err)
	}
	return database, apiCfg
}

func extractAndVerifyEnvironment() (string, string, string, bool) {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")
	originUrl := os.Getenv("ORIGIN_URL")
	debugMode := os.Getenv("DEBUG_MODE")

	if portString == "" || dbUrl == "" || originUrl == "" || debugMode == "" {
		logger.LogFatalRedAndExit("PORT | DB_URL | ORIGIN_URL | DEBUG_MODE is not found in the environment")
	}

	debug, err := strconv.ParseBool(debugMode)
	if err != nil {
		logger.LogFatalRedAndExit("DEBUG_MODE must be a boolean")
	}

	return portString, dbUrl, originUrl, debug
}

func getDatbase(dbUrl string) (*database.Queries, error) {
	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return nil, err
	}

	queries := database.New(conn)

	return queries, nil
}

func getApiConfig(database *database.Queries) (handlers.ApiConfig, error) {
	apiCfg := handlers.ApiConfig{
		DB: database,
	}

	return apiCfg, nil
}
