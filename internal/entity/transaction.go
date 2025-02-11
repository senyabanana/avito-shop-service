package entity

import "time"

type Transaction struct {
	ID        int       `json:"id" db:"id"`
	FromUser  int       `json:"from_user" db:"from_user"`
	ToUser    int       `json:"to_user" db:"to_user"`
	Amount    int       `json:"amount" db:"amount"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
