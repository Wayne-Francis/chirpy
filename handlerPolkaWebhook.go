package main

import (
	"encoding/json"
	"net/http"

	"github.com/Wayne-Francis/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid auth header")
		return
	}
	if key != cfg.polkaKey {
		respondWithError(w, 401, "Incorrect credentials")
		return
	}
	decoder := json.NewDecoder(r.Body)
	request := requestBody{}
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, 500, "couldn't decode parameters")
		return
	}
	if request.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	id, err := uuid.Parse(request.Data.UserID)
	if err != nil {
		respondWithError(w, 400, "couldn't parse Id")
		return
	}
	rowCount, err := cfg.db.UpgradeUser(r.Context(), id)
	if err != nil {
		respondWithError(w, 500, "No user Upgraded")
		return
	}
	if rowCount == 0 {
		respondWithError(w, 404, "No Users Upgraded")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
