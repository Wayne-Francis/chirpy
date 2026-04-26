package main

import (
	"encoding/json"
	"net/http"
	"time"

	"log"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Email string `json:"email"`
	}
	type responseBody struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	request := requestBody{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "couldn't decode parameters")
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), request.Email)
	if err != nil {
		log.Printf("couldn't create user: %s", err)
		respondWithError(w, 500, "couldn't create new user")
		return
	}
	new_user := responseBody{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     request.Email,
	}
	respondWithJSON(w, 201, new_user)
}
