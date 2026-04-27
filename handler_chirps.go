package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Wayne-Francis/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserID    uuid.UUID `json:"user_id"`
	}
	decoder := json.NewDecoder(r.Body)
	request := requestBody{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "couldn't decode parameters")
		return
	}
	if len(request.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	cleanedBody := getCleanedBody(request.Body)
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: request.UserID,
	})
	if err != nil {
		respondWithError(w, 500, "Cannot Generate Chirp")
		return
	}
	new_chirp := responseBody{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}
	respondWithJSON(w, 201, new_chirp)
}

func getCleanedBody(s string) string {
	words := strings.Split(s, " ")
	targets := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range words {
		for _, target := range targets {
			if strings.ToLower(word) == target {
				words[i] = "****"
			}
		}
	}
	cleanbody := strings.Join(words, " ")
	return cleanbody
}
