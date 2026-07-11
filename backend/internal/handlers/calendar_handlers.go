package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/AntonioMartino67/showio/backend/internal/auth"
	"github.com/AntonioMartino67/showio/backend/internal/database"
)

type UpcomingEpisode struct {
	MediaItemID   string  `json:"media_item_id"`
	Title         string  `json:"title"`
	PosterURL     *string `json:"poster_url,omitempty"`
	SeasonNumber  int     `json:"season_number"`
	EpisodeNumber int     `json:"episode_number"`
	AirDate       *string `json:"air_date,omitempty"`
}

// CalendarHandler restituisce, per ogni serie che l'utente sta seguendo,
// il prossimo episodio non ancora visto
func CalendarHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	query := `
		SELECT DISTINCT ON (up.media_item_id)
			mi.id, mi.title, mi.poster_url,
			s.season_number, e.episode_number, e.air_date::text
		FROM user_progress up
		JOIN media_items mi ON mi.id = up.media_item_id
		JOIN seasons s ON s.media_item_id = mi.id
		JOIN episodes e ON e.season_id = s.id
		WHERE up.user_id = $1
		  AND up.status = 'watching'
		  AND (s.season_number, e.episode_number) > (up.current_season, up.current_episode)
		ORDER BY up.media_item_id, s.season_number, e.episode_number
	`

	rows, err := database.Pool.Query(r.Context(), query, userID)
	if err != nil {
		http.Error(w, "Errore durante il recupero del calendario", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	results := []UpcomingEpisode{}
	for rows.Next() {
		var item UpcomingEpisode
		if err := rows.Scan(
			&item.MediaItemID, &item.Title, &item.PosterURL,
			&item.SeasonNumber, &item.EpisodeNumber, &item.AirDate,
		); err != nil {
			http.Error(w, "Errore durante la lettura dei dati", http.StatusInternalServerError)
			return
		}
		results = append(results, item)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}