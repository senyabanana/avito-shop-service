package entity

import "time"

type InventoryItem struct {
	ID        int       `json:"-"`
	UserID    int       `json:"user_id"`
	ItemType  string    `json:"item_type"`
	Quantity  string    `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type MerchItem struct {
	ItemType string `json:"item_type"`
	Price    int    `json:"price"`
}

var MerchItems = map[string]int{
	"t-shirt":    80,
	"cup":        20,
	"book":       50,
	"pen":        10,
	"powerbank":  200,
	"hoody":      300,
	"umbrella":   200,
	"socks":      10,
	"wallet":     50,
	"pink-hoody": 500,
}
