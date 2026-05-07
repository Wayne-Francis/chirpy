package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Wayne-Francis/chirpy/internal/auth"
	"github.com/Wayne-Francis/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
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
	decoder := json.NewDecoder(r.Body)
	request := requestBody{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "couldn't decode parameters")
		return
	}
	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		respondWithError(w, 500, "couldn't hash password")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{Email: request.Email,
		HashedPassword: hashedPassword})
	if err != nil {
		log.Printf("couldn't create user: %s", err)
		respondWithError(w, 500, "couldn't create new user")
		return
	}
	new_user := responseBody{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       request.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	respondWithJSON(w, 201, new_user)
}
