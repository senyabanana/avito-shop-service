package entity

import "time"

type CoinsTransaction struct {
	ID        int       `json:"-"`
	FromUser  string    `json:"fromUser"`
	ToUser    string    `json:"toUser"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
