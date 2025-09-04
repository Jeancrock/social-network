package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"social_network_backend/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// création de compte
func HandleRegister(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Décodage du JSON envoyé par le frontend
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Vérification des champs
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Missing fields", http.StatusBadRequest)
		return
	}

	// Hash du mot de passe
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Création de l'utilisateur à stocker
	u := models.User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashed),
	}

	// Insertion dans la DB
	_, err = db.Exec("INSERT INTO users(id, username, email, password) VALUES(?,?,?,?)",
		u.ID, u.Username, u.Email, u.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error inserting user: %v", err), http.StatusInternalServerError)
		return
	}

	// Réponse (on ne renvoie JAMAIS le mot de passe)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": u.ID})
}

// Suppression de compte
func HandleDeleteAccount(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("❌ Pas de cookie:", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	log.Println("🟢 Suppression de l'utilisateur ID =", cookie.Value)

	// Supprimer les posts
	_, err = db.Exec("DELETE FROM posts WHERE userId=?", cookie.Value)
	if err != nil {
		log.Println("❌ Erreur delete posts:", err)
		http.Error(w, "Failed to delete posts", http.StatusInternalServerError)
		return
	}

	// Supprimer l'utilisateur
	_, err = db.Exec("DELETE FROM users WHERE id=?", cookie.Value)
	if err != nil {
		log.Println("❌ Erreur delete user:", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Supprimer le cookie de session
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	log.Println("✅ Compte supprimé avec succès")
	w.WriteHeader(http.StatusOK)
}
