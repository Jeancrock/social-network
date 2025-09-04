package handlers

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
)

func HandleFollow(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userId := cookie.Value // l'utilisateur connecté
	targetId := r.URL.Query().Get("userId")
	if targetId == "" || targetId == userId {
		http.Error(w, "Invalid user", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPost: // follow
		_, err := db.Exec("INSERT INTO followers(id,userId,followerId) VALUES(?,?,?)",
			uuid.New().String(), targetId, userId)
		if err != nil {
			http.Error(w, "Failed to follow", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	case http.MethodDelete: // unfollow
		_, err := db.Exec("DELETE FROM followers WHERE userId=? AND followerId=?", targetId, userId)
		if err != nil {
			http.Error(w, "Failed to unfollow", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
