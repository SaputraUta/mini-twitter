package store

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisTimeline struct {
	client *redis.Client
}

func NewRedisTimeline(client *redis.Client) *RedisTimeline {
	return &RedisTimeline{client: client}
}

func (r *RedisTimeline) Timeline(userID int64, limit int64) ([]int64, error) {
	key := fmt.Sprintf("timeline:%d", userID)

	vals, err := r.client.LRange(context.Background(), key, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(vals))
	for _, v := range vals {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}
