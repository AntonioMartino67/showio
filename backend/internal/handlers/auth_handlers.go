package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"database/sql"

	"github.com/AntonioMartino67/showio/backend/internal/auth"
	"github.com/AntonioMartino67/showio/backend/internal/database"
	"github.com/AntonioMartino67/showio/backend/internal/email"
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

func generateOTP() string {
	digits := make([]byte, 6)
	rand.Read(digits)
	code := ""
	for _, d := range digits {
		code += fmt.Sprintf("%d", int(d)%10)
	}
	return code
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo della richiesta non valido", http.StatusBadRequest)
		return
	}

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

	otp := generateOTP()
	expiresAt := time.Now().Add(10 * time.Minute)

	var newID string
	query := `
		INSERT INTO users (username, email, password_hash, email_verified, otp_code, otp_expires_at)
		VALUES ($1, $2, $3, false, $4, $5)
		RETURNING id
	`
	err = database.Pool.QueryRow(r.Context(), query, req.Username, req.Email, hash, otp, expiresAt).Scan(&newID)
	if err != nil {
		http.Error(w, "Username o email già in uso", http.StatusConflict)
		return
	}

	if err := email.SendOTPEmail(req.Email, otp); err != nil {
		log.Println("Errore invio email:", err)
		http.Error(w, "Utente creato ma invio email fallito, riprova con /resend-otp", http.StatusInternalServerError)
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

type VerifyOTPRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req VerifyOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo della richiesta non valido", http.StatusBadRequest)
		return
	}

	var userID, storedCode string
	var expiresAt time.Time
	query := `SELECT id, otp_code, otp_expires_at FROM users WHERE email = $1 AND deleted_at IS NULL`
	err := database.Pool.QueryRow(r.Context(), query, req.Email).Scan(&userID, &storedCode, &expiresAt)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	if storedCode == "" || req.Code != storedCode {
		http.Error(w, "Codice non valido", http.StatusBadRequest)
		return
	}
	if time.Now().After(expiresAt) {
		http.Error(w, "Codice scaduto, richiedine uno nuovo", http.StatusBadRequest)
		return
	}

	_, err = database.Pool.Exec(r.Context(),
		`UPDATE users SET email_verified = true, otp_code = NULL, otp_expires_at = NULL WHERE id = $1`,
		userID,
	)
	if err != nil {
		http.Error(w, "Errore interno", http.StatusInternalServerError)
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

type ResendOTPRequest struct {
	Email string `json:"email"`
}

func ResendOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req ResendOTPRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo della richiesta non valido", http.StatusBadRequest)
		return
	}

	var userID string
	var verified bool
	err := database.Pool.QueryRow(r.Context(),
		`SELECT id, email_verified FROM users WHERE email = $1 AND deleted_at IS NULL`,
		req.Email,
	).Scan(&userID, &verified)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}
	if verified {
		http.Error(w, "Email già verificata", http.StatusBadRequest)
		return
	}

	otp := generateOTP()
	expiresAt := time.Now().Add(10 * time.Minute)

	_, err = database.Pool.Exec(r.Context(),
		`UPDATE users SET otp_code = $1, otp_expires_at = $2 WHERE id = $3`,
		otp, expiresAt, userID,
	)
	if err != nil {
		http.Error(w, "Errore interno", http.StatusInternalServerError)
		return
	}

	if err := email.SendOTPEmail(req.Email, otp); err != nil {
		http.Error(w, "Invio email fallito", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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
	var emailVerified bool
	query := `SELECT id, password_hash, email_verified FROM users WHERE email = $1 AND deleted_at IS NULL`
	err := database.Pool.QueryRow(r.Context(), query, req.Email).Scan(&userID, &passwordHash, &emailVerified)
	if err != nil {
		http.Error(w, "Credenziali non valide", http.StatusUnauthorized)
		return
	}

	if !auth.CheckPassword(req.Password, passwordHash) {
		http.Error(w, "Credenziali non valide", http.StatusUnauthorized)
		return
	}

	if !emailVerified {
		http.Error(w, "Email non verificata, controlla la tua casella di posta", http.StatusForbidden)
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
	var avatarURL, passwordHash sql.NullString
	query := `SELECT username, email, avatar_url, password_hash FROM users WHERE id = $1`
	err := database.Pool.QueryRow(r.Context(), query, userID).Scan(&username, &email, &avatarURL, &passwordHash)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":           userID,
		"username":     username,
		"email":        email,
		"avatar_url":   avatarURL.String,
		"has_password": passwordHash.Valid && passwordHash.String != "",
	})
}

func ChangePasswordHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Corpo richiesta non valido", http.StatusBadRequest)
		return
	}
	if len(body.NewPassword) < 6 {
		http.Error(w, "La nuova password deve avere almeno 6 caratteri", http.StatusBadRequest)
		return
	}

	var existingHash sql.NullString
	err := database.Pool.QueryRow(r.Context(),
		`SELECT password_hash FROM users WHERE id = $1`, userID,
	).Scan(&existingHash)
	if err != nil {
		http.Error(w, "Utente non trovato", http.StatusNotFound)
		return
	}

	if existingHash.Valid && existingHash.String != "" {
		if body.CurrentPassword == "" || !auth.CheckPassword(body.CurrentPassword, existingHash.String) {
			http.Error(w, "Password attuale errata", http.StatusForbidden)
			return
		}
	}

	// Controllo contro le ultime 5 password (inclusa quella attuale)
	rows, err := database.Pool.Query(r.Context(),
		`SELECT password_hash FROM password_history WHERE user_id = $1 ORDER BY created_at DESC LIMIT 4`, userID)
	if err != nil {
		http.Error(w, "Errore interno", http.StatusInternalServerError)
		return
	}
	var oldHashes []string
	for rows.Next() {
		var h string
		if err := rows.Scan(&h); err == nil {
			oldHashes = append(oldHashes, h)
		}
	}
	rows.Close()

	if existingHash.Valid && existingHash.String != "" {
		oldHashes = append(oldHashes, existingHash.String)
	}
	for _, h := range oldHashes {
		if auth.CheckPassword(body.NewPassword, h) {
			http.Error(w, "Non puoi riutilizzare una delle ultime 5 password", http.StatusConflict)
			return
		}
	}

	newHash, err := auth.HashPassword(body.NewPassword)
	if err != nil {
		http.Error(w, "Errore interno", http.StatusInternalServerError)
		return
	}

	tx, err := database.Pool.Begin(r.Context())
	if err != nil {
		http.Error(w, "Errore interno", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	if existingHash.Valid && existingHash.String != "" {
		_, err = tx.Exec(r.Context(),
			`INSERT INTO password_history (user_id, password_hash) VALUES ($1, $2)`, userID, existingHash.String)
		if err != nil {
			http.Error(w, "Errore interno", http.StatusInternalServerError)
			return
		}
	}

	_, err = tx.Exec(r.Context(), `UPDATE users SET password_hash = $1 WHERE id = $2`, newHash, userID)
	if err != nil {
		http.Error(w, "Errore aggiornamento password", http.StatusInternalServerError)
		return
	}

	// mantieni solo le ultime 5 nello storico
	_, err = tx.Exec(r.Context(), `
		DELETE FROM password_history
		WHERE user_id = $1 AND id NOT IN (
			SELECT id FROM password_history WHERE user_id = $1 ORDER BY created_at DESC LIMIT 5
		)`, userID)
	if err != nil {
		http.Error(w, "Errore interno", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		http.Error(w, "Errore interno", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func UpdateAvatarHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	var body struct {
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Corpo richiesta non valido", http.StatusBadRequest)
		return
	}

	_, err := database.Pool.Exec(r.Context(),
		`UPDATE users SET avatar_url = $1 WHERE id = $2`, body.AvatarURL, userID)
	if err != nil {
		http.Error(w, "Errore aggiornamento avatar", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}