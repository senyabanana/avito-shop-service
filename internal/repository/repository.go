package repository

import (
	"context"

	"github.com/senyabanana/avito-shop-service/internal/entity"

	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type UserRepository interface {
	CreateUser(ctx context.Context, user entity.User) (int, error)
	GetUser(ctx context.Context, username string) (entity.User, error)
	GetUserBalance(ctx context.Context, userID int) (int, error)
	UpdateCoins(ctx context.Context, userID, amount int) error
}

type TransactionRepository interface {
	GetReceivedTransactions(ctx context.Context, userID int) ([]entity.TransactionDetail, error)
	GetSentTransactions(ctx context.Context, userID int) ([]entity.TransactionDetail, error)
	InsertTransaction(ctx context.Context, fromUserID, toUserID, amount int) error
}

type InventoryRepository interface {
	GetItem(ctx context.Context, itemName string) (entity.MerchItems, error)
	GetUserInventory(ctx context.Context, userID int) ([]entity.InventoryItem, error)
	GetInventoryItem(ctx context.Context, userID, merchID int) (int, error)
	UpdateInventoryItem(ctx context.Context, userID, merchID int) error
	InsertInventoryItem(ctx context.Context, userID, merchID int) error
}

type Repository struct {
	UserRepository
	TransactionRepository
	InventoryRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:        NewUserPostgres(db),
		TransactionRepository: NewTransactionPostgres(db),
		InventoryRepository:   NewInventoryPostgres(db),
	}
}
