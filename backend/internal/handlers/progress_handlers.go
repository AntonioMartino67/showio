package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"encoding/json"
	"net/http"

	"github.com/AntonioMartino67/showio/backend/internal/auth"
	"github.com/AntonioMartino67/showio/backend/internal/database"
)

type AddProgressRequest struct {
	MediaItemID string `json:"media_item_id"`
	Status      string `json:"status"` // "watching", "completed", "dropped", "plan_to_watch"
}

// AddProgressHandler aggiunge un contenuto alla lista personale dell'utente
// (o aggiorna lo status se è già presente)
func AddProgressHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	var req AddProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo della richiesta non valido", http.StatusBadRequest)
		return
	}

	if req.MediaItemID == "" {
		http.Error(w, "media_item_id obbligatorio", http.StatusBadRequest)
		return
	}
	if req.Status == "" {
		req.Status = "plan_to_watch"
	}

	validStatuses := map[string]bool{
		"watching": true, "completed": true, "dropped": true, "plan_to_watch": true,
	}
	if !validStatuses[req.Status] {
		http.Error(w, "Status non valido", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO user_progress (user_id, media_item_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, media_item_id)
		DO UPDATE SET status = $3
		RETURNING id
	`
	var progressID string
	err := database.Pool.QueryRow(r.Context(), query, userID, req.MediaItemID, req.Status).Scan(&progressID)
	if err != nil {
		http.Error(w, "Errore durante il salvataggio", http.StatusInternalServerError)
		return
	}

	// Se il media è una serie TV, sincronizziamo stagioni/episodi da TMDB in background
	go func() {
		var mediaType, externalID string
		err := database.Pool.QueryRow(context.Background(),
			`SELECT type, external_id FROM media_items WHERE id = $1`, req.MediaItemID,
		).Scan(&mediaType, &externalID)
		if err != nil {
			return
		}
		if mediaType == "tv" {
			SyncTVSeasons(req.MediaItemID, externalID)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":     progressID,
		"status": req.Status,
	})
}

type UpdateEpisodeRequest struct {
	CurrentSeason  int `json:"current_season"`
	CurrentEpisode int `json:"current_episode"`
}

// UpdateEpisodeHandler aggiorna a che episodio/stagione è arrivato l'utente
// UpdateEpisodeHandler aggiorna a che episodio/stagione è arrivato l'utente.
// Applica anche una logica di auto-status.
func UpdateEpisodeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	mediaItemID := chi.URLParam(r, "mediaItemId")

	var req UpdateEpisodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo della richiesta non valido", http.StatusBadRequest)
		return
	}

	var currentStatus string
	err := database.Pool.QueryRow(r.Context(),
		`SELECT status FROM user_progress WHERE user_id = $1 AND media_item_id = $2`,
		userID, mediaItemID,
	).Scan(&currentStatus)
	if err != nil {
		http.Error(w, "Progresso non trovato", http.StatusNotFound)
		return
	}

	newStatus := currentStatus
	if currentStatus == "plan_to_watch" {
		newStatus = "watching"
	}

	var lastSeason, lastEpisode int
	err = database.Pool.QueryRow(r.Context(), `
		SELECT s.season_number, e.episode_number
		FROM episodes e
		JOIN seasons s ON s.id = e.season_id
		WHERE s.media_item_id = $1
		ORDER BY s.season_number DESC, e.episode_number DESC
		LIMIT 1
	`, mediaItemID).Scan(&lastSeason, &lastEpisode)
	if err == nil {
		if req.CurrentSeason > lastSeason ||
			(req.CurrentSeason == lastSeason && req.CurrentEpisode >= lastEpisode) {
			newStatus = "completed"
		} else if newStatus == "completed" {
			newStatus = "watching"
		}
	}

	query := `
		UPDATE user_progress
		SET current_season = $1, current_episode = $2, status = $3, last_watched_at = CURRENT_TIMESTAMP
		WHERE user_id = $4 AND media_item_id = $5
	`
	tag, err := database.Pool.Exec(r.Context(), query, req.CurrentSeason, req.CurrentEpisode, newStatus, userID, mediaItemID)
	if err != nil || tag.RowsAffected() == 0 {
		http.Error(w, "Progresso non trovato", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListProgressHandler restituisce l'intera lista personale dell'utente,
// unendo i dati di progresso con i metadati del contenuto (join)
func ListProgressHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	query := `
		SELECT
			up.id, up.status, up.current_season, up.current_episode, up.rating,
			mi.id, mi.title, mi.type, mi.poster_url
		FROM user_progress up
		JOIN media_items mi ON mi.id = up.media_item_id
		WHERE up.user_id = $1
		ORDER BY up.last_watched_at DESC NULLS LAST
	`
	rows, err := database.Pool.Query(r.Context(), query, userID)
	if err != nil {
		http.Error(w, "Errore durante il recupero della lista", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type ProgressItem struct {
		ProgressID     string  `json:"progress_id"`
		Status         string  `json:"status"`
		CurrentSeason  int     `json:"current_season"`
		CurrentEpisode int     `json:"current_episode"`
		Rating         *int    `json:"rating,omitempty"`
		MediaItemID    string  `json:"media_item_id"`
		Title          string  `json:"title"`
		Type           string  `json:"type"`
		PosterURL      *string `json:"poster_url,omitempty"`
	}

	var results []ProgressItem
	for rows.Next() {
		var item ProgressItem
		if err := rows.Scan(
			&item.ProgressID, &item.Status, &item.CurrentSeason, &item.CurrentEpisode, &item.Rating,
			&item.MediaItemID, &item.Title, &item.Type, &item.PosterURL,
		); err != nil {
			http.Error(w, "Errore durante la lettura dei dati", http.StatusInternalServerError)
			return
		}
		results = append(results, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

// RemoveProgressHandler rimuove un titolo dalla lista personale dell'utente
func RemoveProgressHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	mediaItemID := chi.URLParam(r, "mediaItemId")

	query := `DELETE FROM user_progress WHERE user_id = $1 AND media_item_id = $2`
	tag, err := database.Pool.Exec(r.Context(), query, userID, mediaItemID)
	if err != nil {
		http.Error(w, "Errore durante la rimozione", http.StatusInternalServerError)
		return
	}
	if tag.RowsAffected() == 0 {
		http.Error(w, "Elemento non trovato", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type RatingRequest struct {
	Rating int `json:"rating"` // 1-10
}

// UpdateRatingHandler assegna un voto a un titolo della lista
func UpdateRatingHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	mediaItemID := chi.URLParam(r, "mediaItemId")

	var req RatingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Corpo della richiesta non valido", http.StatusBadRequest)
		return
	}
	if req.Rating < 1 || req.Rating > 10 {
		http.Error(w, "Rating deve essere tra 1 e 10", http.StatusBadRequest)
		return
	}

	query := `UPDATE user_progress SET rating = $1 WHERE user_id = $2 AND media_item_id = $3`
	tag, err := database.Pool.Exec(r.Context(), query, req.Rating, userID, mediaItemID)
	if err != nil || tag.RowsAffected() == 0 {
		http.Error(w, "Elemento non trovato", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}