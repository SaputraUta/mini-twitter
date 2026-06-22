package main

import (
	"context"
	"log"
	"os"

	"github.com/SaputraUta/mini-twitter/services/fanout/internal/consumer"
	"github.com/SaputraUta/mini-twitter/services/fanout/internal/service"
	"github.com/SaputraUta/mini-twitter/services/fanout/internal/store"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

func main() {
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

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}

	followers := store.NewPostgresFollowers(pool)
	timelines := store.NewRedisTimeline(rdb)
	svc := service.New(followers, timelines)
	c := consumer.New(js, svc)

	if err := c.Start(); err != nil {
		log.Fatal(err)
	}

	log.Println("fanout-worker running, waiting for tweets...")
	select {}
}
