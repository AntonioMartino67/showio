package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/AntonioMartino67/showio/backend/internal/database"
	"github.com/AntonioMartino67/showio/backend/internal/external"
)

type UnifiedResult struct {
	MediaItemID string `json:"media_item_id"`
	ExternalID  string `json:"external_id"`
	Source      string `json:"source"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	PosterURL   string `json:"poster_url"`
	Overview    string `json:"overview"`
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Parametro 'q' mancante", http.StatusBadRequest)
		return
	}

	var unified []UnifiedResult

	tmdbResults, err := external.SearchMulti(query)
	if err != nil {
		log.Printf("[search] errore TMDB per query %q: %v", query, err)
	} else {
		for _, item := range tmdbResults {
			if item.MediaType != "movie" && item.MediaType != "tv" {
				continue
			}
			title := item.Title
			if title == "" {
				title = item.Name
			}
			unified = append(unified, UnifiedResult{
				ExternalID: fmt.Sprintf("%d", item.ID),
				Source:     "tmdb",
				Title:      title,
				Type:       item.MediaType,
				PosterURL:  item.PosterPath,
				Overview:   item.Overview,
			})
		}
	}

	anilistResults, err := external.SearchAnime(query)
	if err != nil {
		log.Printf("[search] errore AniList per query %q: %v", query, err)
	} else {
		for _, item := range anilistResults {
			title := item.Title.English
			if title == "" {
				title = item.Title.Romaji
			}
			unified = append(unified, UnifiedResult{
				ExternalID: fmt.Sprintf("%d", item.ID),
				Source:     "anilist",
				Title:      title,
				Type:       "anime",
				PosterURL:  item.CoverImage.Large,
				Overview:   item.Description,
			})
		}
	}

	for i, item := range unified {
		q := `
			INSERT INTO media_items (external_id, source, title, type, poster_url, overview)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (external_id, source)
			DO UPDATE SET title = $3, poster_url = $5, overview = $6, last_synced_at = CURRENT_TIMESTAMP
			RETURNING id
		`
		var id string
		err := database.Pool.QueryRow(r.Context(), q,
			item.ExternalID, item.Source, item.Title, item.Type, item.PosterURL, item.Overview).Scan(&id)
		if err == nil {
			unified[i].MediaItemID = id
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(unified)
}