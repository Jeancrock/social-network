package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"social_network_backend/internal/sessions"
	"social_network_backend/models"

	"github.com/google/uuid"
)

// validatePostContent contrôle la taille du contenu du post.
func validatePostContent(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) > 0 && len(s) <= 2000
}

// HandlePosts liste/ajoute des posts. Ajout nécessite une session valide.
func HandlePosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query(`
			SELECT p.id, p.userId, u.username, p.content, p.created 
			FROM posts p JOIN users u ON p.userId = u.id 
			ORDER BY p.created DESC
		`)
		if err != nil {
			http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var posts []models.Post
		for rows.Next() {
			var p models.Post
			if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Content, &p.Created); err != nil {
				http.Error(w, "Scan error", http.StatusInternalServerError)
				return
			}
			posts = append(posts, p)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(posts)

	case http.MethodPost:
		userID, _, err := sessions.GetUserIDFromRequest(db, r)
		if err != nil || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var p models.Post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		if !validatePostContent(p.Content) {
			http.Error(w, "Invalid content", http.StatusBadRequest)
			return
		}

		p.UserID = userID
		p.ID = uuid.New().String()
		p.Created = time.Now().Format(time.RFC3339)

		if _, err = db.Exec("INSERT INTO posts(id,userId,content,created) VALUES(?,?,?,?)",
			p.ID, p.UserID, p.Content, p.Created); err != nil {
			http.Error(w, fmt.Sprintf("Error inserting post: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(p)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleMyPosts récupère les posts de l'utilisateur connecté via session.
func HandleMyPosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, _, err := sessions.GetUserIDFromRequest(db, r)
	if err != nil || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query(`
        SELECT p.id, p.userId, u.username, p.content, p.created 
        FROM posts p 
        JOIN users u ON p.userId = u.id 
        WHERE p.userId = ?
        ORDER BY p.created DESC
    `, userID)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		if err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Content, &p.Created); err != nil {
			http.Error(w, "Scan error", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(posts)
}
