package service

type TimelineStore interface {
	Timeline(userID, limit int64) ([]int64, error)
}

type Service struct {
	store TimelineStore
}

func New(store TimelineStore) *Service {
	return &Service{store: store}
}

func (s *Service) Timeline(userID int64) ([]int64, error) {
	return s.store.Timeline(userID, 20)
}
