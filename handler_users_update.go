package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Wayne-Francis/chirpy/internal/auth"
	"github.com/Wayne-Francis/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type responseBody struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid auth header")
		return
	}
	id, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, 401, "Invalid token")
		return
	}
	decoder := json.NewDecoder(r.Body)
	request := requestBody{}
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "couldn't decode parameters")
		return
	}
	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		respondWithError(w, 500, "couldn't hash password")
		return
	}
	update, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{Email: request.Email,
		HashedPassword: hashedPassword, ID: id})
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	updatedUser := responseBody{
		ID:          update.ID,
		CreatedAt:   update.CreatedAt,
		UpdatedAt:   update.UpdatedAt,
		Email:       update.Email,
		IsChirpyRed: update.IsChirpyRed,
	}
	respondWithJSON(w, 200, updatedUser)
}
