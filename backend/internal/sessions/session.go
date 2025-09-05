package sessions

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

const cookieName = "session_id"

// CreateSession insère une session (id, user_id, expires_at) et renvoie l'ID de session.
func CreateSession(db *sql.DB, userID string, expires time.Time) (string, error) {
	id := uuid.New().String()
	_, err := db.Exec(
		"INSERT INTO sessions(id, userId, expiresAt) VALUES(?,?,?)",
		id,
		userID,
		expires.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return "", err
	}
	return id, nil
}

// DeleteSession supprime une session par ID.
func DeleteSession(db *sql.DB, sessionID string) error {
	_, err := db.Exec("DELETE FROM sessions WHERE id=?", sessionID)
	return err
}

// GetSessionIDFromCookie récupère l'ID de session dans le cookie.
func GetSessionIDFromCookie(r *http.Request) string {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}
	return c.Value
}

// SetSessionCookie écrit le cookie HTTPOnly (SameSite Lax). (Secure=false pour dev local)
func SetSessionCookie(w http.ResponseWriter, sessionID string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    sessionID,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// Secure: true, // à activer derrière HTTPS
	})
}

// ClearSessionCookie supprime le cookie côté client.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// GetUserIDFromRequest vérifie la session (existence + non expiration) et renvoie (userID, sessionID).
func GetUserIDFromRequest(db *sql.DB, r *http.Request) (string, string, error) {
	sid := GetSessionIDFromCookie(r)
	if sid == "" {
		return "", "", errors.New("no session")
	}
	var userID string
	var expiresStr string
	row := db.QueryRow("SELECT userId, expiresAt FROM sessions WHERE id=?", sid)
	if err := row.Scan(&userID, &expiresStr); err != nil {
		return "", "", errors.New("invalid session")
	}
	exp, err := time.Parse(time.RFC3339, expiresStr)
	if err != nil {
		return "", "", errors.New("invalid expiry")
	}
	if time.Now().After(exp) {
		// Session expirée → clean
		_ = DeleteSession(db, sid)
		return "", "", errors.New("session expired")
	}
	return userID, sid, nil
}
