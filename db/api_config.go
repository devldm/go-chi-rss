package db

import "github.com/devldm/go-server-rss/internal/database"

type APIConfig struct {
	DB *database.Queries
}

// Initialize the config
func NewAPIConfig(db *database.Queries) *APIConfig {
	return &APIConfig{
		DB: db,
	}
}
