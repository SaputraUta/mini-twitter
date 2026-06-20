package store

import (
	"context"

	"github.com/SaputraUta/mini-twitter/services/posts/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostresStore(pool *pgxpool.Pool) *PostgresStore {
	return &PostgresStore{pool: pool}
}

func (s *PostgresStore) SaveTweet(t model.Tweet) (int, error) {
	var id int
	err := s.pool.QueryRow(
		context.Background(),
		`INSERT INTO tweets (user_id, text) VALUES ($1, $2) RETURNING id`,
		t.UserID, t.Text,
	).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
