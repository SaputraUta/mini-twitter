package model

type Tweet struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Text   string `json:"text"`
}
