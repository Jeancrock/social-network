package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"social_network_backend/internal/sessions"

	"golang.org/x/crypto/bcrypt"
)

// HandleLogin vérifie les identifiants, crée une session (DB) et place un cookie sécurisé.
func HandleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	creds.Email = strings.TrimSpace(strings.ToLower(creds.Email))
	if !validateEmail(creds.Email) || !validatePassword(creds.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	row := db.QueryRow("SELECT id, password FROM users WHERE email=?", creds.Email)
	var userID, hashed string
	if err := row.Scan(&userID, &hashed); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(creds.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Créer une session DB (expire dans 24h)
	exp := time.Now().Add(24 * time.Hour)
	sessionID, err := sessions.CreateSession(db, userID, exp)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}
	sessions.SetSessionCookie(w, sessionID, exp)

	w.WriteHeader(http.StatusOK)
}

// HandleLogout supprime la session active et vide le cookie.
func HandleLogout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	sessionID := sessions.GetSessionIDFromCookie(r)
	if sessionID != "" {
		_ = sessions.DeleteSession(db, sessionID)
	}
	sessions.ClearSessionCookie(w)
	w.WriteHeader(http.StatusOK)
}

// ensureAuthenticated renvoie l'ID utilisateur à partir de la session, sinon erreur.
func ensureAuthenticated(db *sql.DB, r *http.Request) (string, error) {
	uid, _, err := sessions.GetUserIDFromRequest(db, r)
	if err != nil || uid == "" {
		return "", errors.New("unauthorized")
	}
	return uid, nil
}
