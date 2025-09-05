package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"social_network_backend/internal/sessions"
	"social_network_backend/models"
)

// HandleProfile retourne le profil d'un utilisateur (par défaut l'utilisateur courant via session).
func HandleProfile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	currUserID, _, err := sessions.GetUserIDFromRequest(db, r)
	if err != nil || currUserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		id = currUserID
	}

	// Récup utilisateur
	var u models.User
	row := db.QueryRow("SELECT id, username, email FROM users WHERE id=?", id)
	if err := row.Scan(&u.ID, &u.Username, &u.Email); err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Récup posts
	postsRows, err := db.Query(
		"SELECT id, userId, content, created FROM posts WHERE userId=? ORDER BY created DESC",
		id,
	)
	if err != nil {
		http.Error(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	var posts []models.Post
	for postsRows.Next() {
		var p models.Post
		if err := postsRows.Scan(&p.ID, &p.UserID, &p.Content, &p.Created); err != nil {
			postsRows.Close()
			http.Error(w, "Scan error posts", http.StatusInternalServerError)
			return
		}
		posts = append(posts, p)
	}
	postsRows.Close()
	if posts == nil {
		posts = []models.Post{}
	}

	// Followers (qui me suivent)
	followersRows, err := db.Query(`
		SELECT u.id, u.username, u.email
		FROM followers f JOIN users u ON u.id=f.followerId
		WHERE f.userId=?`, id)
	if err != nil {
		http.Error(w, "Failed to fetch followers", http.StatusInternalServerError)
		return
	}
	var followers []models.User
	for followersRows.Next() {
		var f models.User
		if err := followersRows.Scan(&f.ID, &f.Username, &f.Email); err != nil {
			followersRows.Close()
			http.Error(w, "Scan error followers", http.StatusInternalServerError)
			return
		}
		followers = append(followers, f)
	}
	followersRows.Close()
	if followers == nil {
		followers = []models.User{}
	}

	// Following (que je suis)
	followingRows, err := db.Query(`
		SELECT u.id, u.username, u.email
		FROM followers f JOIN users u ON u.id=f.userId
		WHERE f.followerId=?`, id)
	if err != nil {
		http.Error(w, "Failed to fetch following", http.StatusInternalServerError)
		return
	}
	var following []models.User
	for followingRows.Next() {
		var f models.User
		if err := followingRows.Scan(&f.ID, &f.Username, &f.Email); err != nil {
			followingRows.Close()
			http.Error(w, "Scan error following", http.StatusInternalServerError)
			return
		}
		following = append(following, f)
	}
	followingRows.Close()
	if following == nil {
		following = []models.User{}
	}

	// Groupes
	groupsRows, err := db.Query(`
		SELECT g.name
		FROM group_members gm JOIN groups g ON g.id=gm.groupId
		WHERE gm.userId=?`, id)
	if err != nil {
		http.Error(w, "Failed to fetch groups", http.StatusInternalServerError)
		return
	}
	var groups []string
	for groupsRows.Next() {
		var name string
		if err := groupsRows.Scan(&name); err != nil {
			groupsRows.Close()
			http.Error(w, "Scan error groups", http.StatusInternalServerError)
			return
		}
		groups = append(groups, name)
	}
	groupsRows.Close()
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
	_ = json.NewEncoder(w).Encode(resp)
}
