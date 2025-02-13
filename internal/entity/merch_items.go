package entity

type MerchItems struct {
	ID       int    `json:"id" db:"id"`
	ItemType string `json:"item_type" db:"item_type"`
	Price    int    `json:"price" db:"price"`
}
