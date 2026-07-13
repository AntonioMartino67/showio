package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/AntonioMartino67/showio/backend/internal/auth"
	"github.com/AntonioMartino67/showio/backend/internal/database"
)

type GoogleLoginRequest struct {
	Credential string `json:"credential"`
}

type googleTokenInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	Name          string `json:"name"`
	Aud           string `json:"aud"`
}

// GoogleLoginHandler autentica (o registra al volo) un utente a partire dal
// token ID restituito da Google Identity Services sul frontend.
func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req GoogleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Credential == "" {
		http.Error(w, "Credential mancante", http.StatusBadRequest)
		return
	}

	info, err := verifyGoogleToken(req.Credential)
	if err != nil {
		http.Error(w, "Token Google non valido", http.StatusUnauthorized)
		return
	}

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" || info.Aud != clientID {
		http.Error(w, "Token Google non valido", http.StatusUnauthorized)
		return
	}
	if info.EmailVerified != "true" {
		http.Error(w, "Email Google non verificata", http.StatusUnauthorized)
		return
	}

	// L'utente ha già effettuato l'accesso con Google in passato?
	var userID string
	err = database.Pool.QueryRow(r.Context(),
		`SELECT id FROM users WHERE google_id = $1`, info.Sub,
	).Scan(&userID)

	if err != nil {
		// Nessun account collegato a questo Google ID: se esiste già un utente
		// con la stessa email (registrato con password) lo colleghiamo,
		// altrimenti creiamo un nuovo account.
		err = database.Pool.QueryRow(r.Context(),
			`UPDATE users SET google_id = $1 WHERE email = $2 RETURNING id`,
			info.Sub, info.Email,
		).Scan(&userID)

		if err != nil {
			username := generateUsernameFromEmail(r, info.Email)
			err = database.Pool.QueryRow(r.Context(),
				`INSERT INTO users (username, email, google_id) VALUES ($1, $2, $3) RETURNING id`,
				username, info.Email, info.Sub,
			).Scan(&userID)
			if err != nil {
				http.Error(w, "Errore durante la creazione dell'account", http.StatusInternalServerError)
				return
			}
		}
	}

	token, err := auth.GenerateJWT(userID)
	if err != nil {
		http.Error(w, "Errore interno durante il login", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// verifyGoogleToken valida il token ID interrogando l'endpoint tokeninfo di Google.
// Non richiede librerie esterne: Google firma il token e questo endpoint
// verifica firma/scadenza per noi, restituendo i claim in chiaro.
func verifyGoogleToken(credential string) (*googleTokenInfo, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + credential)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token non valido, status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info googleTokenInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// generateUsernameFromEmail deriva uno username dalla parte locale dell'email,
// aggiungendo un suffisso numerico in caso di collisione
func generateUsernameFromEmail(r *http.Request, email string) string {
	base := strings.Split(email, "@")[0]
	username := base
	for i := 1; i < 50; i++ {
		var exists bool
		err := database.Pool.QueryRow(r.Context(),
			`SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`, username,
		).Scan(&exists)
		if err != nil || !exists {
			return username
		}
		username = fmt.Sprintf("%s%d", base, i)
	}
	return fmt.Sprintf("%s%d", base, os.Getpid())
}