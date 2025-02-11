package entity

import "time"

type Inventory struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	MerchID   int       `json:"merch_id" db:"merch_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type MerchItems struct {
	ID       int    `json:"id" db:"id"`
	ItemType string `json:"item_type" db:"item_type"`
	Price    int    `json:"price" db:"price"`
}
