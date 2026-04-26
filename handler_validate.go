package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Body string `json:"body"`
	}
	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
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
	respclean := responseBody{
		CleanedBody: cleanedBody,
	}
	respondWithJSON(w, 200, respclean)
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
