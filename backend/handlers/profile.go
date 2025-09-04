package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social_network_backend/models"
)

// Handler de profil
func HandleProfile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Lire ID en query ?id= ou sinon utiliser cookie.Value
	id := r.URL.Query().Get("id")
	if id == "" {
		id = cookie.Value
	}

	// Récup user
	var u models.User
	row := db.QueryRow("SELECT id, username, email FROM users WHERE id=?", id)
	if err := row.Scan(&u.ID, &u.Username, &u.Email); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Récup posts
	rows, _ := db.Query(
		"SELECT id, userId, content, created FROM posts WHERE userId=? ORDER BY created DESC",
		id,
	)
	defer rows.Close()
	var posts []models.Post
	for rows.Next() {
		var p models.Post
		rows.Scan(&p.ID, &p.UserID, &p.Content, &p.Created)
		posts = append(posts, p)
	}

	// Assurer tableau vide si null
	if posts == nil {
		posts = []models.Post{}
	}
	// Followers
	rows, _ = db.Query(`
		SELECT u.id, u.username, u.email
		FROM followers f JOIN users u ON u.id=f.followerId
		WHERE f.userId=?`, id)
	defer rows.Close()
	var followers []models.User
	for rows.Next() {
		var f models.User
		rows.Scan(&f.ID, &f.Username, &f.Email)
		followers = append(followers, f)
	}

	// Assurer tableau vide si null
	if followers == nil {
		followers = []models.User{}
	}
	// Following
	rows, _ = db.Query(`
		SELECT u.id, u.username, u.email
		FROM followers f JOIN users u ON u.id=f.userId
		WHERE f.followerId=?`, id)
	defer rows.Close()
	var following []models.User
	for rows.Next() {
		var f models.User
		rows.Scan(&f.ID, &f.Username, &f.Email)
		following = append(following, f)
	}

	// Assurer tableau vide si null
	if following == nil {
		following = []models.User{}
	}

	// Groups
	rows, _ = db.Query(`
		SELECT g.name
		FROM group_members gm JOIN groups g ON g.id=gm.groupId
		WHERE gm.userId=?`, id)
	defer rows.Close()
	var groups []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		groups = append(groups, name)
	}
	// Assurer tableau vide si null

	if groups == nil {
		groups = []string{}
	}

	resp := models.ProfileResponse{
		User:      u,
		Posts:     posts,
		Followers: followers,
		Following: following,
		Groups:    groups,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
