package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const tmdbBaseURL = "https://api.themoviedb.org/3"

type TMDBSearchResult struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`        // presente per i film
	Name         string `json:"name"`         // presente per le serie TV
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	MediaType    string `json:"media_type"`   // "movie" o "tv"
	ReleaseDate  string `json:"release_date"`
	FirstAirDate string `json:"first_air_date"`
}

type tmdbSearchResponse struct {
	Results []TMDBSearchResult `json:"results"`
}

type TMDBTVDetails struct {
	ID              int `json:"id"`
	NumberOfSeasons int `json:"number_of_seasons"`
}

type TMDBEpisode struct {
	EpisodeNumber int    `json:"episode_number"`
	Name          string `json:"name"`
	AirDate       string `json:"air_date"`
}

type tmdbSeasonResponse struct {
	Episodes []TMDBEpisode `json:"episodes"`
}

// SearchMulti cerca film e serie TV contemporaneamente su TMDB
func SearchMulti(query string) ([]TMDBSearchResult, error) {
	token := os.Getenv("TMDB_TOKEN")

	fullURL := fmt.Sprintf("%s/search/multi?query=%s&language=it-IT", tmdbBaseURL, url.QueryEscape(query))

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB ha risposto con status %d", resp.StatusCode)
	}

	var result tmdbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}

// GetTVDetails recupera i dettagli di una serie TV, incluso il numero di stagioni
func GetTVDetails(tmdbID int) (TMDBTVDetails, error) {
	token := os.Getenv("TMDB_TOKEN")

	fullURL := fmt.Sprintf("%s/tv/%d?language=it-IT", tmdbBaseURL, tmdbID)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return TMDBTVDetails{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return TMDBTVDetails{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TMDBTVDetails{}, fmt.Errorf("TMDB ha risposto con status %d", resp.StatusCode)
	}

	var result TMDBTVDetails
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return TMDBTVDetails{}, err
	}

	return result, nil
}

// GetSeasonEpisodes recupera l'elenco degli episodi di una specifica stagione
func GetSeasonEpisodes(tmdbID int, seasonNumber int) ([]TMDBEpisode, error) {
	token := os.Getenv("TMDB_TOKEN")

	fullURL := fmt.Sprintf("%s/tv/%d/season/%d?language=it-IT", tmdbBaseURL, tmdbID, seasonNumber)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB ha risposto con status %d", resp.StatusCode)
	}

	var result tmdbSeasonResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Episodes, nil
}

// GetTrending recupera i titoli di tendenza della settimana (film + serie TV)
func GetTrending() ([]TMDBSearchResult, error) {
	token := os.Getenv("TMDB_TOKEN")

	fullURL := fmt.Sprintf("%s/trending/all/week?language=it-IT", tmdbBaseURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TMDB ha risposto con status %d", resp.StatusCode)
	}

	var result tmdbSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Results, nil
}