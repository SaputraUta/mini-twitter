package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/SaputraUta/mini-twitter/services/timeline/internal/handler"
	"github.com/SaputraUta/mini-twitter/services/timeline/internal/service"
	"github.com/SaputraUta/mini-twitter/services/timeline/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	port := os.Getenv("TIMELINE_SERVICE_PORT")
	if port == "" {
		port = "8002"
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	ropts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.NewClient(ropts)
	defer rdb.Close()

	timelineStore := store.NewRedisTimeline(rdb)
	tweetStore := store.NewPostgresTweets(pool)

	svc := service.New(timelineStore, tweetStore)
	h := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /timeline/{user}", h.Timeline)
	mux.HandleFunc("GET /timeline-db/{user}", h.TimelineDB)
	addr := ":" + port
	log.Printf("timeline-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
