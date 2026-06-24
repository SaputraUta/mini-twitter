package store

import (
	"context"

	"github.com/SaputraUta/mini-twitter/services/timeline/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresTweets struct {
	pool *pgxpool.Pool
}

func NewPostgresTweets(pool *pgxpool.Pool) *PostgresTweets {
	return &PostgresTweets{pool: pool}
}

func (s *PostgresTweets) TweetsByIDs(ids []int64) ([]model.Tweet, error) {
	rows, err := s.pool.Query(context.Background(), `SELECT id, user_id, text, created_at FROM tweets WHERE ID = ANY($1)`, ids)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tweets []model.Tweet
	for rows.Next() {
		var t model.Tweet
		if err := rows.Scan(&t.ID, &t.UserID, &t.Text, &t.CreatedAt); err != nil {
			return nil, err
		}
		tweets = append(tweets, t)
	}
	return tweets, rows.Err()
}
