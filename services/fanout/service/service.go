package service

import "github.com/SaputraUta/mini-twitter/services/fanout/model"

type FollowerStore interface {
	Followers(authorID int64) ([]int64, error)
}

type TimelineStore interface {
	PushTweet(userID, tweetID int64) error
}

type Fanout struct {
	followers FollowerStore
	timelines TimelineStore
}

func New(followers FollowerStore, timelines TimelineStore) *Fanout {
	return &Fanout{followers: followers, timelines: timelines}
}

func (f *Fanout) Handle(ev model.TweetEvent) error {
	ids, err := f.followers.Followers(ev.UserID)
	if err != nil {
		return err
	}
	for _, uid := range ids {
		if err := f.timelines.PushTweet(uid, ev.ID); err != nil {
			return err
		}
	}
	return nil
}
