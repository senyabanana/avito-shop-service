package database

import (
	"fmt"

	"github.com/senyabanana/avito-shop-service/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	UsersTable        = "users"
	TransactionsTable = "transactions"
	InventoryTable    = "inventory"
	MerchItemsTable   = "merch_items"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
