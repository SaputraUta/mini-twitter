package service

import "github.com/SaputraUta/mini-twitter/services/posts/internal/model"

type TweetStore interface {
	SaveTweet(t model.Tweet) (int, error)
}

type Service struct {
	store TweetStore
}

func New(store TweetStore) *Service {
	return &Service{store: store}
}

func (s *Service) CreateTweet(t model.Tweet) (model.Tweet, error) {
	id, err := s.store.SaveTweet(t)
	if err != nil {
		return model.Tweet{}, err
	}
	t.ID = id
	return t, nil
}
