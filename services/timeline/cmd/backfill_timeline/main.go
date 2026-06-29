package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/SaputraUta/mini-twitter/services/timeline/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	ropts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.NewClient(ropts)
	pipe := rdb.Pipeline()
	defer rdb.Close()

	rows, err := pool.Query(ctx,
		`SELECT t.id, t.user_id, t.text, t.created_at FROM tweets t`)
	if err != nil {
		log.Fatal(err)
	}

	counter := 0
	for rows.Next() {
		var tweet model.Tweet
		if err := rows.Scan(&tweet.ID, &tweet.UserID, &tweet.Text, &tweet.CreatedAt); err != nil {
			log.Fatal(err)
		}
		key := fmt.Sprintf("tweet:%d", tweet.ID)
		marshalledTweet, err := json.Marshal(tweet)
		if err != nil {
			log.Fatal(err)
		}
		pipe.Set(ctx, key, marshalledTweet, 0).Err()
		counter++
		if counter%1000 == 0 {
			if _, err := pipe.Exec(ctx); err != nil {
				log.Fatal(err)
			}
		}
	}
	rows.Close()
	if _, err := pipe.Exec(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("Done")
}
