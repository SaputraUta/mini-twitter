package service

import "github.com/SaputraUta/mini-twitter/services/posts/internal/model"

type Publisher interface {
	PublishTweetCreated(t model.Tweet) error
}

type TweetStore interface {
	SaveTweet(t model.Tweet) (int, error)
}

type Service struct {
	store TweetStore
	pub   Publisher
}

func New(store TweetStore, pub Publisher) *Service {
	return &Service{store: store, pub: pub}
}

func (s *Service) CreateTweet(t model.Tweet) (model.Tweet, error) {
	id, err := s.store.SaveTweet(t)
	if err != nil {
		return model.Tweet{}, err
	}
	t.ID = id

	if err := s.pub.PublishTweetCreated(t); err != nil {
		return model.Tweet{}, err
	}
	return t, nil
}
