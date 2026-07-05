package models

import "time"

// User rappresenta un utente registrato
type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"` // "-" esclude questo campo dalle risposte JSON, non va mai esposto al frontend
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

// MediaItem rappresenta un contenuto (serie, film, anime) salvato in cache locale
type MediaItem struct {
	ID            string     `json:"id"`
	ExternalID    string     `json:"external_id"`
	Source        string     `json:"source"` // "tmdb" o "anilist"
	Title         string     `json:"title"`
	Type          string     `json:"type"` // "tv", "movie", "anime"
	PosterURL     *string    `json:"poster_url,omitempty"`
	Overview      *string    `json:"overview,omitempty"`
	Status        *string    `json:"status,omitempty"` // "airing", "ended", "upcoming"
	LastSyncedAt  *time.Time `json:"last_synced_at,omitempty"`
}

// Season rappresenta una stagione di un MediaItem
type Season struct {
	ID            string     `json:"id"`
	MediaItemID   string     `json:"media_item_id"`
	SeasonNumber  int        `json:"season_number"`
	EpisodeCount  *int       `json:"episode_count,omitempty"`
	AirDate       *time.Time `json:"air_date,omitempty"`
}

// Episode rappresenta un episodio di una Season
type Episode struct {
	ID            string     `json:"id"`
	SeasonID      string     `json:"season_id"`
	EpisodeNumber int        `json:"episode_number"`
	Title         *string    `json:"title,omitempty"`
	AirDate       *time.Time `json:"air_date,omitempty"`
	IsWatched     bool       `json:"is_watched"`
}

// UserProgress rappresenta il progresso di un utente su un MediaItem
type UserProgress struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	MediaItemID     string     `json:"media_item_id"`
	Status          string     `json:"status"` // "watching", "completed", "dropped", "plan_to_watch"
	CurrentSeason   int        `json:"current_season"`
	CurrentEpisode  int        `json:"current_episode"`
	Rating          *int       `json:"rating,omitempty"`
	LastWatchedAt   *time.Time `json:"last_watched_at,omitempty"`
}