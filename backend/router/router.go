package router

import (
	"database/sql"
	"net/http"

	"social_network_backend/handlers"
)

// Enregistre toutes les routes de l’API
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleRegister(w, r, db)
	})
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLogin(w, r, db)
	})
	mux.HandleFunc("/api/logout", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleLogout(w, r, db)
	})
	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlePosts(w, r, db)
	})
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleUsers(w, r, db)
	})
	mux.HandleFunc("/api/profile", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleProfile(w, r, db)
	})
	mux.HandleFunc("/api/delete-account", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleDeleteAccount(w, r, db)
	})
	mux.HandleFunc("/api/my-posts", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleMyPosts(w, r, db)
	})
	mux.HandleFunc("/api/follow", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleFollow(w, r, db)
	})
}

// ⚠️ Ici tu reprends toutes tes fonctions : handleRegister, handleLogin, handlePosts…
// Je ne les recopie pas en entier car elles sont déjà dans ton code,
// il faudra juste les adapter en rajoutant `db *sql.DB` en paramètre.
