package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/devldm/go-server-rss/config"
	"github.com/devldm/go-server-rss/db"
	"github.com/devldm/go-server-rss/internal/database"
	"github.com/devldm/go-server-rss/router"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	portString := config.Config("PORT")

	conn, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	nDatabase := database.New(conn)
	// apiCfg := apiConfig{
	// 	DB: db,
	// }

	go startScraping(
		nDatabase,
		10,
		time.Minute,
	)

	router := router.SetupRouter(nDatabase)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("server starting on port %v", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port:", portString)
}
