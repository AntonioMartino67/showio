package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const anilistURL = "https://graphql.anilist.co"

type AniListResult struct {
	ID    int `json:"id"`
	Title struct {
		Romaji  string `json:"romaji"`
		English string `json:"english"`
	} `json:"title"`
	Description string `json:"description"`
	CoverImage  struct {
		Large string `json:"large"`
	} `json:"coverImage"`
	Status string `json:"status"`
}

type anilistGraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type anilistResponse struct {
	Data struct {
		Page struct {
			Media []AniListResult `json:"media"`
		} `json:"Page"`
	} `json:"data"`
}

func SearchAnime(query string) ([]AniListResult, error) {
	graphqlQuery := `
		query ($search: String) {
			Page(page: 1, perPage: 10) {
				media(search: $search, type: ANIME) {
					id
					title {
						romaji
						english
					}
					description
					coverImage {
						large
					}
					status
				}
			}
		}
	`

	reqBody := anilistGraphQLRequest{
		Query:     graphqlQuery,
		Variables: map[string]interface{}{"search": query},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(anilistURL, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AniList ha risposto con status %d", resp.StatusCode)
	}

	var result anilistResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data.Page.Media, nil
}