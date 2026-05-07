package main

import (
	"net/http"
	"sort"
	"time"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	s := r.URL.Query().Get("author_id")
	if s != "" {
		userID, err := uuid.Parse(s)
		if err != nil {
			respondWithError(w, 400, "Bad Request")
			return
		}
		chirps, err := cfg.db.GetChirpsByAuthor(r.Context(), userID)
		if err != nil {
			respondWithError(w, 500, "Cannot Retrieve Chirps")
			return
		}
		sortParam := "asc"
		if r.URL.Query().Get("sort") == "desc" {
			sortParam = "desc"
		}
		var response []responseBody
		for _, chirp := range chirps {
			response = append(response, responseBody{
				ID:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			})
		}
		sort.Slice(response, func(i, j int) bool {
			if sortParam == "desc" {
				return response[i].CreatedAt.After(response[j].CreatedAt)
			}
			return response[i].CreatedAt.Before(response[j].CreatedAt)
		})
		respondWithJSON(w, 200, response)
		return
	}
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, 500, "Cannot Retrieve Chirps")
		return
	}
	sortParam := "asc"
	if r.URL.Query().Get("sort") == "desc" {
		sortParam = "desc"
	}
	var response []responseBody
	for _, chirp := range chirps {
		response = append(response, responseBody{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	sort.Slice(response, func(i, j int) bool {
		if sortParam == "desc" {
			return response[i].CreatedAt.After(response[j].CreatedAt)
		}
		return response[i].CreatedAt.Before(response[j].CreatedAt)
	})
	respondWithJSON(w, 200, response)
}

func (cfg *apiConfig) handlerGetChirpbyId(w http.ResponseWriter, r *http.Request) {
	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Cannot Find Chirp ID")
		return
	}
	response := responseBody{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, 200, response)
}
