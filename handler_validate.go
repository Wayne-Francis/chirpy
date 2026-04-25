package main

import (
	"encoding/json"
	"log"
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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}
