package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/SaputraUta/mini-twitter/services/posts/internal/handler"
	"github.com/SaputraUta/mini-twitter/services/posts/internal/service"
	"github.com/SaputraUta/mini-twitter/services/posts/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	port := os.Getenv("POSTS_SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	dbURL := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	st := store.NewPostresStore(pool)
	svc := service.New(st)
	h := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("POST /tweet", h.CreateTweet)

	addr := ":" + port
	log.Printf("posts-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
