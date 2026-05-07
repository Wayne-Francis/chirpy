package main

import (
	"net/http"

	"github.com/Wayne-Francis/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Missing or invalid auth header")
		return
	}
	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "token doesn't exist, expired, or revoked")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
