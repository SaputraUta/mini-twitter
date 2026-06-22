package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresFollowers struct {
	pool *pgxpool.Pool
}

func NewPostgresFollowers(pool *pgxpool.Pool) *PostgresFollowers {
	return &PostgresFollowers{pool: pool}
}

func (s *PostgresFollowers) Followers(authorID int64) ([]int64, error) {
	rows, err := s.pool.Query(context.Background(),
		`SELECT follower_id FROM follows WHERE followee_id = $1`, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}
