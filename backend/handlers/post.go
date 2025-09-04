package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"social_network_backend/models"

	"github.com/google/uuid"
)

// Handler de posts
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
			rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Content, &p.Created)
			posts = append(posts, p)
		}
		json.NewEncoder(w).Encode(posts)

	case http.MethodPost:
		var p models.Post
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		cookie, err := r.Cookie("session_id")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		p.UserID = cookie.Value
		p.ID = uuid.New().String()
		p.Created = time.Now().Format(time.RFC3339)

		_, err = db.Exec("INSERT INTO posts(id,userId,content,created) VALUES(?,?,?,?)",
			p.ID, p.UserID, p.Content, p.Created)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting post: %v", err), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(p)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Récupérer les posts de l'utilisateur connecté
func HandleMyPosts(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query(`
        SELECT p.id, p.userId, u.username, p.content, p.created 
        FROM posts p 
        JOIN users u ON p.userId = u.id 
        WHERE p.userId = ?
        ORDER BY p.created DESC
    `, cookie.Value)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Content, &p.Created)
		posts = append(posts, p)
	}
	json.NewEncoder(w).Encode(posts)
}
