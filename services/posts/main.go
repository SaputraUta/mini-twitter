package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SaputraUta/mini-twitter/services/posts/internal/handler"
)

func main() {
	port := os.Getenv("POSTS_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	h := handler.New()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /tweet", h.CreateTweet)

	addr := ":" + port
	log.Printf("posts-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
