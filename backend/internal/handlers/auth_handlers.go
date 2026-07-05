package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AntonioMartino67/showio/backend/internal/auth"
	"github.com/AntonioMartino67/showio/backend/internal/database"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo della richiesta non valido", http.StatusBadRequest)
		return
	}

	// Validazione minima
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Username, email e password sono obbligatori", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		http.Error(w, "La password deve avere almeno 8 caratteri", http.StatusBadRequest)
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Errore interno durante la registrazione", http.StatusInternalServerError)
		return
	}

	var newID string
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err = database.Pool.QueryRow(r.Context(), query, req.Username, req.Email, hash).Scan(&newID)
	if err != nil {
		http.Error(w, "Username o email già in uso", http.StatusConflict)
		return
	}

	resp := RegisterResponse{
		ID:       newID,
		Username: req.Username,
		Email:    req.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo della richiesta non valido", http.StatusBadRequest)
		return
	}

	var userID, passwordHash string
	query := `SELECT id, password_hash FROM users WHERE email = $1 AND deleted_at IS NULL`
	err := database.Pool.QueryRow(r.Context(), query, req.Email).Scan(&userID, &passwordHash)
	if err != nil {
		http.Error(w, "Credenziali non valide", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPassword(req.Password, passwordHash) {
		http.Error(w, "Credenziali non valide", http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Errore interno durante il login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	var username, email string
	query := `SELECT username, email FROM users WHERE id = $1`
	err := database.Pool.QueryRow(r.Context(), query, userID).Scan(&username, &email)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":       userID,
		"username": username,
		"email":    email,
	})
}