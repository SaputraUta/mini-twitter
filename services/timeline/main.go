package main

import (
	"log"
	"net/http"
	"os"

	"github.com/SaputraUta/mini-twitter/services/timeline/internal/handler"
	"github.com/SaputraUta/mini-twitter/services/timeline/internal/service"
	"github.com/SaputraUta/mini-twitter/services/timeline/internal/store"
	"github.com/redis/go-redis/v9"
)

func main() {
	port := os.Getenv("TIMELINE_SERVICE_PORT")
	if port == "" {
		port = "8002"
	}

	ropts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.NewClient(ropts)
	defer rdb.Close()

	st := store.NewRedisTimeline(rdb)
	svc := service.New(st)
	h := handler.New(svc)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /timeline/{user}", h.Timeline)
	addr := ":" + port
	log.Printf("timeline-service listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
