package main

import (
	"database/sql"
	"main/internal/database"
	"main/internal/pkg/handlers"
	"main/internal/pkg/logger"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	if portString == "" || dbUrl == "" {
		logger.LogFatalRedAndExit("PORT | DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		logger.LogFatalRedAndExit("can't connect to DB: %s", err)
	}

	queries := database.New(conn)

	apiCfg := handlers.ApiConfig{
		DB: queries,
	}

	// Initialize periodic syncing in the background
	logger.LogBlue("Beginning sync for all repo urls...")
	go startSyncing()

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://example.com"},
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
