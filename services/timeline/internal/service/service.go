package service

import "github.com/SaputraUta/mini-twitter/services/timeline/internal/model"

type TimelineStore interface {
	Timeline(userID, limit int64) ([]int64, error)
}

type TweetStore interface {
	TweetsByIDs(ids []int64) ([]model.Tweet, error)
	TimelineFromFollows(userID int64, limit int64) ([]model.Tweet, error)
}

type Service struct {
	timelines TimelineStore
	tweets    TweetStore
}

func New(timelines TimelineStore, tweets TweetStore) *Service {
	return &Service{timelines: timelines, tweets: tweets}
}

func (s *Service) Timeline(userID int64) ([]model.Tweet, error) {
	ids, err := s.timelines.Timeline(userID, 20)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return []model.Tweet{}, nil
	}

	tweets, err := s.tweets.TweetsByIDs(ids)
	if err != nil {
		return nil, err
	}

	byID := make(map[int64]model.Tweet, len(tweets))
	for _, t := range tweets {
		byID[t.ID] = t
	}

	ordered := make([]model.Tweet, 0, len(ids))
	for _, id := range ids {
		if t, ok := byID[id]; ok {
			ordered = append(ordered, t)
		}
	}
	return ordered, nil
}

func (s *Service) TimelineFromDB(userID int64) ([]model.Tweet, error) {
	return s.tweets.TimelineFromFollows(userID, 20)
}
