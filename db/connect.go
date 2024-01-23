package db

import (
	"database/sql"
	"log"

	"github.com/devldm/go-server-rss/config"
	"github.com/joho/godotenv"
)

func ConnectDB() (*sql.DB, error) {
	godotenv.Load()

	portString := config.Config("PORT")

	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := config.Config("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cant connect to the database", err)
	}

	if err != nil {
		log.Fatal("Cant create db connection", err)
	}

	return conn, err
}
