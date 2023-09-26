package main

import (
	"database/sql"
	"main/internal/database"
	"main/internal/pkg/handlers"
	"main/internal/pkg/logger"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

	database, err := getDatbase(dbUrl)
	if err != nil {
		logger.LogError("error getting database: %s", err)
	}

	apiCfg, err := getApiConfig(database)
	if err != nil {
		logger.LogFatalRedAndExit("can't connect to DB: %s", err)
	}

	uploader, err := getS3Uploader()
	if err != nil {
		logger.LogFatalRedAndExit("error initializing AWS session: %s", err)
	}

	// Initialize periodic syncing in the background
	logger.LogBlue("Beginning sync for all repo urls...")
	go startSyncing(uploader, database, "repos", 10, 10*time.Second)

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

func getS3Uploader() (*s3manager.Uploader, error) {
	// Get AWS API key and secret from environment variables
	awsAccessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Create a session using SharedConfigEnable
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(awsAccessKeyID, awsSecretAccessKey, ""),
	},
	)
	if err != nil {
		return nil, err
	}

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	return uploader, nil
}
