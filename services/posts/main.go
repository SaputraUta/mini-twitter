package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/SaputraUta/mini-twitter/services/posts/internal/events"
	"github.com/SaputraUta/mini-twitter/services/posts/internal/handler"
	"github.com/SaputraUta/mini-twitter/services/posts/internal/service"
	"github.com/SaputraUta/mini-twitter/services/posts/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
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

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	pub, err := events.NewNatsPublisher(nc)
	if err != nil {
		log.Fatal(err)
	}

	svc := service.New(st, pub)
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
