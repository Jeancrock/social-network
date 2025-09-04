package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social_network_backend/models"
)

// Récupérer tous les utilisateurs (sans les mots de passe)
func HandleUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT id, username, email FROM users")
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		rows.Scan(&u.ID, &u.Username, &u.Email)
		users = append(users, u)
	}
	json.NewEncoder(w).Encode(users)
}
