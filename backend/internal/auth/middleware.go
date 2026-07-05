package auth

import (
	"context"
	"net/http"
	"strings"
)

// contextKey evita collisioni con altre chiavi eventualmente usate nel context
type contextKey string

const UserIDKey contextKey = "user_id"

// RequireAuth è un middleware chi-compatibile: controlla il JWT e, se valido,
// inserisce lo user_id nel context della richiesta per gli handler successivi
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Token di autenticazione mancante", http.StatusUnauthorized)
			return
		}

		// Ci aspettiamo il formato "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Formato del token non valido", http.StatusUnauthorized)
			return
		}

		userID, err := ParseJWT(parts[1])
		if err != nil {
			http.Error(w, "Token non valido o scaduto", http.StatusUnauthorized)
			return
		}

		// Rendiamo disponibile lo user_id agli handler successivi tramite il context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}