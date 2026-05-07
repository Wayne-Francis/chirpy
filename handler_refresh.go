package main

import (
	"net/http"
	"time"

	"github.com/Wayne-Francis/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type responseBody struct {
		Token string `json:"token"`
	}
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid auth header")
		return
	}
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "token doesn't exist, expired, or revoked")
		return
	}
	newJWT, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "Could not generate Token")
		return
	}
	newToken := responseBody{
		Token: newJWT,
	}

	respondWithJSON(w, 200, newToken)
}
