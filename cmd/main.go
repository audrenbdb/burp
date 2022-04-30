package main

import (
	"burp"
	"burp/repo/repotest"
	"burp/rest/chi"
	"fmt"
	"log"
	"net/http"
)

func main() {
	addr := "localhost:8080"
	repo := repotest.FakeRepo
	brewer := &burp.Brewer{BeerRepo: repo}
	handler := chi.Handler(brewer)

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	fmt.Println("Starting server at", addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
