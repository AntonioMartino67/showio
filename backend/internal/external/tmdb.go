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