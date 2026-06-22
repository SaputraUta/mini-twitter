package store

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisTimeline struct {
	client *redis.Client
}

func NewRedisTimeline(client *redis.Client) *RedisTimeline {
	return &RedisTimeline{client: client}
}

func (r *RedisTimeline) PushTweet(userID, tweetID int64) error {
	key := fmt.Sprintf("timeline:%d", userID)
	ctx := context.Background()

	if err := r.client.LPush(ctx, key, tweetID).Err(); err != nil {
		return err
	}
	return r.client.LTrim(ctx, key, 0, 799).Err()
}
