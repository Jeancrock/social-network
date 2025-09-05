package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

// RunMigrations applique toutes les migrations SQLite
func RunMigrations(dbPath string) {
	// Les fichiers doivent être nommés comme suit: *.up.sql et *.down.sql
	migrationsPath, err := filepath.Abs("db/migrations")
	if err != nil {
		log.Fatalf("Erreur résolution chemin migrations: %v", err)
	}

	absDBPath, err := filepath.Abs(dbPath)
	if err != nil {
		log.Fatalf("Erreur résolution chemin DB: %v", err)
	}

	// Création du fichier DB s’il n’existe pas
	if _, err := os.Stat(absDBPath); os.IsNotExist(err) {
		file, err := os.Create(absDBPath)
		if err != nil {
			log.Fatalf("Impossible de créer la DB: %v", err)
		}
		file.Close()
	}

	dbConn, err := sql.Open("sqlite3", absDBPath)
	if err != nil {
		log.Fatalf("Impossible d'ouvrir DB: %v", err)
	}

	driver, err := sqlite3.WithInstance(dbConn, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Impossible de créer driver migrate: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatalf("Erreur création migrate: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Erreur migrations: %v", err)
	}
}

// InitDB ouvre la DB et applique les migrations
func InitDB(dbPath string) *sql.DB {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		f, err := os.Create(dbPath)
		if err != nil {
			log.Fatalf("Impossible de créer DB: %v", err)
		}
		f.Close()
	}

	RunMigrations(dbPath)

	dbConn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	if err := dbConn.Ping(); err != nil {
		log.Fatal(err)
	}

	return dbConn
}
