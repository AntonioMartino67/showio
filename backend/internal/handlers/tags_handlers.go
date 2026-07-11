package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/AntonioMartino67/showio/backend/internal/auth"
	"github.com/AntonioMartino67/showio/backend/internal/database"
)

type Tag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type CreateTagRequest struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// ListTagsHandler restituisce tutti i tag creati dall'utente
func ListTagsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	rows, err := database.Pool.Query(r.Context(),
		`SELECT id, name, color FROM tags WHERE user_id = $1 ORDER BY name`, userID)
	if err != nil {
		http.Error(w, "Errore durante il recupero dei tag", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	tags := []Tag{}
	for rows.Next() {
		var t Tag
		if rows.Scan(&t.ID, &t.Name, &t.Color) == nil {
			tags = append(tags, t)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}

// CreateTagHandler crea un nuovo tag personale
func CreateTagHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)

	var req CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
		http.Error(w, "Nome tag obbligatorio", http.StatusBadRequest)
		return
	}
	if req.Color == "" {
		req.Color = "#00d4ff"
	}

	var id string
	err := database.Pool.QueryRow(r.Context(),
		`INSERT INTO tags (user_id, name, color) VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, name) DO UPDATE SET color = $3
		 RETURNING id`,
		userID, req.Name, req.Color,
	).Scan(&id)
	if err != nil {
		http.Error(w, "Errore durante la creazione del tag", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Tag{ID: id, Name: req.Name, Color: req.Color})
}

// DeleteTagHandler elimina un tag (e le sue associazioni, via cascade)
func DeleteTagHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	tagID := chi.URLParam(r, "tagId")

	tag, err := database.Pool.Exec(r.Context(),
		`DELETE FROM tags WHERE id = $1 AND user_id = $2`, tagID, userID)
	if err != nil || tag.RowsAffected() == 0 {
		http.Error(w, "Tag non trovato", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type AssignTagRequest struct {
	TagID string `json:"tag_id"`
}

// AssignTagHandler associa un tag a un titolo della lista dell'utente
func AssignTagHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	mediaItemID := chi.URLParam(r, "mediaItemId")

	var req AssignTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TagID == "" {
		http.Error(w, "tag_id obbligatorio", http.StatusBadRequest)
		return
	}

	var progressID string
	err := database.Pool.QueryRow(r.Context(),
		`SELECT id FROM user_progress WHERE user_id = $1 AND media_item_id = $2`,
		userID, mediaItemID,
	).Scan(&progressID)
	if err != nil {
		http.Error(w, "Titolo non presente nella lista", http.StatusNotFound)
		return
	}

	_, err = database.Pool.Exec(r.Context(),
		`INSERT INTO progress_tags (progress_id, tag_id) VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`, progressID, req.TagID)
	if err != nil {
		http.Error(w, "Errore durante l'associazione del tag", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveTagFromProgressHandler rimuove un tag da un titolo
func RemoveTagFromProgressHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(auth.UserIDKey).(string)
	mediaItemID := chi.URLParam(r, "mediaItemId")
	tagID := chi.URLParam(r, "tagId")

	_, err := database.Pool.Exec(r.Context(), `
		DELETE FROM progress_tags
		WHERE tag_id = $1 AND progress_id = (
			SELECT id FROM user_progress WHERE user_id = $2 AND media_item_id = $3
		)`, tagID, userID, mediaItemID)
	if err != nil {
		http.Error(w, "Errore durante la rimozione del tag", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}