package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/AntonioMartino67/showio/backend/internal/auth"
	"github.com/AntonioMartino67/showio/backend/internal/database"
)

type EpisodeDetail struct {
	SeasonNumber  int     `json:"season_number"`
	EpisodeNumber int     `json:"episode_number"`
	Title         *string `json:"title,omitempty"`
	AirDate       *string `json:"air_date,omitempty"`
}

type MediaDetail struct {
	MediaItemID    string          `json:"media_item_id"`
	Title          string          `json:"title"`
	Type           string          `json:"type"`
	PosterURL      *string         `json:"poster_url,omitempty"`
	Overview       *string         `json:"overview,omitempty"`
	Status         *string         `json:"status,omitempty"` // in lista: watching/completed/dropped/plan_to_watch, nil se non in lista
	CurrentSeason  int             `json:"current_season"`
	CurrentEpisode int             `json:"current_episode"`
	Rating         *int            `json:"rating,omitempty"`
	Episodes       []EpisodeDetail `json:"episodes"`
	Tags           []Tag           `json:"tags"`
}

// MediaDetailHandler restituisce i dettagli di un titolo, incluse le stagioni/episodi
// e il progresso dell'utente corrente (se il titolo è già in lista)
func MediaDetailHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	mediaItemID := chi.URLParam(r, "mediaItemId")

	var detail MediaDetail
	err := database.Pool.QueryRow(r.Context(),
		`SELECT id, title, type, poster_url, overview FROM media_items WHERE id = $1`, mediaItemID,
	).Scan(&detail.MediaItemID, &detail.Title, &detail.Type, &detail.PosterURL, &detail.Overview)
	if err != nil {
		http.Error(w, "Titolo non trovato", http.StatusNotFound)
		return
	}

	var status *string
	var progressID *string
	err = database.Pool.QueryRow(r.Context(),
		`SELECT id, status, current_season, current_episode, rating FROM user_progress WHERE user_id = $1 AND media_item_id = $2`,
		userID, mediaItemID,
	).Scan(&progressID, &status, &detail.CurrentSeason, &detail.CurrentEpisode, &detail.Rating)
	if err == nil {
		detail.Status = status
	}

	detail.Tags = []Tag{}
	if progressID != nil {
		tagRows, err := database.Pool.Query(r.Context(), `
			SELECT t.id, t.name, t.color FROM tags t
			JOIN progress_tags pt ON pt.tag_id = t.id
			WHERE pt.progress_id = $1
			ORDER BY t.name
		`, *progressID)
		if err == nil {
			defer tagRows.Close()
			for tagRows.Next() {
				var t Tag
				if tagRows.Scan(&t.ID, &t.Name, &t.Color) == nil {
					detail.Tags = append(detail.Tags, t)
				}
			}
		}
	}

	detail.Episodes = []EpisodeDetail{}
	rows, err := database.Pool.Query(r.Context(),
		`SELECT s.season_number, e.episode_number, e.title, e.air_date::text
		 FROM episodes e JOIN seasons s ON s.id = e.season_id
		 WHERE s.media_item_id = $1
		 ORDER BY s.season_number, e.episode_number`, mediaItemID,
	)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ep EpisodeDetail
			if rows.Scan(&ep.SeasonNumber, &ep.EpisodeNumber, &ep.Title, &ep.AirDate) == nil {
				detail.Episodes = append(detail.Episodes, ep)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(detail)
}