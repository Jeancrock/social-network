package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"

	"social_network_backend/internal/sessions"
	"social_network_backend/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// validateEmail vérifie un format email basique et longueur.
func validateEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	// Regex simple, suffisante pour une validation de surface
	re := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	return re.MatchString(email)
}

// validateUsername limite la taille et les caractères autorisés.
func validateUsername(u string) bool {
	if len(u) < 3 || len(u) > 32 {
		return false
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9_\. -]+$`)
	return re.MatchString(u)
}

// validatePassword limite la taille minimale.
func validatePassword(p string) bool {
	return len(p) >= 0 && len(p) <= 200
}

// HandleRegister crée un compte utilisateur avec hash du mot de passe et validations.
func HandleRegister(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if !validateUsername(req.Username) || !validateEmail(req.Email) || !validatePassword(req.Password) {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Vérifie l’unicité username
	var exists int
	if err := db.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", req.Username).Scan(&exists); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		http.Error(w, "Username already taken", http.StatusConflict)
		return
	}

	// Vérifie l’unicité email
	if err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email=?", req.Email).Scan(&exists); err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}
	if exists > 0 {
		http.Error(w, "Email already registered", http.StatusConflict)
		return
	}

	// Hash du mot de passe
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	u := models.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashed),
	}

	if _, err = db.Exec("INSERT INTO users(id, username, email, password) VALUES(?,?,?,?)",
		u.ID, u.Username, u.Email, u.Password); err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": u.ID})
}

// HandleDeleteAccount supprime le compte de l’utilisateur connecté (via session) et nettoie la session.
func HandleDeleteAccount(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, sessionID, err := sessions.GetUserIDFromRequest(db, r)
	if err != nil {
		log.Println("❌ Session invalide:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("🟢 Suppression de l'utilisateur ID =", userID)

	// Supprimer les posts de l'utilisateur
	if _, err = db.Exec("DELETE FROM posts WHERE userId=?", userID); err != nil {
		log.Println("❌ Erreur delete posts:", err)
		http.Error(w, "Failed to delete posts", http.StatusInternalServerError)
		return
	}

	// Supprimer les relations followers (où il est suivi ou suiveur)
	if _, err = db.Exec("DELETE FROM followers WHERE userId=? OR followerId=?", userID, userID); err != nil {
		log.Println("❌ Erreur delete followers:", err)
		http.Error(w, "Failed to delete followers", http.StatusInternalServerError)
		return
	}

	// Supprimer les memberships de groupes
	if _, err = db.Exec("DELETE FROM group_members WHERE userId=?", userID); err != nil {
		log.Println("❌ Erreur delete group_members:", err)
		http.Error(w, "Failed to delete memberships", http.StatusInternalServerError)
		return
	}

	// Supprimer l'utilisateur
	if _, err = db.Exec("DELETE FROM users WHERE id=?", userID); err != nil {
		log.Println("❌ Erreur delete user:", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Supprimer la session DB et le cookie
	if sessionID != "" {
		_ = sessions.DeleteSession(db, sessionID)
	}
	sessions.ClearSessionCookie(w)

	log.Println("✅ Compte supprimé avec succès")
	w.WriteHeader(http.StatusOK)
}
