package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Wayne-Francis/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	dbPF := os.Getenv("PLATFORM")
	dbs := os.Getenv("SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
		return
	}
	dbQueries := database.New(db)
	apiCfg := &apiConfig{
		db:       dbQueries,
		platform: dbPF,
		secret:   dbs,
	}
	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)
	mux.HandleFunc("POST /api/users", apiCfg.handlerUsers)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpbyId)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)
	mux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUsers)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirps)
	mux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
		return
	}
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
}
