package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"social_network_backend/cors"
	"social_network_backend/db"
	"social_network_backend/router"
)

// main démarre le serveur HTTP et initialise la base (migrations incluses).
func main() {
	dbPath := os.Getenv("SQLITE_PATH")
	if dbPath == "" {
		dbPath = "./db/social.db"
	}

	// Init DB + migrations
	database := db.InitDB(dbPath)
	defer database.Close()

	mux := http.NewServeMux()
	router.RegisterRoutes(mux, database)

	// Serveur avec timeouts pour éviter les connexions pendantes
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           cors.EnableCORS(mux),
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("Backend running on :8080")
	log.Fatal(srv.ListenAndServe())
}
