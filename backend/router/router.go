package router

import (
	"database/sql"
	"net/http"

	"social_network_backend/handlers"
)

// RegisterRoutes enregistre toutes les routes de l’API avec la DB injectée.
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	// Auth & compte
	mux.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRegister(w, r, db)
	})
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLogin(w, r, db)
	})
	mux.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLogout(w, r, db)
	})
	mux.HandleFunc("/api/delete-account", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleDeleteAccount(w, r, db)
	})

	// Utilisateurs
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleUsers(w, r, db)
	})
	mux.HandleFunc("/api/profile", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleProfile(w, r, db)
	})
	mux.HandleFunc("/api/follow", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleFollow(w, r, db)
	})
	mux.HandleFunc("/api/search", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleSearch(w, r, db)
	})

	// Posts
	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlePosts(w, r, db)
	})
	mux.HandleFunc("/api/my-posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMyPosts(w, r, db)
	})
}
