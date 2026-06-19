package store

import (
	"sync"

	"github.com/SaputraUta/mini-twitter/services/posts/internal/model"
)

type MemoryStore struct {
	mu     sync.Mutex
	tweets []model.Tweet
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{nextID: 1}
}

func (s *MemoryStore) SaveTweet(t model.Tweet) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t.ID = s.nextID
	s.nextID++
	s.tweets = append(s.tweets, t)
	return t.ID, nil
}
