package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"social_network_backend/cors"
	"social_network_backend/router"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	dbPath := os.Getenv("SQLITE_PATH")
	if dbPath == "" {
		dbPath = "./db/social.db"
	}

	var err error
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	defer db.Close()

	createTables()

	mux := http.NewServeMux()
	router.RegisterRoutes(mux, db)

	log.Println("Backend running on :8080")
	log.Fatal(http.ListenAndServe(":8080", cors.EnableCORS(mux)))
}

// ----------------- UTIL -----------------
func createTables() {
	// Création des tables utilisateurs
	db.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			id TEXT PRIMARY KEY,
			username TEXT,
			email TEXT UNIQUE,
			password TEXT
		)
	`)

	// Création de la table followers
	db.Exec(`
		CREATE TABLE IF NOT EXISTS followers(
			id TEXT PRIMARY KEY,
			userId TEXT,       -- l'utilisateur suivi
			followerId TEXT,   -- l'utilisateur qui suit
			FOREIGN KEY(userId) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(followerId) REFERENCES users(id) ON DELETE CASCADE
		)
	`)

	// Création de la table groupes
	db.Exec(`
		CREATE TABLE IF NOT EXISTS groups(
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL	
		)
	`)

	// Création de la table membres du groupe
	db.Exec(`
		CREATE TABLE IF NOT EXISTS group_members(
			id TEXT PRIMARY KEY,
			userId TEXT,
			groupId TEXT,
			FOREIGN KEY(userId) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(groupId) REFERENCES groups(id) ON DELETE CASCADE
		)
	`)

	db.Exec(`
	CREATE TABLE IF NOT EXISTS posts(
		id TEXT PRIMARY KEY,
		userId TEXT,
		content TEXT,
		created TEXT,
		FOREIGN KEY(userId) REFERENCES users(id)
	)
	`)
}
