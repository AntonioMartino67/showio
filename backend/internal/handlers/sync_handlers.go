package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/AntonioMartino67/showio/backend/internal/database"
	"github.com/AntonioMartino67/showio/backend/internal/external"
)

// SyncTVSeasons scarica da TMDB tutte le stagioni ed episodi di una serie
// e li salva/aggiorna nel DB locale (tabelle seasons ed episodes)
func SyncTVSeasons(mediaItemID string, tmdbExternalID string) error {
	tmdbID, err := strconv.Atoi(tmdbExternalID)
	if err != nil {
		return err
	}

	details, err := external.GetTVDetails(tmdbID)
	if err != nil {
		return err
	}

	for seasonNum := 1; seasonNum <= details.NumberOfSeasons; seasonNum++ {
		episodes, err := external.GetSeasonEpisodes(tmdbID, seasonNum)
		if err != nil {
			continue // se una stagione fallisce, proviamo comunque le altre
		}

		var seasonID string
		err = database.Pool.QueryRow(context.Background(),
			`INSERT INTO seasons (media_item_id, season_number, episode_count)
			 VALUES ($1, $2, $3)
			 ON CONFLICT (media_item_id, season_number)
			 DO UPDATE SET episode_count = $3
			 RETURNING id`,
			mediaItemID, seasonNum, len(episodes),
		).Scan(&seasonID)
		if err != nil {
			continue
		}

		for _, ep := range episodes {
			var airDate interface{}
			if ep.AirDate == "" {
				airDate = nil
			} else {
				airDate = ep.AirDate
			}

			database.Pool.Exec(context.Background(),
				`INSERT INTO episodes (season_id, episode_number, title, air_date)
				 VALUES ($1, $2, $3, $4)
				 ON CONFLICT (season_id, episode_number)
				 DO UPDATE SET title = $3, air_date = $4`,
				seasonID, ep.EpisodeNumber, ep.Name, airDate,
			)
		}
	}

	return nil
}

// SyncAllHandler risincronizza stagioni/episodi di tutte le serie TV
// attualmente tracciate da almeno un utente. Protetto da una chiave segreta
// passata nell'header X-Cron-Secret, pensato per essere chiamato da un job schedulato.
func SyncAllHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("CRON_SECRET")
	if secret == "" || r.Header.Get("X-Cron-Secret") != secret {
		http.Error(w, "Non autorizzato", http.StatusUnauthorized)
		return
	}

	rows, err := database.Pool.Query(context.Background(),
		`SELECT DISTINCT mi.id, mi.external_id
		 FROM media_items mi
		 JOIN user_progress up ON up.media_item_id = mi.id
		 WHERE mi.type = 'tv' AND mi.source = 'tmdb' AND up.status = 'watching'`,
	)
	if err != nil {
		http.Error(w, "Errore durante il recupero delle serie", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type tvItem struct {
		MediaItemID string
		ExternalID  string
	}
	var items []tvItem
	for rows.Next() {
		var item tvItem
		if err := rows.Scan(&item.MediaItemID, &item.ExternalID); err != nil {
			continue
		}
		items = append(items, item)
	}
	rows.Close()

	synced := 0
	for _, item := range items {
		if err := SyncTVSeasons(item.MediaItemID, item.ExternalID); err == nil {
			synced++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"serie_sincronizzate": synced,
		"totale_trovate":      len(items),
	})
}