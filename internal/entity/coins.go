package entity

import "time"

type CoinsTransaction struct {
	ID        int       `json:"-"`
	FromUser  string    `json:"from_user"`
	ToUser    string    `json:"to_user"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
