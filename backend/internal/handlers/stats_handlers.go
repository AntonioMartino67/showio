package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AntonioMartino67/showio/backend/internal/auth"
	"github.com/AntonioMartino67/showio/backend/internal/database"
)

type StatsResponse struct {
	TotalTitles      int            `json:"total_titles"`
	ByStatus         map[string]int `json:"by_status"`
	ByType           map[string]int `json:"by_type"`
	TotalEpisodes    int            `json:"total_episodes_watched"`
	AverageRating    *float64       `json:"average_rating,omitempty"`
	RatedTitlesCount int            `json:"rated_titles_count"`
}

// StatsHandler restituisce statistiche aggregate sulla libreria dell'utente
func StatsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	resp := StatsResponse{
		ByStatus: map[string]int{},
		ByType:   map[string]int{},
	}

	rows, err := database.Pool.Query(r.Context(), `
		SELECT up.status, mi.type, up.current_episode, up.rating
		FROM user_progress up
		JOIN media_items mi ON mi.id = up.media_item_id
		WHERE up.user_id = $1
	`, userID)
	if err != nil {
		http.Error(w, "Errore durante il recupero delle statistiche", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ratingSum, ratingCount int
	for rows.Next() {
		var status, mediaType string
		var currentEpisode int
		var rating *int
		if err := rows.Scan(&status, &mediaType, &currentEpisode, &rating); err != nil {
			continue
		}
		resp.TotalTitles++
		resp.ByStatus[status]++
		resp.ByType[mediaType]++
		resp.TotalEpisodes += currentEpisode
		if rating != nil {
			ratingSum += *rating
			ratingCount++
		}
	}

	resp.RatedTitlesCount = ratingCount
	if ratingCount > 0 {
		avg := float64(ratingSum) / float64(ratingCount)
		resp.AverageRating = &avg
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}