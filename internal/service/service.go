package service

import (
	"context"

	"github.com/senyabanana/avito-shop-service/internal/entity"
	"github.com/senyabanana/avito-shop-service/internal/repository"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
	"github.com/sirupsen/logrus"
)

type Authorization interface {
	GetUser(ctx context.Context, username string) (entity.User, error)
	CreateUser(ctx context.Context, username, password string) error
	GenerateToken(ctx context.Context, username, password string) (string, error)
	ParseToken(accessToken string) (int, error)
}

type Transaction interface {
	GetUserInfo(ctx context.Context, userID int) (entity.InfoResponse, error)
	SendCoin(ctx context.Context, fromUserID int, toUsername string, amount int) error
}

type Inventory interface {
	BuyItem(ctx context.Context, userID int, itemName string) error
}

type Service struct {
	Authorization
	Transaction
	Inventory
}

func NewService(repos *repository.Repository, trManager *manager.Manager, log *logrus.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos.UserRepository, trManager, log),
		Transaction:   NewTransactionService(repos.UserRepository, repos.TransactionRepository, repos.InventoryRepository, trManager, log),
		Inventory:     NewInventoryService(repos.UserRepository, repos.InventoryRepository, trManager, log),
	}
}
