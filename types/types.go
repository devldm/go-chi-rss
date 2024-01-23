package types

import "github.com/devldm/go-server-rss/internal/database"

type ApiConfig struct {
	DB *database.Queries
}
