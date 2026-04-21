package main

import (
	"fmt"
	"net/http"
)

func main() {

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", healthzHandler)
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
}

func healthzHandler(w http.ResponseWriter, req *http.Request) {

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))

}
