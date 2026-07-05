package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AntonioMartino67/showio/backend/internal/database"
	"github.com/AntonioMartino67/showio/backend/internal/external"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Parametro 'q' mancante", http.StatusBadRequest)
		return
	}

	results, err := external.SearchMulti(query)
	if err != nil {
		http.Error(w, "Errore durante la ricerca su TMDB", http.StatusBadGateway)
		return
	}

	for _, item := range results {
		if item.MediaType != "movie" && item.MediaType != "tv" {
			continue
		}

		title := item.Title
		if title == "" {
			title = item.Name
		}

		mediaType := "movie"
		if item.MediaType == "tv" {
			mediaType = "tv"
		}

		q := `
			INSERT INTO media_items (external_id, source, title, type, poster_url, overview)
			VALUES ($1, 'tmdb', $2, $3, $4, $5)
			ON CONFLICT (external_id, source)
			DO UPDATE SET title = $2, poster_url = $4, overview = $5, last_synced_at = CURRENT_TIMESTAMP
		`
		database.Pool.Exec(r.Context(), q,
			fmt.Sprintf("%d", item.ID), title, mediaType, item.PosterPath, item.Overview)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}