package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

// escapeLike échappe les métacaractères % et _ pour un LIKE sûr (selon collations).
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	return s
}

// HandleSearch recherche des utilisateurs et des groupes par mot-clé, avec échappement LIKE et LIMIT.
func HandleSearch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	query := strings.TrimSpace(r.URL.Query().Get("query"))
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}
	// Limite côté appli pour éviter des requêtes trop lourdes
	const limit = 50
	like := "%" + escapeLike(query) + "%"

	// --- Users ---
	userRows, err := db.Query(`SELECT id, username FROM users WHERE username LIKE ? ESCAPE '\' LIMIT ?`, like, limit)
	if err != nil {
		http.Error(w, "Failed to search users", http.StatusInternalServerError)
		return
	}
	defer userRows.Close()

	var users []map[string]string
	for userRows.Next() {
		var id, username string
		if err := userRows.Scan(&id, &username); err != nil {
			http.Error(w, "Error scanning user", http.StatusInternalServerError)
			return
		}
		users = append(users, map[string]string{"id": id, "username": username})
	}

	// --- Groups ---
	groupRows, err := db.Query(`SELECT id, name FROM groups WHERE name LIKE ? ESCAPE '\' LIMIT ?`, like, limit)
	if err != nil {
		http.Error(w, "Failed to search groups", http.StatusInternalServerError)
		return
	}
	defer groupRows.Close()

	var groups []map[string]string
	for groupRows.Next() {
		var id, name string
		if err := groupRows.Scan(&id, &name); err != nil {
			http.Error(w, "Error scanning group", http.StatusInternalServerError)
			return
		}
		groups = append(groups, map[string]string{"id": id, "name": name})
	}

	if users == nil {
		users = []map[string]string{}
	}
	if groups == nil {
		groups = []map[string]string{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"users":  users,
		"groups": groups,
	})
}
