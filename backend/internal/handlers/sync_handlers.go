package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/AntonioMartino67/showio/backend/internal/database"
	"github.com/AntonioMartino67/showio/backend/internal/email"
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

// getLastKnownEpisode ritorna l'ultima stagione/episodio noto per un MediaItem
func getLastKnownEpisode(mediaItemID string) (int, int) {
	var season, episode int
	database.Pool.QueryRow(context.Background(), `
		SELECT s.season_number, e.episode_number
		FROM episodes e
		JOIN seasons s ON s.id = e.season_id
		WHERE s.media_item_id = $1
		ORDER BY s.season_number DESC, e.episode_number DESC
		LIMIT 1
	`, mediaItemID).Scan(&season, &episode)
	return season, episode
}

// fixCompletedStatus riporta a "watching" chi era rimasto a "completed"
// ma ora è indietro rispetto al nuovo ultimo episodio noto
func fixCompletedStatus(mediaItemID string, newSeason, newEpisode int) {
	database.Pool.Exec(context.Background(), `
		UPDATE user_progress
		SET status = 'watching'
		WHERE media_item_id = $1 AND status = 'completed'
		AND (current_season < $2 OR (current_season = $2 AND current_episode < $3))
	`, mediaItemID, newSeason, newEpisode)
}

// notifyNewContent avvisa via email chi sta guardando (watching) questo MediaItem
func notifyNewContent(mediaItemID string, seasonNumber int) {
	var title string
	database.Pool.QueryRow(context.Background(),
		`SELECT title FROM media_items WHERE id = $1`, mediaItemID,
	).Scan(&title)

	rows, err := database.Pool.Query(context.Background(),
		`SELECT u.email, u.username FROM users u
		 JOIN user_progress up ON up.user_id = u.id
		 WHERE up.media_item_id = $1 AND up.status = 'watching'
		 AND u.notify_new_seasons = true AND u.deleted_at IS NULL`,
		mediaItemID,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var to, username string
		if err := rows.Scan(&to, &username); err != nil {
			continue
		}
		go email.SendNewSeasonEmail(to, username, title, seasonNumber)
	}
}

// SyncAllHandler risincronizza stagioni/episodi di tutte le serie TV
// tracciate da almeno un utente in "watching" o "completed" (chi ha finito
// gli episodi usciti finora va comunque ricontrollato per i nuovi).
// Protetto da una chiave segreta passata nell'header X-Cron-Secret,
// pensato per essere chiamato da un job schedulato.
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
		 WHERE mi.type = 'tv' AND mi.source = 'tmdb' AND up.status IN ('watching', 'completed')`,
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
		oldSeason, oldEpisode := getLastKnownEpisode(item.MediaItemID)

		if err := SyncTVSeasons(item.MediaItemID, item.ExternalID); err != nil {
			continue
		}
		synced++

		newSeason, newEpisode := getLastKnownEpisode(item.MediaItemID)
		if newSeason > oldSeason || (newSeason == oldSeason && newEpisode > oldEpisode) {
			fixCompletedStatus(item.MediaItemID, newSeason, newEpisode)
			notifyNewContent(item.MediaItemID, newSeason)
		}
	}

	// Auto-drop: segna come "dropped" i titoli in "watching" fermi da troppo tempo
	autoDropped := 0
	dropTag, err := database.Pool.Exec(context.Background(),
		`UPDATE user_progress
		 SET status = 'dropped'
		 WHERE status = 'watching' AND last_watched_at < NOW() - INTERVAL '30 days'`,
	)
	if err == nil {
		autoDropped = int(dropTag.RowsAffected())
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"serie_sincronizzate": synced,
		"totale_trovate":      len(items),
		"auto_droppati":       autoDropped,
	})
}