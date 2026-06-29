package store

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/SaputraUta/mini-twitter/services/timeline/internal/model"
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

func (r *RedisTimeline) TweetsByIDs(ids []int64) ([]model.Tweet, error) {
	var keys []string
	var tweets []model.Tweet
	for _, id := range ids {
		key := fmt.Sprintf("tweet:%d", id)
		keys = append(keys, key)
	}

	rawTweets, err := r.client.MGet(context.Background(), keys...).Result()

	if err != nil {
		return nil, err
	}

	for _, tweetInterface := range rawTweets {
		var tweet model.Tweet
		if tweetInterface != nil {
			if tweetStr, ok := tweetInterface.(string); ok {
				tweetByte := []byte(tweetStr)
				err := json.Unmarshal(tweetByte, &tweet)

				if err != nil {
					return nil, err
				}
				tweets = append(tweets, tweet)
			}
		}
	}
	return tweets, nil
}
