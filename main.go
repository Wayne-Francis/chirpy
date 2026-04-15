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
	mux.Handle("/", http.FileServer(http.Dir(".")))
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
}
