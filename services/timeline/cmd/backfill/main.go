package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	numUsers     = 10000
	timelineSize = 400
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
	defer rdb.Close()

	for u := 1; u <= numUsers; u++ {
		rows, err := pool.Query(ctx,
			`SELECT t.id FROM tweets t
			 JOIN follows f ON t.user_id = f.followee_id
			 WHERE f.follower_id = $1
			 ORDER BY t.id DESC LIMIT $2`, u, timelineSize)
		if err != nil {
			log.Fatal(err)
		}

		var ids []interface{}
		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				log.Fatal(err)
			}
			ids = append(ids, id)
		}
		rows.Close()

		key := fmt.Sprintf("timeline:%d", u)
		rdb.Del(ctx, key)
		if len(ids) > 0 {
			if err := rdb.RPush(ctx, key, ids...).Err(); err != nil {
				log.Fatal(err)
			}
		}

		if u%1000 == 0 {
			log.Printf("backfilled %d/%d users", u, numUsers)
		}
	}
	log.Println("done")
}
