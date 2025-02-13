package entity

type Inventory struct {
	ID       int `json:"id" db:"id"`
	UserID   int `json:"user_id" db:"user_id"`
	MerchID  int `json:"merch_id" db:"merch_id"`
	Quantity int `json:"quantity" db:"quantity"`
}

type MerchItems struct {
	ID       int    `json:"id" db:"id"`
	ItemType string `json:"item_type" db:"item_type"`
	Price    int    `json:"price" db:"price"`
}
