package handlers

import (
	"database/sql"
	"net/http"

	"social_network_backend/internal/sessions"

	"github.com/google/uuid"
)

// HandleFollow permet de suivre / ne plus suivre un utilisateur. Nécessite une session valide.
func HandleFollow(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userID, _, err := sessions.GetUserIDFromRequest(db, r)
	if err != nil || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	targetID := r.URL.Query().Get("userId")
	if targetID == "" || targetID == userID {
		http.Error(w, "Invalid user", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost: // follow
		if _, err := db.Exec("INSERT INTO followers(id,userId,followerId) VALUES(?,?,?)",
			uuid.New().String(), targetID, userID); err != nil {
			http.Error(w, "Failed to follow", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete: // unfollow
		if _, err := db.Exec("DELETE FROM followers WHERE userId=? AND followerId=?", targetID, userID); err != nil {
			http.Error(w, "Failed to unfollow", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
