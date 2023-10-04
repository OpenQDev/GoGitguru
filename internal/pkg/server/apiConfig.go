package server

import "main/internal/database"

type ApiConfig struct {
	DB *database.Queries
}
