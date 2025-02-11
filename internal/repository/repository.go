package repository

import (
	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetUser(username string) (entity.User, error)
}

type UserTransaction interface {
	GetUserBalance(userID int) (int, error)
	GetUserInventory(userID int) ([]entity.InventoryItem, error)
	GetReceivedTransactions(userID int) ([]entity.TransactionDetail, error)
	GetSentTransactions(userID int) ([]entity.TransactionDetail, error)
}

type Inventory interface {
}

type Repository struct {
	Authorization
	UserTransaction
	Inventory
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization:   NewAuthPostgres(db),
		UserTransaction: NewUserTransactionPostgres(db),
	}
}
