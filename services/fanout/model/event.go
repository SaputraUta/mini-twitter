package model

type TweetEvent struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Text   string `json:"text"`
}
