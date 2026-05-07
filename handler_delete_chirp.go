package main

import (
	"net/http"

	"github.com/Wayne-Francis/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirps(w http.ResponseWriter, r *http.Request) {
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
	if id != chirp.UserID {
		respondWithError(w, 403, "Incorrect ID")
		return
	}
	err = cfg.db.DeleteChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, 404, "Chirp Not Found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
