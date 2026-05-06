package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Wayne-Francis/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}
	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
		Token     string    `json:"token"`
	}
	decoder := json.NewDecoder(r.Body)
	request := requestBody{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 400, "Bad Request")
		return
	}
	var expiresIn time.Duration
	if request.ExpiresInSeconds == nil {
		expiresIn = time.Hour
	} else {
		seconds := *request.ExpiresInSeconds
		if seconds > 3600 {
			seconds = 3600
		}
		expiresIn = time.Duration(seconds) * time.Second
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	match, err := auth.CheckPasswordHash(request.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	if !match {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}
	t, err := auth.MakeJWT(user.ID, cfg.secret, expiresIn)
	if err != nil {
		respondWithError(w, 500, "Could not generate Token")
		return
	}
	authUser := responseBody{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     t,
	}
	respondWithJSON(w, 200, authUser)
}
