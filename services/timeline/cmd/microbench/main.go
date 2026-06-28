// Micro-benchmark: isolate the DATA-LAYER cost of the two timeline strategies,
// with no HTTP / JSON-to-client / k6 in the way.
//
//	approach-1: Redis LRANGE(20) + Postgres hydrate (WHERE id = ANY)
//	approach-2: Postgres JOIN (tweets x follows, ORDER BY id DESC LIMIT 20)
//
// Both run under the same connection budget and concurrency so the only
// difference is the access path itself.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const (
	totalOps    = 50000
	concurrency = 50
	numUsers    = 10000
)

func main() {
	ctx := context.Background()

	cfg, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	cfg.MaxConns = concurrency
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	ropts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Fatal(err)
	}
	ropts.PoolSize = concurrency
	rdb := redis.NewClient(ropts)
	defer rdb.Close()

	// warm caches
	for u := int64(1); u <= 500; u++ {
		redisPath(ctx, pool, rdb, u)
		joinPath(ctx, pool, u)
	}

	run("approach-1 (Redis LRANGE + Postgres hydrate)", func(uid int64) { redisPath(ctx, pool, rdb, uid) })
	run("approach-2 (Postgres JOIN)", func(uid int64) { joinPath(ctx, pool, uid) })
}

func run(name string, op func(int64)) {
	durs := make([]time.Duration, totalOps)
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < totalOps; i++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(i int) {
			defer wg.Done()
			defer func() { <-sem }()
			uid := int64((i % numUsers) + 1)
			t0 := time.Now()
			op(uid)
			durs[i] = time.Since(t0)
		}(i)
	}
	wg.Wait()
	total := time.Since(start)

	sort.Slice(durs, func(a, b int) bool { return durs[a] < durs[b] })
	pct := func(q float64) time.Duration { return durs[int(float64(len(durs))*q)] }
	var sum time.Duration
	for _, d := range durs {
		sum += d
	}

	fmt.Printf("\n%s\n", name)
	fmt.Printf("  ops=%d  concurrency=%d\n", totalOps, concurrency)
	fmt.Printf("  avg=%-10v p50=%-10v p95=%-10v p99=%-10v max=%v\n",
		sum/totalOps, pct(0.50), pct(0.95), pct(0.99), durs[len(durs)-1])
	fmt.Printf("  throughput=%.0f ops/s\n", float64(totalOps)/total.Seconds())
}

func redisPath(ctx context.Context, pool *pgxpool.Pool, rdb *redis.Client, uid int64) {
	key := fmt.Sprintf("timeline:%d", uid)
	vals, err := rdb.LRange(ctx, key, 0, 19).Result()
	if err != nil {
		log.Fatal(err)
	}
	if len(vals) == 0 {
		return
	}
	ids := make([]int64, 0, len(vals))
	for _, v := range vals {
		id, _ := strconv.ParseInt(v, 10, 64)
		ids = append(ids, id)
	}
	rows, err := pool.Query(ctx,
		`SELECT id, user_id, text, created_at FROM tweets WHERE id = ANY($1)`, ids)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id, userID int64
		var text string
		var createdAt time.Time
		rows.Scan(&id, &userID, &text, &createdAt)
	}
	rows.Close()
}

func joinPath(ctx context.Context, pool *pgxpool.Pool, uid int64) {
	rows, err := pool.Query(ctx,
		`SELECT t.id, t.user_id, t.text, t.created_at
		 FROM tweets t JOIN follows f ON t.user_id = f.followee_id
		 WHERE f.follower_id = $1 ORDER BY t.id DESC LIMIT 20`, uid)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var id, userID int64
		var text string
		var createdAt time.Time
		rows.Scan(&id, &userID, &text, &createdAt)
	}
	rows.Close()
}
